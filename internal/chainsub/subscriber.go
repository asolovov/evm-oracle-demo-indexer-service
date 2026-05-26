package chainsub

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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

// EventStore is the persistence surface the subscriber needs.
// Implemented by *repository.Repository in production.
type EventStore interface {
	InsertEvent(ctx context.Context, e *models.Event) (bool, error)
	UpsertAggregator(ctx context.Context, aggregator common.Address, assetID common.Hash) error
	AggregatorRegistry(ctx context.Context) (map[common.Address]common.Hash, error)
}

// Metrics is the optional metrics surface; nil-safe.
type Metrics interface {
	ObserveSeen(kind models.EventKind)
	ObserveDecodeError()
}

// Subscriber owns the WS subscription to the chain. It is a single-
// goroutine consumer: Run() blocks until ctx is canceled.
type Subscriber struct {
	wsURL          string
	rpcURL         string
	registryAddr   common.Address
	store          EventStore
	metrics        Metrics
	reconnectWait  time.Duration

	mu         sync.RWMutex
	mapping    map[common.Address]common.Hash // aggregator -> asset_id
	parser     *Parser
	client     *ethclient.Client
}

// Config bundles the constructor knobs.
type Config struct {
	WSURL           string
	RPCURL          string
	RegistryAddress common.Address
	ReconnectWait   time.Duration // default 2s
	Metrics         Metrics
}

// New builds a Subscriber. Does NOT dial; that happens on Run().
func New(store EventStore, cfg Config) *Subscriber {
	wait := cfg.ReconnectWait
	if wait <= 0 {
		wait = 2 * time.Second
	}
	return &Subscriber{
		wsURL:         cfg.WSURL,
		rpcURL:        cfg.RPCURL,
		registryAddr:  cfg.RegistryAddress,
		store:         store,
		metrics:       cfg.Metrics,
		reconnectWait: wait,
		mapping:       make(map[common.Address]common.Hash),
	}
}

// AssetIDFor implements AggregatorResolver — backed by the in-memory
// mapping the subscriber maintains.
func (s *Subscriber) AssetIDFor(addr common.Address) (common.Hash, bool) {
	s.mu.RLock()
	h, ok := s.mapping[addr]
	s.mu.RUnlock()
	return h, ok
}

// Client returns the underlying ethclient.Client for callers (e.g.
// the confirmer + backfill reconciler) that share the dial. Only
// valid after a successful Run() pass has connected.
func (s *Subscriber) Client() *ethclient.Client { return s.client }

// Run is a blocking loop: connect, subscribe, drain. On WS error,
// disconnect, or context cancellation, it reconnects after
// reconnectWait. Returns nil on ctx cancellation.
func (s *Subscriber) Run(ctx context.Context) error {
	for {
		if err := s.runOnce(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			logger.Log().Warnf("chainsub: session error: %v — reconnecting in %s", err, s.reconnectWait)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.reconnectWait):
		}
	}
}

func (s *Subscriber) runOnce(ctx context.Context) error {
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

	// Seed the aggregator->asset mapping: DB cache first, then refresh
	// from the on-chain registry (authoritative).
	cached, err := s.store.AggregatorRegistry(ctx)
	if err != nil {
		return fmt.Errorf("load cached registry: %w", err)
	}
	s.replaceMapping(cached)

	if err := s.refreshFromRegistry(ctx, rpcClient); err != nil {
		logger.Log().Warnf("chainsub: registry refresh failed: %v (continuing with cached mapping of size %d)", err, len(cached))
	}

	parser, err := NewParser(s.registryAddr, s)
	if err != nil {
		return fmt.Errorf("init parser: %w", err)
	}
	s.parser = parser
	s.client = wsClient

	// Build the address filter from the current mapping + the registry.
	addresses := s.subscribedAddresses()

	logger.Log().Infof("chainsub: subscribed to %d address(es) on %s", len(addresses), s.wsURL)

	logs := make(chan types.Log, 256)
	sub, err := wsClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
		Addresses: addresses,
	}, logs)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case log := <-logs:
			if err := s.handleLog(ctx, rpcClient, log); err != nil {
				logger.Log().Warnf("chainsub: handleLog: %v", err)
			}
		}
	}
}

func (s *Subscriber) handleLog(ctx context.Context, rpc *ethclient.Client, log types.Log) error {
	evt, err := s.parser.Parse(log)
	switch {
	case errors.Is(err, ErrUnknownEvent):
		return nil // not a log we care about (e.g. Ownable*)
	case errors.Is(err, ErrAssetIDUnknown):
		// Warning, not fatal — the parser still returned a usable event
		// (with empty asset_id) so we persist + log.
		logger.Log().Warnf("chainsub: aggregator %s emitted %s but is not in registry mapping",
			log.Address.Hex(), evt.Kind)
	case err != nil:
		if s.metrics != nil {
			s.metrics.ObserveDecodeError()
		}
		return fmt.Errorf("parse log (tx=%s idx=%d): %w", log.TxHash.Hex(), log.Index, err)
	}

	// For PriceFulfilled, best-effort backfill round_id at the
	// emitting block. Failure is non-fatal — we just leave round_id
	// nil and document it.
	if evt.Kind == models.EventKindPriceFulfilled && rpc != nil {
		if roundID, rerr := s.fetchRoundID(ctx, rpc, log.Address, log.BlockNumber); rerr == nil {
			evt.PriceFulfilled.RoundID = roundID
		}
	}

	// For AssetRegistered, update the live mapping + the persistent
	// cache so subsequent PriceRequested logs resolve correctly.
	if evt.Kind == models.EventKindAssetRegistered && evt.AssetRegistered != nil {
		s.recordAggregator(evt.AssetRegistered.Aggregator, evt.AssetRegistered.AssetID)
		if perr := s.store.UpsertAggregator(ctx, evt.AssetRegistered.Aggregator, evt.AssetRegistered.AssetID); perr != nil {
			logger.Log().Warnf("chainsub: persist aggregator mapping: %v", perr)
		}
		// NB: this NEW aggregator's logs won't reach the current
		// subscription until the next reconnect. v1 limitation —
		// documented in the README. We rely on operator-driven
		// restart for now.
	}

	inserted, err := s.store.InsertEvent(ctx, evt)
	if err != nil {
		return fmt.Errorf("persist event: %w", err)
	}
	if inserted && s.metrics != nil {
		s.metrics.ObserveSeen(evt.Kind)
	}
	return nil
}

// refreshFromRegistry enumerates the OracleRegistry on chain and
// upserts every (asset_id, aggregator) pair it returns.
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

// fetchRoundID calls PriceAggregator.latestRoundData() at the
// fulfilling block. Best-effort: returns (nil, nil) when the call
// reverts (e.g. the round is not yet readable at that block height).
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
