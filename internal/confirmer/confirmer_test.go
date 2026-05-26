package confirmer

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// ---- test doubles ----------------------------------------------------

type fakeStore struct {
	mu       sync.Mutex
	pending  []*models.Event
	bumps    map[int64]uint32
	orphaned map[int64]bool
}

func newFakeStore(evts ...*models.Event) *fakeStore {
	return &fakeStore{
		pending:  evts,
		bumps:    make(map[int64]uint32),
		orphaned: make(map[int64]bool),
	}
}

func (f *fakeStore) PendingEvents(_ context.Context, _ uint32, _ int) ([]*models.Event, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := make([]*models.Event, len(f.pending))
	copy(cp, f.pending)
	return cp, nil
}

func (f *fakeStore) UpdateConfirmations(_ context.Context, id int64, confirmations uint32) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bumps[id] = confirmations
	for _, e := range f.pending {
		if e.ID == id {
			e.Confirmations = confirmations
		}
	}
	return nil
}

func (f *fakeStore) MarkOrphaned(_ context.Context, id int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.orphaned[id] = true
	for _, e := range f.pending {
		if e.ID == id {
			e.Orphaned = true
		}
	}
	return nil
}

type fakeChain struct {
	head    uint64
	headers map[uint64]*types.Header // canonical header per height
}

func (f *fakeChain) HeaderByNumber(_ context.Context, n *big.Int) (*types.Header, error) {
	target := f.head
	if n != nil {
		target = n.Uint64()
	}
	h, ok := f.headers[target]
	if !ok {
		return nil, errors.New("no header at block")
	}
	return h, nil
}

func (f *fakeChain) HeaderByHash(_ context.Context, h common.Hash) (*types.Header, error) {
	for _, hdr := range f.headers {
		if hdr.Hash() == h {
			return hdr, nil
		}
	}
	return nil, errors.New("not found")
}

// makeHeader returns a header with Number=n and Nonce=tag — Nonce
// alters the hash without changing semantics, so tests can construct
// canonical-vs-reorg pairs at the same height.
func makeHeader(n uint64, tag uint64) *types.Header {
	return &types.Header{Number: new(big.Int).SetUint64(n), Nonce: types.EncodeNonce(tag)}
}

type fakePublisher struct {
	mu       sync.Mutex
	received []*models.Event
}

func (f *fakePublisher) Publish(e *models.Event) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *e
	f.received = append(f.received, &cp)
	return 1
}

func (f *fakePublisher) snapshot() []*models.Event {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*models.Event, len(f.received))
	copy(out, f.received)
	return out
}

type fakeMetrics struct {
	confirmed, orphaned []models.EventKind
}

func (f *fakeMetrics) ObserveConfirmed(k models.EventKind) {
	f.confirmed = append(f.confirmed, k)
}
func (f *fakeMetrics) ObserveOrphaned(k models.EventKind) {
	f.orphaned = append(f.orphaned, k)
}
func (f *fakeMetrics) ObserveLagSeconds(_ float64) {}

// ---- helpers ---------------------------------------------------------

// mkEvent ties the event's block_hash to the canonical header at its
// height so the confirmer accepts it as in-chain.
func mkEvent(id int64, block uint64, kind models.EventKind, hash common.Hash) *models.Event {
	return &models.Event{
		ID:          id,
		Kind:        kind,
		BlockNumber: block,
		BlockHash:   hash,
	}
}

// ---- tests -----------------------------------------------------------

func TestProcess_FutureBlockIgnored(t *testing.T) {
	canonical100 := makeHeader(100, 0)
	chain := &fakeChain{head: 50, headers: map[uint64]*types.Header{50: makeHeader(50, 0)}}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, canonical100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}

	c := New(store, chain, pub, Config{Threshold: 5})
	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if _, ok := store.bumps[1]; ok {
		t.Error("future event should not have produced a confirmation bump")
	}
	if got := pub.snapshot(); len(got) != 0 {
		t.Errorf("future event must not have been published; got %d events", len(got))
	}
}

func TestProcess_BumpButNotYetFinal(t *testing.T) {
	canonical100 := makeHeader(100, 0)
	chain := &fakeChain{
		head: 102, // depth 2 < 5
		headers: map[uint64]*types.Header{
			102: makeHeader(102, 0),
			100: canonical100,
		},
	}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, canonical100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}
	c := New(store, chain, pub, Config{Threshold: 5})

	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if got := store.bumps[1]; got != 2 {
		t.Errorf("confirmations bumped to %d, want 2", got)
	}
	if got := pub.snapshot(); len(got) != 0 {
		t.Errorf("publisher saw %d events, want 0", len(got))
	}
}

func TestProcess_CrossesThresholdPublishes(t *testing.T) {
	canonical100 := makeHeader(100, 0)
	chain := &fakeChain{
		head: 110,
		headers: map[uint64]*types.Header{
			110: makeHeader(110, 0),
			100: canonical100,
		},
	}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, canonical100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}
	metrics := &fakeMetrics{}
	c := New(store, chain, pub, Config{Threshold: 5, Metrics: metrics})

	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if got := store.bumps[1]; got != 5 {
		t.Errorf("confirmations recorded = %d, want 5 (capped at threshold)", got)
	}
	if got := pub.snapshot(); len(got) != 1 || got[0].Confirmations != 5 {
		t.Errorf("publisher snapshot = %+v", got)
	}
	if len(metrics.confirmed) != 1 || metrics.confirmed[0] != models.EventKindPriceRequested {
		t.Errorf("metrics.confirmed = %v", metrics.confirmed)
	}
}

func TestProcess_ReorgMarksOrphanedAndNoPublish(t *testing.T) {
	// `recorded` is the header originally observed; canonical is a
	// different header at the same height (post-reorg).
	originalAt100 := makeHeader(100, 0)
	postReorgAt100 := makeHeader(100, 99) // different Nonce -> different hash
	chain := &fakeChain{
		head: 110,
		headers: map[uint64]*types.Header{
			110: makeHeader(110, 0),
			100: postReorgAt100, // the chain now reports this as canonical
		},
	}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, originalAt100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}
	metrics := &fakeMetrics{}
	c := New(store, chain, pub, Config{Threshold: 5, Metrics: metrics})

	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if !store.orphaned[1] {
		t.Error("event should have been marked orphaned")
	}
	if got := pub.snapshot(); len(got) != 0 {
		t.Errorf("orphaned event must not be published; got %d events", len(got))
	}
	if len(metrics.orphaned) != 1 || metrics.orphaned[0] != models.EventKindPriceRequested {
		t.Errorf("metrics.orphaned = %v", metrics.orphaned)
	}
}

func TestProcess_NoDoublePublishOnSecondTick(t *testing.T) {
	canonical100 := makeHeader(100, 0)
	chain := &fakeChain{
		head: 110,
		headers: map[uint64]*types.Header{
			110: makeHeader(110, 0),
			100: canonical100,
		},
	}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, canonical100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}
	c := New(store, chain, pub, Config{Threshold: 5})

	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("Tick: %v", err)
	}
	// On the second tick the pending query still returns the event
	// (in real usage PendingEvents would skip it because
	// confirmations >= threshold; the fake returns the whole slice).
	// The processor must NOT publish twice.
	if err := c.Tick(context.Background()); err != nil {
		t.Fatalf("second Tick: %v", err)
	}
	if got := pub.snapshot(); len(got) != 1 {
		t.Errorf("publisher saw %d events, want 1", len(got))
	}
}

func TestRun_ContextCancelExits(t *testing.T) {
	canonical100 := makeHeader(100, 0)
	chain := &fakeChain{
		head: 102,
		headers: map[uint64]*types.Header{
			102: makeHeader(102, 0),
			100: canonical100,
		},
	}
	evt := mkEvent(1, 100, models.EventKindPriceRequested, canonical100.Hash())
	store := newFakeStore(evt)
	pub := &fakePublisher{}
	c := New(store, chain, pub, Config{Threshold: 5, Interval: 50e6}) // 50ms

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.Run(ctx); err != nil {
		t.Errorf("Run returned %v, want nil after context cancel", err)
	}
}
