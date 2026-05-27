package streamhub

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

func eventOf(kind models.EventKind, asset common.Hash) *models.Event {
	return &models.Event{Kind: kind, AssetID: asset}
}

func TestFilterMatches(t *testing.T) {
	asset := common.HexToHash("0xaa")
	other := common.HexToHash("0xbb")

	cases := []struct {
		name   string
		filter Filter
		evt    *models.Event
		want   bool
	}{
		{"empty filter accepts anything", Filter{}, eventOf(models.EventKindPriceRequested, asset), true},
		{"empty filter rejects nil", Filter{}, nil, false},
		{"kind match", Filter{Kinds: []models.EventKind{models.EventKindPriceFulfilled}}, eventOf(models.EventKindPriceFulfilled, asset), true},
		{"kind mismatch", Filter{Kinds: []models.EventKind{models.EventKindPriceFulfilled}}, eventOf(models.EventKindPriceRequested, asset), false},
		{"asset match", Filter{AssetID: &asset}, eventOf(models.EventKindPriceRequested, asset), true},
		{"asset mismatch", Filter{AssetID: &asset}, eventOf(models.EventKindPriceRequested, other), false},
		{"both match", Filter{Kinds: []models.EventKind{models.EventKindPriceRequested}, AssetID: &asset}, eventOf(models.EventKindPriceRequested, asset), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.filter.Matches(c.evt); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestHubFanOutWithFilters(t *testing.T) {
	asset := common.HexToHash("0xaa")
	otherAsset := common.HexToHash("0xbb")

	hub := New(8, nil)
	defer hub.Shutdown()

	subAll := hub.Subscribe(Filter{})
	subRequested := hub.Subscribe(Filter{Kinds: []models.EventKind{models.EventKindPriceRequested}})
	subAsset := hub.Subscribe(Filter{AssetID: &asset})

	hub.Publish(eventOf(models.EventKindPriceRequested, asset))      // all three see this
	hub.Publish(eventOf(models.EventKindPriceFulfilled, asset))      // subAll + subAsset
	hub.Publish(eventOf(models.EventKindPriceFulfilled, otherAsset)) // subAll only

	drain := func(s *Subscription, want int) []*models.Event {
		t.Helper()
		got := make([]*models.Event, 0, want)
		timeout := time.After(200 * time.Millisecond)
		for len(got) < want {
			select {
			case e, ok := <-s.Events():
				if !ok {
					t.Fatalf("channel closed early after %d events", len(got))
				}
				got = append(got, e)
			case <-timeout:
				t.Fatalf("timeout: want %d events, got %d", want, len(got))
			}
		}
		return got
	}

	if got := drain(subAll, 3); len(got) != 3 {
		t.Errorf("subAll got %d events, want 3", len(got))
	}
	if got := drain(subRequested, 1); got[0].Kind != models.EventKindPriceRequested {
		t.Errorf("subRequested received wrong kind %v", got[0].Kind)
	}
	got := drain(subAsset, 2)
	for _, e := range got {
		if e.AssetID != asset {
			t.Errorf("subAsset received wrong asset %v", e.AssetID)
		}
	}
}

func TestHubSlowConsumerDropped(t *testing.T) {
	var (
		dropID     atomic.Uint64
		dropReason atomic.Value
		dropCh     = make(chan struct{}, 1)
	)

	hub := New(2, func(id uint64, reason DropReason) {
		dropID.Store(id)
		dropReason.Store(reason)
		select {
		case dropCh <- struct{}{}:
		default:
		}
	})
	defer hub.Shutdown()

	slow := hub.Subscribe(Filter{})

	// Push more than the buffer can hold — buffer is 2, push 5.
	asset := common.HexToHash("0x01")
	for i := 0; i < 5; i++ {
		hub.Publish(eventOf(models.EventKindPriceRequested, asset))
	}

	select {
	case <-dropCh:
	case <-time.After(time.Second):
		t.Fatal("expected slow subscriber to be dropped")
	}

	if id := dropID.Load(); id != slow.ID() {
		t.Errorf("dropped wrong subscriber: got id %d, want %d", id, slow.ID())
	}
	if r, _ := dropReason.Load().(DropReason); r != DropReasonSlow {
		t.Errorf("drop reason = %q, want %q", r, DropReasonSlow)
	}

	if hub.Subscribers() != 0 {
		t.Errorf("hub still has %d subscribers after slow-drop", hub.Subscribers())
	}
}

func TestHubFastConsumerSurvivesAlongsideSlow(t *testing.T) {
	hub := New(4, nil)
	defer hub.Shutdown()

	_ = hub.Subscribe(Filter{}) // slow — never drained
	fast := hub.Subscribe(Filter{})

	var received atomic.Int32
	done := make(chan struct{})
	go func() {
		defer close(done)
		for e := range fast.Events() {
			_ = e
			received.Add(1)
		}
	}()

	// Drive the publisher slowly enough that `fast` can drain
	// between publishes — `slow` will accumulate and eventually be
	// dropped, but `fast` should never miss anything.
	const total = 50
	asset := common.HexToHash("0x01")
	for i := 0; i < total; i++ {
		hub.Publish(eventOf(models.EventKindPriceRequested, asset))
		time.Sleep(200 * time.Microsecond)
	}

	hub.Shutdown()
	<-done
	if got := int(received.Load()); got != total {
		t.Errorf("fast subscriber received %d events, want %d", got, total)
	}
}

func TestHubCancelIsIdempotent(t *testing.T) {
	hub := New(1, nil)
	defer hub.Shutdown()

	sub := hub.Subscribe(Filter{})
	sub.Cancel()
	sub.Cancel()

	// Channel should be closed exactly once — a second close would panic.
	select {
	case _, ok := <-sub.Events():
		if ok {
			t.Error("expected closed channel")
		}
	default:
	}
}

func TestHubShutdownClosesAllSubscribers(t *testing.T) {
	hub := New(1, nil)
	subs := make([]*Subscription, 0, 4)
	for i := 0; i < 4; i++ {
		subs = append(subs, hub.Subscribe(Filter{}))
	}
	hub.Shutdown()

	var wg sync.WaitGroup
	wg.Add(len(subs))
	for _, s := range subs {
		go func(s *Subscription) {
			defer wg.Done()
			select {
			case _, ok := <-s.Events():
				if ok {
					t.Errorf("subscriber %d still receiving after shutdown", s.ID())
				}
			case <-time.After(time.Second):
				t.Errorf("subscriber %d never closed", s.ID())
			}
		}(s)
	}
	wg.Wait()

	if sub := hub.Subscribe(Filter{}); sub != nil {
		t.Error("Subscribe after Shutdown should return nil")
	}
}

func TestHubPublishNilSafe(t *testing.T) {
	hub := New(1, nil)
	defer hub.Shutdown()
	if got := hub.Publish(nil); got != 0 {
		t.Errorf("Publish(nil) returned %d, want 0", got)
	}
}
