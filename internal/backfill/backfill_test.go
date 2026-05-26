package backfill

import (
	"context"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/chainsub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// ---- doubles ----------------------------------------------------------

type fakeFetcher struct {
	mu      sync.Mutex
	headNum uint64
	logs    []types.Log // returned slice is filtered by FromBlock/ToBlock at query time
	queries []ethereum.FilterQuery
}

func (f *fakeFetcher) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return &types.Header{Number: new(big.Int).SetUint64(f.headNum)}, nil
}

func (f *fakeFetcher) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queries = append(f.queries, q)

	out := make([]types.Log, 0)
	from := q.FromBlock.Uint64()
	to := q.ToBlock.Uint64()
	for _, l := range f.logs {
		if l.BlockNumber >= from && l.BlockNumber <= to {
			out = append(out, l)
		}
	}
	return out, nil
}

type fakeStore struct {
	mu              sync.Mutex
	inserts         []*models.Event
	cursor          uint64
	cursorWrites    []uint64
	aggregatorMap   map[common.Address]common.Hash
}

func (s *fakeStore) InsertEvent(_ context.Context, e *models.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, prev := range s.inserts {
		if prev.TxHash == e.TxHash && prev.LogIndex == e.LogIndex {
			return false, nil
		}
	}
	s.inserts = append(s.inserts, e)
	return true, nil
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
func (s *fakeStore) AggregatorRegistry(_ context.Context) (map[common.Address]common.Hash, error) {
	return s.aggregatorMap, nil
}

type stubResolver struct{ mapping map[common.Address]common.Hash }

func (s *stubResolver) AssetIDFor(addr common.Address) (common.Hash, bool) {
	h, ok := s.mapping[addr]
	return h, ok
}

// ---- helpers ----------------------------------------------------------

// buildPriceRequestedLog mirrors the helper in chainsub/parser_test.go but
// is duplicated here to keep test packages independent.
func buildPriceRequestedLog(aggregator common.Address, reqID *big.Int, requester common.Address, block uint64, txTag byte) types.Log {
	priceRequestedTopic := common.HexToHash("0x33b362cfc336cd3829a3a7896832b11a70bdbc0194a75fb7f1003919b23c8600")
	reqIDBytes := common.LeftPadBytes(reqID.Bytes(), 32)
	requesterBytes := common.LeftPadBytes(requester.Bytes(), 32)
	return types.Log{
		Address:     aggregator,
		Topics:      []common.Hash{priceRequestedTopic, common.BytesToHash(reqIDBytes), common.BytesToHash(requesterBytes)},
		Data:        nil,
		BlockNumber: block,
		BlockHash:   common.HexToHash("0xbeef"),
		TxHash:      common.BytesToHash(append([]byte{txTag}, make([]byte, 31)...)),
		Index:       0,
	}
}

// ---- tests ------------------------------------------------------------

func TestBackfill_NothingToReplay(t *testing.T) {
	fetcher := &fakeFetcher{headNum: 100}
	store := &fakeStore{cursor: 100, aggregatorMap: map[common.Address]common.Hash{}}
	parser, _ := chainsub.NewParser(common.HexToAddress("0x1"), &stubResolver{})
	r := New(fetcher, store, parser, Config{ChunkSize: 10})

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(fetcher.queries) != 0 {
		t.Errorf("expected no eth_getLogs calls, got %d", len(fetcher.queries))
	}
}

func TestBackfill_ChunkingAndCursorAdvance(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")

	var logs []types.Log
	for i, blk := range []uint64{5, 17, 28, 42, 51, 66, 79} {
		logs = append(logs, buildPriceRequestedLog(aggregator, big.NewInt(int64(i+1)), requester, blk, byte(i+1)))
	}

	fetcher := &fakeFetcher{headNum: 80, logs: logs}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{aggregator: asset}}
	parser, _ := chainsub.NewParser(common.HexToAddress("0x1"), &stubResolver{mapping: map[common.Address]common.Hash{aggregator: asset}})

	r := New(fetcher, store, parser, Config{ChunkSize: 25, CursorEvery: 50})

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(store.inserts) != len(logs) {
		t.Errorf("inserted %d events, want %d", len(store.inserts), len(logs))
	}
	// Chunking: with chunk size 25 and gap [1..80] we expect chunks
	// [1..25], [26..50], [51..75], [76..80] -> 4 queries.
	if len(fetcher.queries) != 4 {
		t.Errorf("expected 4 eth_getLogs chunks, got %d (queries=%+v)", len(fetcher.queries), fetcher.queries)
	}
	// Final cursor must be the head height.
	if store.cursor != 80 {
		t.Errorf("final cursor = %d, want 80", store.cursor)
	}
	// Cursor should have been persisted at least once mid-flight.
	if len(store.cursorWrites) < 2 {
		t.Errorf("expected at least 2 cursor writes (mid + final), got %d", len(store.cursorWrites))
	}
}

func TestBackfill_DefaultStartUsedWhenCursorZero(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")

	logs := []types.Log{
		buildPriceRequestedLog(aggregator, big.NewInt(1), requester, 1_000, 1),
		buildPriceRequestedLog(aggregator, big.NewInt(2), requester, 1_500, 2),
		buildPriceRequestedLog(aggregator, big.NewInt(3), requester, 999, 3),
	}

	fetcher := &fakeFetcher{headNum: 2_000, logs: logs}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{aggregator: asset}}
	parser, _ := chainsub.NewParser(common.HexToAddress("0x1"), &stubResolver{mapping: map[common.Address]common.Hash{aggregator: asset}})

	r := New(fetcher, store, parser, Config{ChunkSize: 5_000, DefaultStart: 1_000})

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// The log at block 999 must be ignored because DefaultStart=1000.
	if len(store.inserts) != 2 {
		t.Errorf("expected 2 inserts (block 1000 + 1500), got %d", len(store.inserts))
	}
	gotBlocks := []uint64{store.inserts[0].BlockNumber, store.inserts[1].BlockNumber}
	sort.Slice(gotBlocks, func(i, j int) bool { return gotBlocks[i] < gotBlocks[j] })
	if gotBlocks[0] != 1_000 || gotBlocks[1] != 1_500 {
		t.Errorf("unexpected inserted blocks: %v", gotBlocks)
	}
}

func TestBackfill_ContextCancelStops(t *testing.T) {
	fetcher := &fakeFetcher{headNum: 10}
	store := &fakeStore{aggregatorMap: map[common.Address]common.Hash{}}
	parser, _ := chainsub.NewParser(common.HexToAddress("0x1"), &stubResolver{})
	r := New(fetcher, store, parser, Config{ChunkSize: 1})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.Run(ctx)
	if err == nil || err != context.Canceled {
		t.Errorf("Run on cancelled ctx returned %v, want context.Canceled", err)
	}
}
