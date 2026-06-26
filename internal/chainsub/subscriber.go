package chainsub

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
)

// Backoff bounds for the reconnect loop. Public-RPC WS resets are the
// steady state; these are operational constants, not config.
const (
	minBackoff   = time.Second
	maxBackoff   = 30 * time.Second
	healthyReset = 60 * time.Second // a session living this long resets backoff
)

// EventStore is the persistence surface the subscriber needs.
type EventStore interface {
	InsertEvent(ctx context.Context, e *models.Event) (bool, error)
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
}

// AssetMapping is one seeded aggregator->asset entry.
type AssetMapping struct {
	Aggregator common.Address
	AssetID    common.Hash
}

// Subscriber is the single chain-observer. One goroutine (Run) owns the
// WS connection end-to-end: it dials, subscribes to the registry + the
// configured aggregators, and publishes each observed log the moment it
// is persisted (emit-on-ingest, no confirmation gate). On disconnect it
// reconnects with exponential backoff + jitter.
//
// LIVE-ONLY by design: no historical eth_getLogs, no eth_call, no
// cursor — the only chain operation is the WS log subscription, which
// every free RPC tier supports. The aggregator->asset mapping is
// seeded from config (the indexer never enumerates the registry on
// chain); live AssetRegistered events extend it at runtime.
//
// Because the client never leaves this goroutine there is no shared
// mutable client state, hence no data race.
type Subscriber struct {
	wsURL        string
	registryAddr common.Address
	store        EventStore
	publisher    Publisher
	metrics      Metrics

	mu      sync.RWMutex
	mapping map[common.Address]common.Hash // aggregator -> asset_id
}

// Config bundles the constructor knobs.
type Config struct {
	WSURL           string
	RegistryAddress common.Address
	Assets          []AssetMapping // seeded aggregator->asset mapping
	Metrics         Metrics
}

// New builds a Subscriber, seeding the aggregator->asset mapping from
// config. It does NOT dial; Run() owns the connection.
func New(store EventStore, publisher Publisher, cfg Config) *Subscriber {
	mapping := make(map[common.Address]common.Hash, len(cfg.Assets))
	for _, a := range cfg.Assets {
		mapping[a.Aggregator] = a.AssetID
	}
	return &Subscriber{
		wsURL:        cfg.WSURL,
		registryAddr: cfg.RegistryAddress,
		store:        store,
		publisher:    publisher,
		metrics:      cfg.Metrics,
		mapping:      mapping,
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
// used entirely within this call stack — no struct field, no sharing.
func (s *Subscriber) runSession(ctx context.Context) error {
	client, err := ethclient.DialContext(ctx, s.wsURL)
	if err != nil {
		return fmt.Errorf("dial ws: %w", err)
	}
	defer client.Close()

	parser, err := NewParser(s.registryAddr, s)
	if err != nil {
		return fmt.Errorf("init parser: %w", err)
	}
	addresses := s.subscribedAddresses()

	logs := make(chan types.Log, 1024)
	sub, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{Addresses: addresses}, logs)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	logger.Log().Infof("chainsub: live on %s — watching %d address(es)", s.wsURL, len(addresses))

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
			if err := s.ingest(ctx, parser, lg); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				logger.Log().Warnf("chainsub: ingest log (tx=%s idx=%d): %v", lg.TxHash.Hex(), lg.Index, err)
			}
		}
	}
}

// ingest is the single path every log flows through: parse → persist →
// publish-on-insert.
func (s *Subscriber) ingest(ctx context.Context, parser *Parser, lg types.Log) error {
	evt, err := parser.Parse(lg)
	switch {
	case errors.Is(err, ErrUnknownEvent):
		return nil // a log we don't care about (e.g. Ownable*)
	case errors.Is(err, ErrAssetIDUnknown):
		logger.Log().Warnf("chainsub: aggregator %s emitted %s but is not in the asset mapping", lg.Address.Hex(), evt.Kind)
	case err != nil:
		if s.metrics != nil {
			s.metrics.ObserveDecodeError()
		}
		return fmt.Errorf("parse log: %w", err)
	}

	// A live AssetRegistered extends the in-memory mapping so this
	// aggregator's price logs resolve their asset_id on the next
	// reconnect (when subscribedAddresses picks it up).
	if evt.Kind == models.EventKindAssetRegistered && evt.AssetRegistered != nil {
		s.recordAggregator(evt.AssetRegistered.Aggregator, evt.AssetRegistered.AssetID)
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
	var delivered int
	if s.publisher != nil {
		delivered = s.publisher.Publish(evt)
	}
	logger.Log().Infof("chainsub: seen %s asset=%s block=%d log_index=%d tx=%s -> published to %d subscriber(s)",
		evt.Kind, shortHash(evt.AssetID.Hex()), evt.BlockNumber, evt.LogIndex, shortHash(evt.TxHash.Hex()), delivered)
	return nil
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
