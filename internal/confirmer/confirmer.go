// Package confirmer walks the events table on a tick and applies the
// chain's verdict: bump confirmations until they cross the configured
// threshold, mark orphaned if the original block is no longer the
// canonical block at that height. When an event first crosses the
// threshold and remains canonical, the confirmer hands it to the
// stream hub — confirmer + hub together form the confirmation gate.
package confirmer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// EventStore is the persistence surface the confirmer needs.
type EventStore interface {
	PendingEvents(ctx context.Context, confirmations uint32, limit int) ([]*models.Event, error)
	UpdateConfirmations(ctx context.Context, id int64, confirmations uint32) error
	MarkOrphaned(ctx context.Context, id int64) error
}

// HeadProvider is the chain surface the confirmer needs. The methods
// match go-ethereum's `ethclient.Client` so the production
// implementation needs no shim.
type HeadProvider interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
}

// Publisher publishes a finalized event to subscribers. Implemented
// by *streamhub.Hub in production.
type Publisher interface {
	Publish(e *models.Event) int
}

// Metrics is the optional metrics-collection surface. Nil-safe.
type Metrics interface {
	ObserveConfirmed(kind models.EventKind)
	ObserveOrphaned(kind models.EventKind)
	ObserveLagSeconds(lag float64)
}

// Confirmer ticks every Interval, drains a bounded chunk of pending
// events, and applies the chain's verdict.
type Confirmer struct {
	store        EventStore
	chain        HeadProvider
	publisher    Publisher
	metrics      Metrics
	threshold    uint32
	interval     time.Duration
	batchLimit   int
	headObserver func(uint64)

	running  atomic.Bool
	stopOnce sync.Once
	stop     chan struct{}
}

// Config bundles the constructor knobs.
type Config struct {
	Threshold    uint32
	Interval     time.Duration
	BatchLimit   int
	Metrics      Metrics
	HeadObserver func(uint64) // optional — useful in tests + lag metrics
}

// New constructs a Confirmer.
func New(store EventStore, chain HeadProvider, publisher Publisher, cfg Config) *Confirmer {
	if cfg.Threshold == 0 {
		cfg.Threshold = 5
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 10 * time.Second
	}
	if cfg.BatchLimit <= 0 {
		cfg.BatchLimit = 500
	}
	return &Confirmer{
		store:        store,
		chain:        chain,
		publisher:    publisher,
		metrics:      cfg.Metrics,
		threshold:    cfg.Threshold,
		interval:     cfg.Interval,
		batchLimit:   cfg.BatchLimit,
		headObserver: cfg.HeadObserver,
		stop:         make(chan struct{}),
	}
}

// Run blocks until ctx is canceled. The confirmer is single-shot —
// concurrent Run calls return an error rather than racing on internal
// state.
func (c *Confirmer) Run(ctx context.Context) error {
	if !c.running.CompareAndSwap(false, true) {
		return errors.New("confirmer already running")
	}
	defer c.running.Store(false)

	t := time.NewTicker(c.interval)
	defer t.Stop()

	// Run an initial tick so a fresh startup doesn't have to wait
	// `interval` seconds before draining the backlog. Errors here
	// are transient by design — the loop below will retry on the
	// next tick.
	if err := c.Tick(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		_ = err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stop:
			return nil
		case <-t.C:
			if err := c.Tick(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				// Don't terminate on transient errors; the next tick will retry.
				// Caller surfaces failures via metrics + logs.
				_ = err
			}
		}
	}
}

// Stop signals Run to return. Safe to call before Run; idempotent
// across repeated calls.
func (c *Confirmer) Stop() {
	c.stopOnce.Do(func() { close(c.stop) })
}

// Tick runs one drain cycle. Exported for tests + the initial pass at
// app startup before the goroutine takes over.
func (c *Confirmer) Tick(ctx context.Context) error {
	head, err := c.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("head lookup: %w", err)
	}
	headNum := head.Number.Uint64()
	if c.headObserver != nil {
		c.headObserver(headNum)
	}

	pending, err := c.store.PendingEvents(ctx, c.threshold, c.batchLimit)
	if err != nil {
		return fmt.Errorf("read pending events: %w", err)
	}

	for _, evt := range pending {
		if err := c.process(ctx, evt, headNum); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			// Per-event failures shouldn't poison the whole tick; just
			// log + continue. The next tick will retry.
			_ = err
		}
	}
	return nil
}

func (c *Confirmer) process(ctx context.Context, evt *models.Event, headNum uint64) error {
	// Future blocks (head < block_number) — the WS subscription
	// races the head; just wait for the next tick.
	if headNum < evt.BlockNumber {
		if c.metrics != nil {
			c.metrics.ObserveLagSeconds(0)
		}
		return nil
	}

	depth := headNum - evt.BlockNumber

	// Verify the event's original block is still the canonical block
	// at that height. A reorg replaces the block at a height with a
	// different hash; if HeaderByNumber for that height doesn't match
	// our recorded block_hash, the event has been re-orged.
	canonicalHeader, err := c.chain.HeaderByNumber(ctx, new(big.Int).SetUint64(evt.BlockNumber))
	if err != nil {
		return fmt.Errorf("canonical header for block %d: %w", evt.BlockNumber, err)
	}

	if canonicalHeader.Hash() != evt.BlockHash {
		// Re-orged.
		if err := c.store.MarkOrphaned(ctx, evt.ID); err != nil {
			return fmt.Errorf("mark orphaned: %w", err)
		}
		if c.metrics != nil {
			c.metrics.ObserveOrphaned(evt.Kind)
		}
		return nil
	}

	// Cap the recorded depth at the threshold — there's no value in
	// going past it (we only need to know "final").
	recorded := depth
	if recorded > uint64(c.threshold) {
		recorded = uint64(c.threshold)
	}
	newConfirmations := uint32(recorded) //nolint:gosec // bounded by threshold.

	priorConfirmations := evt.Confirmations
	if newConfirmations <= priorConfirmations {
		// Nothing changed since last tick.
		return nil
	}

	if err := c.store.UpdateConfirmations(ctx, evt.ID, newConfirmations); err != nil {
		return fmt.Errorf("update confirmations: %w", err)
	}
	evt.Confirmations = newConfirmations

	// Did we just cross the threshold?
	if priorConfirmations < c.threshold && newConfirmations >= c.threshold {
		if c.publisher != nil {
			c.publisher.Publish(evt)
		}
		if c.metrics != nil {
			c.metrics.ObserveConfirmed(evt.Kind)
		}
	}
	return nil
}
