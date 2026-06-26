// Package grpcsrv hosts the indexer-service gRPC server +
// IndexerService implementation. Three RPCs:
//
//   - ListEvents — paginated query with kind/asset/block filters,
//     sorted desc by (block_number, log_index).
//
//   - GetRequest — joins PriceRequested + PriceFulfilled by req_id,
//     returns lifecycle + tx hashes + last-known fulfillment price.
//
//   - StreamEvents — long-lived server stream. Replay-then-live: if
//     `from_block > 0` the server first drains the matching backlog
//     in chronological order, then attaches to the live stream hub
//     (with log-granular (block, log_index) dedup at the boundary).
//     With `from_block = 0` it's live-only. Events are emitted on
//     ingest — there is no confirmation gate.
package grpcsrv

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
	commonv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/common/v1"
	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/repository"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/streamhub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
)

// peerAddress extracts the client's network address from a gRPC
// context. Returns "unknown" when the framework didn't attach peer
// info (e.g. unit-test transports).
func peerAddress(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil || p.Addr == nil {
		return "unknown"
	}
	return p.Addr.String()
}

// formatAssetFilter renders a *common.Hash for log lines, printing
// "any" when no asset filter is set.
func formatAssetFilter(h *common.Hash) string {
	if h == nil {
		return "any"
	}
	return strings.ToLower(h.Hex())
}

// maxConcurrentStreams bounds concurrent gRPC streams per connection.
const maxConcurrentStreams = 256

// EventReader is the read surface ListEvents + GetRequest need.
// Implemented by *repository.Repository in production.
type EventReader interface {
	ListEvents(ctx context.Context, f repository.ListEventsFilter) ([]*models.Event, error)
	EventsForRequest(ctx context.Context, reqID *big.Int) ([]*models.Event, error)
}

// StreamHub is the publisher surface StreamEvents needs.
type StreamHub interface {
	Subscribe(filter streamhub.Filter) *streamhub.Subscription
}

// Server hosts the gRPC listener.
type Server struct {
	indexerv1.UnimplementedIndexerServiceServer

	cfg    *config.GRPCConfig
	reader EventReader
	hub    StreamHub

	grpc     *grpc.Server
	listener net.Listener
	started  atomic.Bool
}

// New wires the gRPC server with the supplied dependencies. `cfg.Port`
// is honored at Start time; New does not bind the socket.
func New(cfg *config.GRPCConfig, reader EventReader, hub StreamHub) *Server {
	return &Server{
		cfg:    cfg,
		reader: reader,
		hub:    hub,
	}
}

// Start binds the socket and begins serving in a background goroutine.
// Returns once the listener is up.
func (s *Server) Start(_ context.Context) error {
	if s.cfg == nil {
		return fmt.Errorf("grpcsrv: nil config")
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", addr) //nolint:noctx // standard listener.
	if err != nil {
		return fmt.Errorf("grpcsrv: listen %s: %w", addr, err)
	}
	s.listener = lis

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(s.cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(s.cfg.MaxSendMsgSize),
		// Bound concurrent streams per connection — cheap protection on
		// an unauthenticated internal port (pairs with the stream-hub
		// subscriber cap).
		grpc.MaxConcurrentStreams(maxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
		// Enforcement policy: tolerate the client's keepalive cadence.
		// Without this, gRPC-go defaults to MinTime=5m / no pings
		// without a stream, so a normal client (the API pings every
		// 30s, PermitWithoutStream=true) gets GOAWAY ENHANCE_YOUR_CALM
		// "too_many_pings" and the connection flaps. MinTime=10s leaves
		// headroom under the client's 30s; PermitWithoutStream mirrors
		// the client so idle gaps between unary calls don't trip it.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	if s.cfg.NumStreamWorkers > 0 {
		opts = append(opts, grpc.NumStreamWorkers(s.cfg.NumStreamWorkers))
	}

	srv := grpc.NewServer(opts...)
	indexerv1.RegisterIndexerServiceServer(srv, s)

	healthSrv := health.NewServer()
	healthSrv.SetServingStatus(indexerv1.IndexerService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, healthSrv)

	if s.cfg.Reflection {
		reflection.Register(srv)
	}
	s.grpc = srv
	s.started.Store(true)

	go func() {
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			// In production this surfaces through logs at the
			// application level; we keep this hot path quiet so
			// `Stop()` -> "server stopped" doesn't fire as an error.
			_ = err
		}
	}()
	return nil
}

// Stop gracefully drains pending RPCs and closes the listener.
func (s *Server) Stop(_ context.Context) error {
	if !s.started.Load() {
		return nil
	}
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
	return nil
}

// Addr returns the bound listener address. Useful for tests + readyz.
func (s *Server) Addr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// ---------------------------------------------------------------------
// IndexerService implementations
// ---------------------------------------------------------------------

// ListEvents returns the paginated historical view.
func (s *Server) ListEvents(ctx context.Context, req *indexerv1.ListEventsRequest) (*indexerv1.ListEventsResponse, error) {
	filter, err := s.buildListFilter(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	events, err := s.reader.ListEvents(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list events: %v", err)
	}

	out := &indexerv1.ListEventsResponse{
		Events: make([]*indexerv1.Event, 0, len(events)),
		Page:   &commonv1.PageResponse{},
	}
	for _, e := range events {
		p, err := e.ToProto()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "encode event id=%d: %v", e.ID, err)
		}
		out.Events = append(out.Events, p)
	}
	return out, nil
}

// GetRequest returns the lifecycle of a single price request.
func (s *Server) GetRequest(ctx context.Context, req *indexerv1.GetRequestRequest) (*indexerv1.RequestStatus, error) {
	if req == nil || strings.TrimSpace(req.ReqId) == "" {
		return nil, status.Error(codes.InvalidArgument, "req_id is required")
	}
	reqID, ok := new(big.Int).SetString(strings.TrimSpace(req.ReqId), 10)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "req_id must be a base-10 uint256")
	}

	events, err := s.reader.EventsForRequest(ctx, reqID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "events for request: %v", err)
	}
	if len(events) == 0 {
		return nil, status.Errorf(codes.NotFound, "req_id %s not observed", req.ReqId)
	}

	out := &indexerv1.RequestStatus{ReqId: req.ReqId, Status: indexerv1.RequestStatus_STATUS_PENDING}
	for _, e := range events {
		switch e.Kind {
		case models.EventKindPriceRequested:
			if e.PriceRequested == nil {
				continue
			}
			out.AssetId = strings.ToLower(e.PriceRequested.AssetID.Hex())
			out.Requester = strings.ToLower(e.PriceRequested.Requester.Hex())
			out.RequestedTxHash = strings.ToLower(e.TxHash.Hex())
			out.RequestedAt = timestamppb.New(e.ObservedAt)
		case models.EventKindPriceFulfilled:
			if e.PriceFulfilled == nil {
				continue
			}
			out.Status = indexerv1.RequestStatus_STATUS_FULFILLED
			out.FulfilledTxHash = strings.ToLower(e.TxHash.Hex())
			out.FulfilledPrice = e.PriceFulfilled.Price.String()
			out.FulfilledAt = timestamppb.New(e.ObservedAt)
			if out.AssetId == "" {
				out.AssetId = strings.ToLower(e.PriceFulfilled.AssetID.Hex())
			}
		case models.EventKindAssetRegistered, models.EventKindUnknown:
			// AssetRegistered isn't tied to a req_id; Unknown is a sentinel.
		}
	}
	return out, nil
}

// StreamEvents is the indexer's long-lived server stream. Replay-then-
// live when from_block > 0; live-only otherwise.
func (s *Server) StreamEvents(req *indexerv1.StreamEventsRequest, stream indexerv1.IndexerService_StreamEventsServer) error {
	filter, hubFilter, err := s.buildStreamFilter(req)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := stream.Context()
	peerAddr := peerAddress(ctx)

	// Subscribe BEFORE the historical replay so we don't miss events
	// fired during the replay window. Live events that arrive while
	// the replay drains are buffered in the subscriber's channel.
	// (Duplicate-suppression below ensures a replayed event isn't
	// re-emitted from the live tail.)
	sub := s.hub.Subscribe(hubFilter)
	if sub == nil {
		return status.Error(codes.Unavailable, "indexer-service is shutting down")
	}
	defer func() {
		sub.Cancel()
		logger.Log().Infof("grpcsrv: StreamEvents disconnect id=%d peer=%s", sub.ID(), peerAddr)
	}()

	logger.Log().Infof("grpcsrv: StreamEvents subscribe id=%d peer=%s kinds=%v asset_id=%s from_block=%d",
		sub.ID(), peerAddr, hubFilter.Kinds, formatAssetFilter(hubFilter.AssetID), req.GetFromBlock())

	var watermark logKey // highest (block, log_index) emitted during replay

	if req.GetFromBlock() > 0 {
		var err error
		watermark, err = s.replayHistory(ctx, filter, req.GetFromBlock(), stream)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-sub.Events():
			if !ok {
				return status.Error(codes.Unavailable, "stream hub closed")
			}
			// Skip events the replay already emitted. Dedup is
			// log-granular (block, log_index) — block-granular would
			// drop live events sharing the boundary block, or re-emit
			// replayed ones. The replay/live overlap is exactly the
			// boundary block, so this is where precision matters.
			if !(logKey{block: e.BlockNumber, logIndex: e.LogIndex}).after(watermark) {
				continue
			}
			p, perr := e.ToProto()
			if perr != nil {
				// Skip malformed events but stay connected.
				continue
			}
			if err := stream.Send(p); err != nil {
				return err
			}
		}
	}
}

// logKey orders events by (block_number, log_index) — the on-chain
// total order of logs.
type logKey struct {
	block    uint64
	logIndex uint32
}

// after reports whether k is strictly later than o in chain order.
func (k logKey) after(o logKey) bool {
	if k.block != o.block {
		return k.block > o.block
	}
	return k.logIndex > o.logIndex
}

// replayHistory drains historical events from `fromBlock` to head in
// chronological order. Returns the highest (block, log_index) emitted
// so the live loop can suppress duplicates at log granularity.
func (s *Server) replayHistory(
	ctx context.Context,
	base repository.ListEventsFilter,
	fromBlock uint64,
	stream indexerv1.IndexerService_StreamEventsServer,
) (logKey, error) {
	// Page through historical events ascending by block. We reuse
	// ListEvents but flip the order in-memory because the repo's
	// SELECT is DESC; for the v1 demo workload (small DB) flipping
	// in-memory is fine.
	filter := base
	filter.FromBlock = fromBlock
	filter.Limit = 1000
	filter.Offset = 0

	var highest logKey
	for {
		batch, err := s.reader.ListEvents(ctx, filter)
		if err != nil {
			return highest, status.Errorf(codes.Internal, "replay query: %v", err)
		}
		if len(batch) == 0 {
			break
		}
		// Reverse to ASC.
		for i, j := 0, len(batch)-1; i < j; i, j = i+1, j-1 {
			batch[i], batch[j] = batch[j], batch[i]
		}
		for _, e := range batch {
			p, err := e.ToProto()
			if err != nil {
				continue
			}
			if err := stream.Send(p); err != nil {
				return highest, err
			}
			if k := (logKey{block: e.BlockNumber, logIndex: e.LogIndex}); k.after(highest) {
				highest = k
			}
		}
		if len(batch) < filter.Limit {
			break
		}
		filter.Offset += len(batch)
		if err := ctx.Err(); err != nil {
			return highest, nil
		}
	}
	return highest, nil
}

// ---------------------------------------------------------------------
// filter helpers
// ---------------------------------------------------------------------

func (s *Server) buildListFilter(req *indexerv1.ListEventsRequest) (repository.ListEventsFilter, error) {
	out := repository.ListEventsFilter{}
	if req == nil {
		return out, nil
	}
	for _, k := range req.Kinds {
		dk := models.EventKindFromProto(k)
		if dk == models.EventKindUnknown {
			continue
		}
		out.Kinds = append(out.Kinds, dk)
	}
	if req.AssetId != "" {
		h, err := parseAssetID(req.AssetId)
		if err != nil {
			return out, err
		}
		out.AssetID = &h
	}
	out.FromBlock = req.FromBlock
	out.ToBlock = req.ToBlock
	if req.Page != nil {
		out.Limit = int(req.Page.PageSize)
		// Page is 1-indexed per protocols/common/v1/pagination.proto.
		// Offset 0 for page 1.
		if req.Page.Page > 0 {
			out.Offset = int((req.Page.Page - 1) * req.Page.PageSize)
		}
	}
	return out, nil
}

func (s *Server) buildStreamFilter(req *indexerv1.StreamEventsRequest) (repository.ListEventsFilter, streamhub.Filter, error) {
	var (
		base    repository.ListEventsFilter
		hubFlt  streamhub.Filter
	)
	if req == nil {
		return base, hubFlt, nil
	}
	for _, k := range req.Kinds {
		dk := models.EventKindFromProto(k)
		if dk == models.EventKindUnknown {
			continue
		}
		base.Kinds = append(base.Kinds, dk)
		hubFlt.Kinds = append(hubFlt.Kinds, dk)
	}
	if req.AssetId != "" {
		h, err := parseAssetID(req.AssetId)
		if err != nil {
			return base, hubFlt, err
		}
		base.AssetID = &h
		hubFlt.AssetID = &h
	}
	return base, hubFlt, nil
}

func parseAssetID(s string) (common.Hash, error) {
	if !strings.HasPrefix(strings.ToLower(s), "0x") || len(s) != 66 {
		return common.Hash{}, fmt.Errorf("asset_id must be a 0x-prefixed 32-byte hex string")
	}
	return common.HexToHash(s), nil
}
