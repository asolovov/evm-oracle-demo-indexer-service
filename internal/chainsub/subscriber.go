package chainsub

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/contracts/oracleregistry"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/contracts/priceaggregator"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
)

// Backoff bounds for the reconnect loop. Hardcoded consts (not config):
// public-RPC WS resets are the steady state, and these are operational
// constants no demo operator needs to tune.
const (
	minBackoff   = time.Second
	maxBackoff   = 30 * time.Second
	healthyReset = 60 * time.Second // a session living this long resets backoff
	roundIDLimit = 3 * time.Second  // bound on the best-effort latestRoundData call
)

// EventStore is the persistence surface the subscriber needs.
// Implemented by *repository.Repository.
type EventStore interface {
	InsertEvent(ctx context.Context, e *models.Event) (bool, error)
	UpsertAggregator(ctx context.Context, aggregator common.Address, assetID common.Hash) error
	AggregatorRegistry(ctx context.Context) (map[common.Address]common.Hash, error)
	ChainCursor(ctx context.Context) (*models.ChainCursor, error)
	UpdateChainCursor(ctx context.Context, block uint64) error
}

// Publisher fans a freshly-ingested event out to live StreamEvents
// subscribers. Implemented by *streamhub.Hub. Nil-safe at the call site.
type Publisher interface {
	Publish(e *models.Event) int
}

// Metrics is the optional metrics surface; nil-safe.
type Metrics interface {
	ObserveSeen(kind models.EventKind)
	ObserveDecodeError()
	ObserveLagBlocks(v float64)
}

// logFetcher is the slice of the chain client catch-up needs. Matches
// *ethclient.Client; narrowed to an interface so catch-up is unit-
// testable without a node.
type logFetcher interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// Subscriber is the single chain-observer. One goroutine (Run) owns the
// chain client end-to-end: it dials, catches up the gap since the
// persisted cursor via eth_getLogs, live-subscribes, and on any
// disconnect reconnects (exponential backoff + jitter) and catches up
// again. Because the client never leaves this goroutine, there is no
// shared mutable client state and therefore no data race.
//
// There is NO confirmation gate: an event is published the moment it is
// persisted (emit-on-ingest). Reorg exposure is a documented downstream
// contract — a consumer that takes irreversible action owns its own
// finality guard.
type Subscriber struct {
	wsURL        string
	rpcURL       string
	registryAddr common.Address
	store        EventStore
	publisher    Publisher
	metrics      Metrics

	defaultStart uint64
	chunkSize    uint64
	cursorEvery  uint64

	// mapping (aggregator -> asset_id) is only mutated inside the Run
	// goroutine, but AssetIDFor is part of the AggregatorResolver
	// contract, so it stays mutex-guarded as a defensive measure.
	mu      sync.RWMutex
	mapping map[common.Address]common.Hash
}

// Config bundles the constructor knobs.
type Config struct {
	WSURL           string
	RPCURL          string
	RegistryAddress common.Address
	DefaultStart    uint64 // cursor seed on a cold DB (chain_cursor == 0)
	ChunkSize       uint64 // eth_getLogs catch-up chunk size (default 1000)
	CursorEvery     uint64 // persist the cursor every N blocks (default 100)
	Metrics         Metrics
}

// New builds a Subscriber. It does NOT dial; Run() owns the connection.
func New(store EventStore, publisher Publisher, cfg Config) *Subscriber {
	if cfg.ChunkSize == 0 {
		cfg.ChunkSize = 1000
	}
	if cfg.CursorEvery == 0 {
		cfg.CursorEvery = 100
	}
	return &Subscriber{
		wsURL:        cfg.WSURL,
		rpcURL:       cfg.RPCURL,
		registryAddr: cfg.RegistryAddress,
		store:        store,
		publisher:    publisher,
		metrics:      cfg.Metrics,
		defaultStart: cfg.DefaultStart,
		chunkSize:    cfg.ChunkSize,
		cursorEvery:  cfg.CursorEvery,
		mapping:      make(map[common.Address]common.Hash),
	}
}

// AssetIDFor implements AggregatorResolver.
func (s *Subscriber) AssetIDFor(addr common.Address) (common.Hash, bool) {
	s.mu.RLock()
	h, ok := s.mapping[addr]
	s.mu.RUnlock()
	return h, ok
}

// Run is the blocking reconnect loop. Returns nil on ctx cancellation.
func (s *Subscriber) Run(ctx context.Context) error {
	attempt := 0
	for {
		if ctx.Err() != nil {
			return nil
		}
		started := time.Now()
		err := s.runSession(ctx)
		if ctx.Err() != nil {
			return nil
		}
		// A session that stayed up past healthyReset is "healthy" —
		// reset the backoff so an occasional blip doesn't permanently
		// inflate the delay. A fast-dying session escalates.
		if time.Since(started) >= healthyReset {
			attempt = 0
		}
		delay := backoffFor(attempt)
		attempt++
		logger.Log().Warnf("chainsub: session ended (%v); reconnecting in %s (attempt %d)", err, delay.Round(time.Millisecond), attempt)
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
		}
	}
}

// backoffFor returns a full-jittered delay in [0, min(maxBackoff,
// minBackoff*2^attempt)].
func backoffFor(attempt int) time.Duration {
	capDelay := minBackoff << min(attempt, 16)
	if capDelay <= 0 || capDelay > maxBackoff {
		capDelay = maxBackoff
	}
	//nolint:gosec // jitter does not need crypto-grade randomness.
	return time.Duration(rand.Int63n(int64(capDelay) + 1))
}

// runSession owns one connection lifetime. The client is created and
// used entirely within this call stack — no struct field, no sharing —
// which is what removes the data race the previous design had.
func (s *Subscriber) runSession(ctx context.Context) error {
	wsClient, err := ethclient.DialContext(ctx, s.wsURL)
	if err != nil {
		return fmt.Errorf("dial ws: %w", err)
	}
	defer wsClient.Close()

	rpcClient := wsClient
	if s.rpcURL != "" && s.rpcURL != s.wsURL {
		rpcClient, err = ethclient.DialContext(ctx, s.rpcURL)
		if err != nil {
			return fmt.Errorf("dial rpc: %w", err)
		}
		defer rpcClient.Close()
	}

	cached, err := s.store.AggregatorRegistry(ctx)
	if err != nil {
		return fmt.Errorf("load cached registry: %w", err)
	}
	s.replaceMapping(cached)
	if rerr := s.refreshFromRegistry(ctx, rpcClient); rerr != nil {
		logger.Log().Warnf("chainsub: registry refresh failed: %v (continuing with cached mapping of %d)", rerr, len(cached))
	}

	parser, err := NewParser(s.registryAddr, s)
	if err != nil {
		return fmt.Errorf("init parser: %w", err)
	}
	addresses := s.subscribedAddresses()

	// Subscribe live FIRST into a buffered channel, so logs emitted
	// while catch-up runs are not lost. The live tail and the catch-up
	// range overlap at the head boundary; InsertEvent is idempotent on
	// (tx_hash, log_index), so the overlap is harmless and the gap is
	// closed.
	logs := make(chan types.Log, 1024)
	sub, err := wsClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{Addresses: addresses}, logs)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	head, err := rpcClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("head lookup: %w", err)
	}
	if err := s.catchUp(ctx, rpcClient, parser, addresses, head.Number.Uint64()); err != nil {
		return fmt.Errorf("catch-up: %w", err)
	}

	logger.Log().Infof("chainsub: live on %s — %d address(es), caught up through %d", s.wsURL, len(addresses), head.Number.Uint64())
	return s.runLive(ctx, rpcClient, parser, sub, logs)
}

// runLive drains the live subscription until the connection dies or ctx
// is canceled.
func (s *Subscriber) runLive(
	ctx context.Context,
	rpc *ethclient.Client,
	parser *Parser,
	sub ethereum.Subscription,
	logs <-chan types.Log,
) error {
	lastCursor := uint64(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			if err == nil {
				return errors.New("subscription closed")
			}
			return fmt.Errorf("subscription error: %w", err)
		case lg := <-logs:
			if err := s.ingest(ctx, rpc, parser, lg, true); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				logger.Log().Warnf("chainsub: ingest live log (tx=%s idx=%d): %v", lg.TxHash.Hex(), lg.Index, err)
			}
			// Advance the cursor as we go so the next reconnect's
			// catch-up window stays small. Park at the PREVIOUS block —
			// logs arrive in order, so block-1 is fully drained.
			if lg.BlockNumber > 0 && lg.BlockNumber-1 > lastCursor && lg.BlockNumber-1-lastCursor >= s.cursorEvery {
				if err := s.store.UpdateChainCursor(ctx, lg.BlockNumber-1); err != nil {
					logger.Log().Warnf("chainsub: advance cursor to %d: %v", lg.BlockNumber-1, err)
				} else {
					lastCursor = lg.BlockNumber - 1
				}
			}
		}
	}
}

// catchUp replays [cursor+1, head] in per-address chunks via
// eth_getLogs. Per-address (not multi-address) because some public RPC
// providers block multi-address filters. A per-(chunk,address) failure
// is a warning, not terminal; the cursor parks at the last fully-clean
// block so the gap is retried on the next connect.
func (s *Subscriber) catchUp(ctx context.Context, client logFetcher, parser *Parser, addresses []common.Address, head uint64) error {
	cursor, err := s.store.ChainCursor(ctx)
	if err != nil {
		return fmt.Errorf("read chain cursor: %w", err)
	}
	from := cursor.LastProcessedBlock + 1
	if cursor.LastProcessedBlock == 0 && s.defaultStart > 0 {
		from = s.defaultStart
	}
	if from > head {
		s.observeLag(0)
		return nil
	}

	logger.Log().Infof("chainsub: catch-up [%d..%d] in %d-block chunks (%d address(es))", from, head, s.chunkSize, len(addresses))

	persistAt := from + s.cursorEvery - 1
	failedChunks := 0
	highestClean := from - 1

	for from <= head {
		if err := ctx.Err(); err != nil {
			return err
		}
		to := from + s.chunkSize - 1
		if to > head {
			to = head
		}
		chunkFailed, derr := s.drainChunk(ctx, client, parser, from, to, addresses)
		if derr != nil {
			return derr
		}
		if chunkFailed {
			failedChunks++
		} else {
			highestClean = to
			if to >= persistAt {
				if err := s.store.UpdateChainCursor(ctx, to); err != nil {
					return fmt.Errorf("persist cursor at %d: %w", to, err)
				}
				persistAt = to + s.cursorEvery
			}
		}
		from = to + 1
	}

	finalCursor := head
	if failedChunks > 0 {
		finalCursor = highestClean
	}
	if finalCursor > cursor.LastProcessedBlock {
		if err := s.store.UpdateChainCursor(ctx, finalCursor); err != nil {
			return fmt.Errorf("persist final cursor at %d: %w", finalCursor, err)
		}
	}
	s.observeLag(float64(head - finalCursor))

	if failedChunks > 0 {
		logger.Log().Warnf("chainsub: catch-up reached %d but %d chunk(s) had per-address failures; the gap retries on the next connect", finalCursor, failedChunks)
	}
	return nil
}

// drainChunk pulls one eth_getLogs per address for [from, to] and
// ingests every log. Returns chunkFailed=true if any per-address call
// errored; a non-nil error only for ctx cancel/deadline.
func (s *Subscriber) drainChunk(ctx context.Context, client logFetcher, parser *Parser, from, to uint64, addresses []common.Address) (bool, error) {
	chunkFailed := false
	for _, addr := range addresses {
		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(to),
			Addresses: []common.Address{addr},
		}
		fetched, err := client.FilterLogs(ctx, query)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return chunkFailed, err
			}
			logger.Log().Warnf("chainsub: catch-up eth_getLogs [%d..%d] addr=%s: %v", from, to, addr.Hex(), err)
			chunkFailed = true
			if s.metrics != nil {
				s.metrics.ObserveDecodeError()
			}
			continue
		}
		for _, lg := range fetched {
			// allowRoundID=false during catch-up: the archival eth_call
			// per fulfilled event would hammer the RPC and slow the
			// replay; round_id is best-effort and backfilled live. rpc
			// is nil here precisely because it's never consulted.
			if perr := s.ingest(ctx, nil, parser, lg, false); perr != nil {
				if errors.Is(perr, context.Canceled) {
					return chunkFailed, perr
				}
				logger.Log().Warnf("chainsub: catch-up ingest [%d..%d] addr=%s: %v", from, to, addr.Hex(), perr)
			}
		}
	}
	return chunkFailed, nil
}

// ingest is the single path every log flows through — catch-up and live
// both call it. Parse → (optional best-effort round_id) → persist →
// publish-on-insert.
func (s *Subscriber) ingest(ctx context.Context, rpc *ethclient.Client, parser *Parser, lg types.Log, allowRoundID bool) error {
	evt, err := parser.Parse(lg)
	switch {
	case errors.Is(err, ErrUnknownEvent):
		return nil // a log we don't care about (e.g. Ownable*)
	case errors.Is(err, ErrAssetIDUnknown):
		logger.Log().Warnf("chainsub: aggregator %s emitted %s but is not in the registry mapping yet", lg.Address.Hex(), evt.Kind)
	case err != nil:
		if s.metrics != nil {
			s.metrics.ObserveDecodeError()
		}
		return fmt.Errorf("parse log: %w", err)
	}

	if allowRoundID && evt.Kind == models.EventKindPriceFulfilled && rpc != nil {
		rctx, cancel := context.WithTimeout(ctx, roundIDLimit)
		if roundID, rerr := s.fetchRoundID(rctx, rpc, lg.Address, lg.BlockNumber); rerr == nil {
			evt.PriceFulfilled.RoundID = roundID
		}
		cancel()
	}

	if evt.Kind == models.EventKindAssetRegistered && evt.AssetRegistered != nil {
		s.recordAggregator(evt.AssetRegistered.Aggregator, evt.AssetRegistered.AssetID)
		if perr := s.store.UpsertAggregator(ctx, evt.AssetRegistered.Aggregator, evt.AssetRegistered.AssetID); perr != nil {
			logger.Log().Warnf("chainsub: persist aggregator mapping: %v", perr)
		}
		// The new aggregator's own logs reach the filter on the next
		// reconnect — which, on a public RPC, is rarely more than a
		// couple minutes away (catch-up then backfills the gap).
	}

	inserted, err := s.store.InsertEvent(ctx, evt)
	if err != nil {
		return fmt.Errorf("persist event: %w", err)
	}
	if !inserted {
		return nil // idempotent replay — already have it
	}
	if s.metrics != nil {
		s.metrics.ObserveSeen(evt.Kind)
	}
	// Emit-on-ingest: no confirmation gate. Publish straight to live
	// subscribers (nil-safe).
	var delivered int
	if s.publisher != nil {
		delivered = s.publisher.Publish(evt)
	}
	logger.Log().Infof("chainsub: seen %s asset=%s block=%d log_index=%d tx=%s -> published to %d subscriber(s)",
		evt.Kind, shortHash(evt.AssetID.Hex()), evt.BlockNumber, evt.LogIndex, shortHash(evt.TxHash.Hex()), delivered)
	return nil
}

// refreshFromRegistry enumerates the OracleRegistry on chain and upserts
// every (asset_id, aggregator) pair it returns.
func (s *Subscriber) refreshFromRegistry(ctx context.Context, c *ethclient.Client) error {
	caller, err := oracleregistry.NewOracleRegistryCaller(s.registryAddr, c)
	if err != nil {
		return fmt.Errorf("bind registry caller: %w", err)
	}
	opts := &bind.CallOpts{Context: ctx}
	assets, err := caller.ListAssets(opts)
	if err != nil {
		return fmt.Errorf("listAssets: %w", err)
	}
	for _, raw := range assets {
		assetID := common.BytesToHash(raw[:])
		addr, err := caller.GetAggregator(opts, raw)
		if err != nil {
			return fmt.Errorf("getAggregator(%s): %w", assetID.Hex(), err)
		}
		s.recordAggregator(addr, assetID)
		if perr := s.store.UpsertAggregator(ctx, addr, assetID); perr != nil {
			return fmt.Errorf("persist aggregator %s: %w", addr.Hex(), perr)
		}
	}
	return nil
}

// fetchRoundID calls PriceAggregator.latestRoundData() at the fulfilling
// block. Best-effort: returns (nil, nil) on revert.
func (s *Subscriber) fetchRoundID(ctx context.Context, c *ethclient.Client, aggregator common.Address, blockNumber uint64) (*big.Int, error) {
	caller, err := priceaggregator.NewPriceAggregatorCaller(aggregator, c)
	if err != nil {
		return nil, fmt.Errorf("bind aggregator caller: %w", err)
	}
	rd, err := caller.LatestRoundData(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(blockNumber),
	})
	if err != nil {
		return nil, err
	}
	if rd.RoundId == nil {
		return nil, nil
	}
	return new(big.Int).Set(rd.RoundId), nil
}

func (s *Subscriber) observeLag(v float64) {
	if s.metrics != nil {
		s.metrics.ObserveLagBlocks(v)
	}
}

func (s *Subscriber) replaceMapping(m map[common.Address]common.Hash) {
	s.mu.Lock()
	s.mapping = make(map[common.Address]common.Hash, len(m))
	for k, v := range m {
		s.mapping[k] = v
	}
	s.mu.Unlock()
}

func (s *Subscriber) recordAggregator(addr common.Address, asset common.Hash) {
	s.mu.Lock()
	s.mapping[addr] = asset
	s.mu.Unlock()
}

func (s *Subscriber) subscribedAddresses() []common.Address {
	s.mu.RLock()
	out := make([]common.Address, 0, len(s.mapping)+1)
	out = append(out, s.registryAddr)
	for addr := range s.mapping {
		out = append(out, addr)
	}
	s.mu.RUnlock()
	return out
}
