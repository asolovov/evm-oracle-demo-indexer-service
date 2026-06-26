package chainsub

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// ---- doubles ----------------------------------------------------------

type fakeStore struct {
	mu      sync.Mutex
	inserts []*models.Event
}

func (s *fakeStore) InsertEvent(_ context.Context, e *models.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, prev := range s.inserts {
		if prev.TxHash == e.TxHash && prev.LogIndex == e.LogIndex {
			return false, nil // idempotent replay
		}
	}
	s.inserts = append(s.inserts, e)
	return true, nil
}

type fakePublisher struct {
	mu       sync.Mutex
	received []*models.Event
}

func (p *fakePublisher) Publish(e *models.Event) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.received = append(p.received, e)
	return 1
}

func (p *fakePublisher) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.received)
}

var (
	testRegistry   = common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef")
	testAggregator = common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	testAsset      = common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
)

func newSub(store EventStore, pub Publisher) *Subscriber {
	return New(store, pub, Config{
		WSURL:           "ws://unused-in-unit-tests",
		RegistryAddress: testRegistry,
		Assets:          []AssetMapping{{Aggregator: testAggregator, AssetID: testAsset}},
	})
}

// reqLogAt builds a PriceRequested log with a distinct (tx_hash,
// log_index) per tag.
func reqLogAt(agg common.Address, reqID *big.Int, requester common.Address, block uint64, tag byte) types.Log {
	reqIDBytes := common.LeftPadBytes(reqID.Bytes(), 32)
	requesterBytes := common.LeftPadBytes(requester.Bytes(), 32)
	return types.Log{
		Address:     agg,
		Topics:      []common.Hash{priceRequestedTopic, common.BytesToHash(reqIDBytes), common.BytesToHash(requesterBytes)},
		BlockNumber: block,
		BlockHash:   common.HexToHash("0xbeef"),
		TxHash:      common.BytesToHash(append([]byte{tag}, make([]byte, 31)...)),
		Index:       uint(tag),
	}
}

// ---- tests -----------------------------------------------------------

func TestBackoffFor_BoundedAndJittered(t *testing.T) {
	for attempt := 0; attempt < 20; attempt++ {
		capDelay := minBackoff << min(attempt, 16)
		if capDelay <= 0 || capDelay > maxBackoff {
			capDelay = maxBackoff
		}
		for i := 0; i < 50; i++ {
			d := backoffFor(attempt)
			if d < 0 || d > capDelay {
				t.Fatalf("attempt %d: delay %s out of [0,%s]", attempt, d, capDelay)
			}
		}
	}
}

func TestNew_SeedsMappingFromConfig(t *testing.T) {
	s := newSub(&fakeStore{}, &fakePublisher{})
	got, ok := s.AssetIDFor(testAggregator)
	if !ok || got != testAsset {
		t.Errorf("AssetIDFor(seeded aggregator) = %s,%v; want %s,true", got.Hex(), ok, testAsset.Hex())
	}
	if _, ok := s.AssetIDFor(common.HexToAddress("0xdead")); ok {
		t.Error("AssetIDFor(unknown) should be false")
	}
}

func TestSubscribedAddresses_RegistryPlusAggregators(t *testing.T) {
	s := newSub(&fakeStore{}, &fakePublisher{})
	addrs := s.subscribedAddresses()
	if len(addrs) != 2 {
		t.Fatalf("want 2 addresses (registry + 1 aggregator), got %d", len(addrs))
	}
	set := map[common.Address]bool{addrs[0]: true, addrs[1]: true}
	if !set[testRegistry] || !set[testAggregator] {
		t.Errorf("addresses missing registry or aggregator: %v", addrs)
	}
}

func TestIngest_PublishesOncePerNewEvent(t *testing.T) {
	store := &fakeStore{}
	pub := &fakePublisher{}
	s := newSub(store, pub)
	parser, err := NewParser(s.registryAddr, s)
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")
	lg := reqLogAt(testAggregator, big.NewInt(7), requester, 42, 1)

	// First ingest: persisted + published.
	if err := s.ingest(context.Background(), parser, lg); err != nil {
		t.Fatalf("ingest 1: %v", err)
	}
	// Replay: idempotent, not re-published.
	if err := s.ingest(context.Background(), parser, lg); err != nil {
		t.Fatalf("ingest 2: %v", err)
	}
	if pub.count() != 1 {
		t.Errorf("published %d times, want exactly 1 (replays must not re-publish)", pub.count())
	}
	if len(store.inserts) != 1 {
		t.Errorf("stored %d, want 1", len(store.inserts))
	}
	if store.inserts[0].AssetID != testAsset {
		t.Errorf("asset_id not resolved from config mapping: %s", store.inserts[0].AssetID.Hex())
	}
}

func TestIngest_AssetRegisteredExtendsMapping(t *testing.T) {
	store := &fakeStore{}
	s := newSub(store, &fakePublisher{})
	parser, _ := NewParser(s.registryAddr, s)

	newAgg := common.HexToAddress("0x1111111111111111111111111111111111111111")
	newAsset := common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")
	// AssetRegistered(assetId indexed, aggregator indexed) from the registry.
	lg := types.Log{
		Address: testRegistry,
		Topics: []common.Hash{
			assetRegisteredTopic,
			newAsset,
			common.BytesToHash(common.LeftPadBytes(newAgg.Bytes(), 32)),
		},
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xfeed"),
		Index:       0,
	}
	if err := s.ingest(context.Background(), parser, lg); err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if got, ok := s.AssetIDFor(newAgg); !ok || got != newAsset {
		t.Errorf("live AssetRegistered did not extend the mapping: %s,%v", got.Hex(), ok)
	}
}
