// Package streamhub is the in-memory pub/sub at the heart of the
// indexer's StreamEvents server. It IS the confirmation gate — events
// only flow through the hub once the confirmer has marked them final
// (confirmations >= N, orphaned == false). Subscribers (oracle-service
// + rest-api) get a clean post-confirmation feed and never see
// in-flight or re-orged events.
//
// Slow consumers are dropped on overflow rather than allowed to block
// the publish path, so a stalled subscriber can't stall the rest of
// the system.
package streamhub

import (
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// Filter constrains which events a subscriber receives.
//
//   - Empty Kinds means "any kind".
//   - Nil AssetID means "any asset".
//
// A filter is evaluated for every event before fanout — non-matching
// events never touch a subscriber's channel.
type Filter struct {
	Kinds   []models.EventKind
	AssetID *common.Hash
}

// Matches reports whether the event satisfies the filter.
func (f Filter) Matches(e *models.Event) bool {
	if e == nil {
		return false
	}
	if len(f.Kinds) > 0 {
		ok := false
		for _, k := range f.Kinds {
			if k == e.Kind {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if f.AssetID != nil {
		if e.AssetID != *f.AssetID {
			return false
		}
	}
	return true
}

// DropReason describes why a subscriber was dropped. Surfaced through
// metrics so operators can tell slow-consumer events apart from
// shutdown events.
type DropReason string

const (
	// DropReasonSlow — the subscriber's buffer overflowed.
	DropReasonSlow DropReason = "slow"
	// DropReasonContext — the subscriber's context was canceled.
	DropReasonContext DropReason = "context_canceled"
)

// Subscription is the read-side of a stream-hub subscriber.
type Subscription struct {
	id     uint64
	hub    *Hub
	filter Filter
	ch     chan *models.Event
	once   sync.Once
}

// Events returns the receive channel. Closed when the subscription is
// canceled (either by the caller via Cancel or by the hub on slow-
// consumer drop / shutdown).
func (s *Subscription) Events() <-chan *models.Event { return s.ch }

// ID returns the subscription's identifier (useful for logs).
func (s *Subscription) ID() uint64 { return s.id }

// Cancel unsubscribes and releases resources. Idempotent.
func (s *Subscription) Cancel() {
	s.once.Do(func() {
		s.hub.remove(s.id, DropReasonContext)
	})
}

// DropFunc is invoked when the hub drops a subscriber. Optional —
// applications wire it up to metrics + logs.
type DropFunc func(id uint64, reason DropReason)

// Hub is the central fan-out. Goroutine-safe.
type Hub struct {
	mu       sync.RWMutex
	subs     map[uint64]*Subscription
	nextID   atomic.Uint64
	bufSize  int
	onDrop   DropFunc
	shutdown atomic.Bool
}

// New builds a Hub with the given per-subscriber buffer size. A
// buffer < 1 is clamped to 1 to keep the publish path non-blocking.
func New(bufSize int, onDrop DropFunc) *Hub {
	if bufSize < 1 {
		bufSize = 1
	}
	return &Hub{
		subs:    make(map[uint64]*Subscription),
		bufSize: bufSize,
		onDrop:  onDrop,
	}
}

// Subscribe registers a new consumer with the supplied filter and
// returns a Subscription. After Shutdown returns nil.
func (h *Hub) Subscribe(filter Filter) *Subscription {
	if h.shutdown.Load() {
		return nil
	}

	id := h.nextID.Add(1)
	sub := &Subscription{
		id:     id,
		hub:    h,
		filter: filter,
		ch:     make(chan *models.Event, h.bufSize),
	}

	h.mu.Lock()
	h.subs[id] = sub
	h.mu.Unlock()
	return sub
}

// Publish hands an event to every matching subscriber. Non-blocking:
// if a subscriber's buffer is full the subscriber is dropped (its
// channel is closed and it's removed from the hub) and Publish moves
// on. Returns the number of subscribers the event was successfully
// delivered to.
func (h *Hub) Publish(e *models.Event) int {
	if e == nil || h.shutdown.Load() {
		return 0
	}

	// Snapshot the subscribers under the read lock; deliver outside
	// the lock so a blocked Publish caller never blocks Subscribe /
	// Cancel.
	h.mu.RLock()
	snap := make([]*Subscription, 0, len(h.subs))
	for _, s := range h.subs {
		snap = append(snap, s)
	}
	h.mu.RUnlock()

	var dropped []uint64
	delivered := 0
	for _, s := range snap {
		if !s.filter.Matches(e) {
			continue
		}
		select {
		case s.ch <- e:
			delivered++
		default:
			dropped = append(dropped, s.id)
		}
	}

	for _, id := range dropped {
		h.remove(id, DropReasonSlow)
	}
	return delivered
}

// Subscribers returns the current subscriber count. Useful for
// metrics + tests.
func (h *Hub) Subscribers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.subs)
}

// Shutdown closes every subscriber channel and rejects new
// subscriptions / publishes. Idempotent.
func (h *Hub) Shutdown() {
	if !h.shutdown.CompareAndSwap(false, true) {
		return
	}
	h.mu.Lock()
	for id, s := range h.subs {
		close(s.ch)
		delete(h.subs, id)
	}
	h.mu.Unlock()
}

func (h *Hub) remove(id uint64, reason DropReason) {
	h.mu.Lock()
	sub, ok := h.subs[id]
	if ok {
		delete(h.subs, id)
	}
	h.mu.Unlock()
	if !ok {
		return
	}
	close(sub.ch)
	if h.onDrop != nil {
		h.onDrop(id, reason)
	}
}
