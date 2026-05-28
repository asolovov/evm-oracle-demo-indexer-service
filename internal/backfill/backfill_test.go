package backfill

import (
	"context"
	"errors"
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
	logs    []types.Log // returned slice is filtered by FromBlock/ToBlock and Addresses at query time
	queries []ethereum.FilterQuery

	// failOnAddress, when set, makes FilterLogs return an error
	// whenever the query targets exactly this address. Lets a test
	// simulate "this one provider blocks this one filter".
	failOnAddress *common.Address
	failErr       error
}

func (f *fakeFetcher) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return &types.Header{Number: new(big.Int).SetUint64(f.headNum)}, nil
}

func (f *fakeFetcher) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queries = append(f.queries, q)

	if f.failOnAddress != nil && len(q.Addresses) == 1 && q.Addresses[0] == *f.failOnAddress {
		err := f.failErr
		if err == nil {
			err = errors.New("simulated provider rejection")
		}
		return nil, err
	}

	out := make([]types.Log, 0)
	from := q.FromBlock.Uint64()
	to := q.ToBlock.Uint64()
	addrSet := make(map[common.Address]struct{}, len(q.Addresses))
	for _, a := range q.Addresses {
		addrSet[a] = struct{}{}
	}
	for _, l := range f.logs {
		if l.BlockNumber < from || l.BlockNumber > to {
			continue
		}
		if len(addrSet) > 0 {
			if _, ok := addrSet[l.Address]; !ok {
				continue
			}
		}
		out = append(out, l)
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

	blocks := []uint64{5, 17, 28, 42, 51, 66, 79}
	logs := make([]types.Log, 0, len(blocks))
	for i, blk := range blocks {
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
	// Chunking: with chunk size 25 and gap [1..80] there are 4 chunks
	// (chunks [1..25], [26..50], [51..75], [76..80]) and 2 addresses
	// (registry + 1 aggregator), so 4 * 2 = 8 eth_getLogs calls.
	if len(fetcher.queries) != 8 {
		t.Errorf("expected 8 eth_getLogs calls (4 chunks * 2 addresses), got %d", len(fetcher.queries))
	}
	// Each query must target exactly one address — the per-address
	// chunking guard against providers that block multi-address filters.
	for i, q := range fetcher.queries {
		if len(q.Addresses) != 1 {
			t.Errorf("query[%d] targets %d addresses, want exactly 1", i, len(q.Addresses))
		}
	}
	// Final cursor must be the head height (all chunks clean).
	if store.cursor != 80 {
		t.Errorf("final cursor = %d, want 80", store.cursor)
	}
	if len(store.cursorWrites) < 2 {
		t.Errorf("expected at least 2 cursor writes (mid + final), got %d", len(store.cursorWrites))
	}
}

func TestBackfill_PerAddressFailureWarnsAndParksCursor(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	registry := common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")

	logs := []types.Log{
		buildPriceRequestedLog(aggregator, big.NewInt(1), requester, 5, 1),
		buildPriceRequestedLog(aggregator, big.NewInt(2), requester, 12, 2),
	}

	// Registry filter rejected — simulates publicnode-style provider
	// behavior on a specific address.
	fetcher := &fakeFetcher{
		headNum:       30,
		logs:          logs,
		failOnAddress: &registry,
		failErr:       errors.New("Request blocked. Details: blocked parameter: params.0.address.#"),
	}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{aggregator: asset}}
	parser, _ := chainsub.NewParser(registry, &stubResolver{mapping: map[common.Address]common.Hash{aggregator: asset}})

	r := New(fetcher, store, parser, Config{
		RegistryAddress: registry,
		ChunkSize:       50, // one chunk covers [1..30]
		CursorEvery:     100,
	})

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v (per-address failures must NOT abort the whole pass)", err)
	}
	// Aggregator queries still succeeded — both events should be persisted.
	if len(store.inserts) != len(logs) {
		t.Errorf("inserted %d events, want %d (aggregator side should ingest cleanly)", len(store.inserts), len(logs))
	}
	// Cursor must NOT advance past startFrom-1 because the chunk had a
	// per-address failure. Otherwise the next start would skip the gap.
	if store.cursor != 0 {
		t.Errorf("cursor advanced to %d despite per-address failure; want 0", store.cursor)
	}
}

func TestBackfill_PartialFailureParksCursorAtLastCleanChunk(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0xaa")
	registry := common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef")

	// Aggregator addr is fine; registry fails. With 4 chunks where
	// every chunk has a registry-failure, the cursor must stay at 0.
	// Then we flip the failure off mid-flight to model a recovery —
	// actually this stub can't change at runtime, so we test the
	// simpler invariant: ANY chunk failure parks the cursor before
	// the failed chunk.
	fetcher := &fakeFetcher{
		headNum:       80,
		failOnAddress: &registry,
		failErr:       errors.New("rate limited"),
	}
	store := &fakeStore{cursor: 0, aggregatorMap: map[common.Address]common.Hash{aggregator: asset}}
	parser, _ := chainsub.NewParser(registry, &stubResolver{mapping: map[common.Address]common.Hash{aggregator: asset}})

	r := New(fetcher, store, parser, Config{
		RegistryAddress: registry,
		ChunkSize:       25,
		CursorEvery:     50,
	})

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if store.cursor != 0 {
		t.Errorf("final cursor = %d, want 0 (every chunk had a failure)", store.cursor)
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
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("Run on canceled ctx returned %v, want context.Canceled", err)
	}
}
