package grpcsrv

import (
	"context"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
	commonv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/common/v1"
	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/repository"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/streamhub"
)

type fakeReader struct {
	mu       sync.Mutex
	list     []*models.Event
	forReq   map[string][]*models.Event
	queries  []repository.ListEventsFilter
}

func (f *fakeReader) ListEvents(_ context.Context, q repository.ListEventsFilter, _ uint32) ([]*models.Event, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queries = append(f.queries, q)

	out := make([]*models.Event, 0, len(f.list))
	for _, e := range f.list {
		if q.FromBlock > 0 && e.BlockNumber < q.FromBlock {
			continue
		}
		if q.ToBlock > 0 && e.BlockNumber > q.ToBlock {
			continue
		}
		if q.AssetID != nil && e.AssetID != *q.AssetID {
			continue
		}
		if len(q.Kinds) > 0 {
			ok := false
			for _, k := range q.Kinds {
				if k == e.Kind {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		out = append(out, e)
	}
	// Repo returns DESC by (block_number, log_index) — mirror that
	// here so the server's reverse-to-ASC step lands correctly.
	sort.Slice(out, func(i, j int) bool {
		if out[i].BlockNumber == out[j].BlockNumber {
			return out[i].LogIndex > out[j].LogIndex
		}
		return out[i].BlockNumber > out[j].BlockNumber
	})
	return out, nil
}

func (f *fakeReader) EventsForRequest(_ context.Context, reqID *big.Int, _ uint32) ([]*models.Event, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.forReq[reqID.String()], nil
}

func sampleEvent(id int64, block uint64, kind models.EventKind, asset common.Hash) *models.Event {
	e := &models.Event{
		ID:              id,
		Kind:            kind,
		BlockNumber:     block,
		LogIndex:        uint32(id),
		AssetID:         asset,
		BlockHash:       common.HexToHash("0xbb"),
		TxHash:          common.HexToHash("0xaa"),
		ContractAddress: common.HexToAddress("0xcc"),
		ObservedAt:      time.Unix(int64(1_700_000_000+block), 0).UTC(),
	}
	switch kind {
	case models.EventKindPriceRequested:
		e.ReqID = big.NewInt(id)
		e.PriceRequested = &models.PriceRequestedPayload{
			ReqID:     big.NewInt(id),
			AssetID:   asset,
			Requester: common.HexToAddress("0xdd"),
		}
	case models.EventKindPriceFulfilled:
		e.ReqID = big.NewInt(id)
		e.PriceFulfilled = &models.PriceFulfilledPayload{
			ReqID:     big.NewInt(id),
			AssetID:   asset,
			Price:     big.NewInt(100 + id),
			Timestamp: big.NewInt(1700000000),
		}
	case models.EventKindAssetRegistered:
		e.AssetRegistered = &models.AssetRegisteredPayload{AssetID: asset, Aggregator: common.HexToAddress("0xee")}
	}
	return e
}

func startTestServer(t *testing.T, reader EventReader, hub StreamHub) (indexerv1.IndexerServiceClient, func()) {
	t.Helper()
	srv := New(&config.GRPCConfig{
		Host:           "127.0.0.1",
		Port:           0,
		MaxRecvMsgSize: 4 * 1024 * 1024,
		MaxSendMsgSize: 4 * 1024 * 1024,
	}, reader, hub, 5)

	if err := srv.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	addr := srv.Addr()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	cleanup := func() {
		_ = conn.Close()
		_ = srv.Stop(context.Background())
	}
	return indexerv1.NewIndexerServiceClient(conn), cleanup
}

func TestListEvents_FilteringAndOrdering(t *testing.T) {
	asset := common.HexToHash("0xa1")
	other := common.HexToHash("0xa2")
	reader := &fakeReader{list: []*models.Event{
		sampleEvent(1, 100, models.EventKindPriceRequested, asset),
		sampleEvent(2, 110, models.EventKindPriceFulfilled, asset),
		sampleEvent(3, 120, models.EventKindPriceRequested, other),
	}}
	hub := streamhub.New(8, nil)
	defer hub.Shutdown()

	cli, cleanup := startTestServer(t, reader, hub)
	defer cleanup()

	resp, err := cli.ListEvents(context.Background(), &indexerv1.ListEventsRequest{
		Kinds:   []indexerv1.EventKind{indexerv1.EventKind_EVENT_KIND_PRICE_REQUESTED},
		AssetId: asset.Hex(),
	})
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(resp.Events) != 1 || resp.Events[0].GetPriceRequested() == nil {
		t.Errorf("ListEvents response = %+v", resp.Events)
	}

	// Page request.
	resp, err = cli.ListEvents(context.Background(), &indexerv1.ListEventsRequest{
		Page: &commonv1.PageRequest{Page: 1, PageSize: 100},
	})
	if err != nil {
		t.Fatalf("ListEvents page: %v", err)
	}
	if len(resp.Events) != 3 {
		t.Errorf("got %d events, want 3", len(resp.Events))
	}
}

func TestListEvents_RejectsBadAssetID(t *testing.T) {
	cli, cleanup := startTestServer(t, &fakeReader{}, streamhub.New(1, nil))
	defer cleanup()

	_, err := cli.ListEvents(context.Background(), &indexerv1.ListEventsRequest{AssetId: "not-an-asset"})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("got %v, want InvalidArgument", err)
	}
}

func TestGetRequest_PendingAndFulfilled(t *testing.T) {
	asset := common.HexToHash("0xa1")
	pendingReq := sampleEvent(7, 100, models.EventKindPriceRequested, asset)
	fulfilled := sampleEvent(8, 110, models.EventKindPriceFulfilled, asset)
	fulfilled.ReqID = big.NewInt(7)
	fulfilled.PriceFulfilled.ReqID = big.NewInt(7)

	reader := &fakeReader{
		forReq: map[string][]*models.Event{
			"7":  {pendingReq, fulfilled},
			"99": {pendingReq},
		},
	}
	cli, cleanup := startTestServer(t, reader, streamhub.New(1, nil))
	defer cleanup()

	// Fulfilled
	got, err := cli.GetRequest(context.Background(), &indexerv1.GetRequestRequest{ReqId: "7"})
	if err != nil {
		t.Fatalf("GetRequest: %v", err)
	}
	if got.Status != indexerv1.RequestStatus_STATUS_FULFILLED {
		t.Errorf("status = %v, want FULFILLED", got.Status)
	}
	if got.FulfilledPrice == "" {
		t.Errorf("FulfilledPrice missing")
	}

	// Pending
	got, err = cli.GetRequest(context.Background(), &indexerv1.GetRequestRequest{ReqId: "99"})
	if err != nil {
		t.Fatalf("GetRequest pending: %v", err)
	}
	if got.Status != indexerv1.RequestStatus_STATUS_PENDING {
		t.Errorf("status = %v, want PENDING", got.Status)
	}

	// Not found
	_, err = cli.GetRequest(context.Background(), &indexerv1.GetRequestRequest{ReqId: "12345"})
	if status.Code(err) != codes.NotFound {
		t.Errorf("got %v, want NotFound", err)
	}

	// Bad reqId
	_, err = cli.GetRequest(context.Background(), &indexerv1.GetRequestRequest{ReqId: ""})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("empty reqId: got %v, want InvalidArgument", err)
	}
	_, err = cli.GetRequest(context.Background(), &indexerv1.GetRequestRequest{ReqId: "garbage"})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("garbage reqId: got %v, want InvalidArgument", err)
	}
}

func TestStreamEvents_LiveOnly(t *testing.T) {
	hub := streamhub.New(8, nil)
	defer hub.Shutdown()
	cli, cleanup := startTestServer(t, &fakeReader{}, hub)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := cli.StreamEvents(ctx, &indexerv1.StreamEventsRequest{})
	if err != nil {
		t.Fatalf("StreamEvents: %v", err)
	}

	// Give the server time to attach the subscription.
	time.Sleep(50 * time.Millisecond)

	asset := common.HexToHash("0xa1")
	hub.Publish(sampleEvent(1, 100, models.EventKindPriceRequested, asset))

	got, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if got.GetPriceRequested() == nil {
		t.Errorf("expected PriceRequested payload, got %+v", got)
	}
}

func TestStreamEvents_ReplayThenLive(t *testing.T) {
	asset := common.HexToHash("0xa1")
	reader := &fakeReader{list: []*models.Event{
		sampleEvent(1, 100, models.EventKindPriceRequested, asset),
		sampleEvent(2, 105, models.EventKindPriceFulfilled, asset),
	}}
	hub := streamhub.New(8, nil)
	defer hub.Shutdown()
	cli, cleanup := startTestServer(t, reader, hub)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := cli.StreamEvents(ctx, &indexerv1.StreamEventsRequest{FromBlock: 100})
	if err != nil {
		t.Fatalf("StreamEvents: %v", err)
	}

	// First two are the replay (ASC by block).
	first, err := stream.Recv()
	if err != nil {
		t.Fatalf("first Recv: %v", err)
	}
	if first.Meta.BlockNumber != 100 {
		t.Errorf("replay[0] block = %d, want 100", first.Meta.BlockNumber)
	}
	second, err := stream.Recv()
	if err != nil {
		t.Fatalf("second Recv: %v", err)
	}
	if second.Meta.BlockNumber != 105 {
		t.Errorf("replay[1] block = %d, want 105", second.Meta.BlockNumber)
	}

	// Now publish a live event at block 200 — should pass the
	// highest-replayed gate (105) and reach the client.
	hub.Publish(sampleEvent(3, 200, models.EventKindPriceRequested, asset))

	live, err := stream.Recv()
	if err != nil {
		t.Fatalf("live Recv: %v", err)
	}
	if live.Meta.BlockNumber != 200 {
		t.Errorf("live block = %d, want 200", live.Meta.BlockNumber)
	}
}
