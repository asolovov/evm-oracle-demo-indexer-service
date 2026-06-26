package chainsub

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// ---- doubles ----------------------------------------------------------

type fakeStore struct {
	mu            sync.Mutex
	inserts       []*models.Event
	cursor        uint64
	cursorWrites  []uint64
	aggregatorMap map[common.Address]common.Hash
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

func (s *fakeStore) UpsertAggregator(_ context.Context, _ common.Address, _ common.Hash) error {
	return nil
}

func (s *fakeStore) AggregatorRegistry(_ context.Context) (map[common.Address]common.Hash, error) {
	return s.aggregatorMap, nil
}

func (s *fakeStore) ChainCursor(_ context.Context) (*models.ChainCursor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &models.ChainCursor{LastProcessedBlock: s.cursor, UpdatedAt: time.Now()}, nil
}

func (s *fakeStore) UpdateChainCursor(_ context.Context, block uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cursor = block
	s.cursorWrites = append(s.cursorWrites, block)
	return nil
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

type fakeFetcher struct {
	mu            sync.Mutex
	logs          []types.Log
	queries       []ethereum.FilterQuery
	failOnAddress *common.Address
}

func (f *fakeFetcher) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queries = append(f.queries, q)
	if f.failOnAddress != nil && len(q.Addresses) == 1 && q.Addresses[0] == *f.failOnAddress {
		return nil, errors.New("simulated provider rejection")
	}
	from, to := q.FromBlock.Uint64(), q.ToBlock.Uint64()
	addrs := make(map[common.Address]struct{}, len(q.Addresses))
	for _, a := range q.Addresses {
		addrs[a] = struct{}{}
	}
	out := make([]types.Log, 0)
	for _, l := range f.logs {
		if l.BlockNumber < from || l.BlockNumber > to {
			continue
		}
		if _, ok := addrs[l.Address]; !ok {
			continue
		}
		out = append(out, l)
	}
	return out, nil
}

func newSub(store EventStore, pub Publisher, agg common.Address, asset common.Hash) *Subscriber {
	s := New(store, pub, Config{
		RegistryAddress: common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef"),
		ChunkSize:       25,
		CursorEvery:     50,
	})
	s.recordAggregator(agg, asset)
	return s
}

// reqLogAt builds a PriceRequested log with a DISTINCT (tx_hash,
// log_index) per tag so idempotent inserts don't collapse them (the
// shared parser_test helper hardcodes one tx_hash).
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

func mustParser(t *testing.T, s *Subscriber) *Parser {
	t.Helper()
	p, err := NewParser(s.registryAddr, s)
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	return p
}

// ---- backoff ----------------------------------------------------------

func TestBackoffFor_BoundedAndJittered(t *testing.T) {
	// Every value must be within [0, cap] where cap grows exponentially
	// but never exceeds maxBackoff.
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
	// High attempts are capped at maxBackoff.
	if got := minBackoff << min(99, 16); got < maxBackoff {
		// sanity: the shift saturates well above the cap
		t.Logf("shifted floor %s", got)
	}
}

// ---- catch-up ---------------------------------------------------------

func TestCatchUp_PublishesAndAdvancesCursor(t *testing.T) {
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")

	blocks := []uint64{5, 17, 28, 51, 79}
	logs := make([]types.Log, 0, len(blocks))
	for i, blk := range blocks {
		logs = append(logs, reqLogAt(agg, big.NewInt(int64(i+1)), requester, blk, byte(i+1)))
	}
	fetcher := &fakeFetcher{logs: logs}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{agg: asset}}
	pub := &fakePublisher{}
	s := newSub(store, pub, agg, asset)

	if err := s.catchUp(context.Background(), fetcher, mustParser(t, s), s.subscribedAddresses(), 80); err != nil {
		t.Fatalf("catchUp: %v", err)
	}

	// Every event persisted AND published (emit-on-ingest, no gate).
	if len(store.inserts) != len(logs) {
		t.Errorf("inserted %d, want %d", len(store.inserts), len(logs))
	}
	if pub.count() != len(logs) {
		t.Errorf("published %d, want %d (catch-up must publish, not just persist)", pub.count(), len(logs))
	}
	// Cursor advanced to head.
	if store.cursor != 80 {
		t.Errorf("final cursor = %d, want 80", store.cursor)
	}
	// Per-address chunking: every query targets exactly one address.
	for i, q := range fetcher.queries {
		if len(q.Addresses) != 1 {
			t.Errorf("query[%d] targets %d addresses, want 1", i, len(q.Addresses))
		}
	}
}

func TestCatchUp_PartialFailureParksCursor(t *testing.T) {
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	registry := common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef")

	// Registry filter rejected on every chunk → no chunk is fully clean.
	fetcher := &fakeFetcher{failOnAddress: &registry}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{agg: asset}}
	pub := &fakePublisher{}
	s := newSub(store, pub, agg, asset)

	if err := s.catchUp(context.Background(), fetcher, mustParser(t, s), s.subscribedAddresses(), 80); err != nil {
		t.Fatalf("catchUp: %v", err)
	}
	if store.cursor != 0 {
		t.Errorf("cursor advanced to %d despite every chunk failing; want 0", store.cursor)
	}
}

func TestCatchUp_NothingToDo(t *testing.T) {
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	fetcher := &fakeFetcher{}
	store := &fakeStore{cursor: 100, aggregatorMap: map[common.Address]common.Hash{agg: asset}}
	s := newSub(store, &fakePublisher{}, agg, asset)

	if err := s.catchUp(context.Background(), fetcher, mustParser(t, s), s.subscribedAddresses(), 100); err != nil {
		t.Fatalf("catchUp: %v", err)
	}
	if len(fetcher.queries) != 0 {
		t.Errorf("expected no getLogs calls when cursor==head, got %d", len(fetcher.queries))
	}
}

// ---- ingest publish-on-insert + idempotency ---------------------------

func TestIngest_PublishesOncePerNewEvent(t *testing.T) {
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")

	store := &fakeStore{aggregatorMap: map[common.Address]common.Hash{agg: asset}}
	pub := &fakePublisher{}
	s := newSub(store, pub, agg, asset)
	parser := mustParser(t, s)
	lg := buildPriceRequestedLog(agg, big.NewInt(7), requester, 42)

	// First ingest: persisted + published.
	if err := s.ingest(context.Background(), nil, parser, lg, false); err != nil {
		t.Fatalf("ingest 1: %v", err)
	}
	// Replay of the same log: idempotent, NOT re-published.
	if err := s.ingest(context.Background(), nil, parser, lg, false); err != nil {
		t.Fatalf("ingest 2: %v", err)
	}
	if pub.count() != 1 {
		t.Errorf("published %d times, want exactly 1 (replays must not re-publish)", pub.count())
	}
	if len(store.inserts) != 1 {
		t.Errorf("stored %d, want 1", len(store.inserts))
	}
}
