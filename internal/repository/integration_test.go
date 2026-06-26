//go:build integration

package repository

import (
	"context"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// setupPostgres spins up a postgres:16-alpine container, applies the
// 0001_init migration, and returns a connected repo + teardown.
func setupPostgres(t *testing.T) (*Repository, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("evm_indexer"),
		tcpostgres.WithUsername("indexer_user"),
		tcpostgres.WithPassword("indexerpw"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("dsn: %v", err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}

	// Apply migration directly (testcontainers' MigrateUp helper has
	// changed shape over versions; we just read the file and Exec).
	migPath := filepath.Join("..", "..", "migrations", "0001_init.up.sql")
	sql, err := os.ReadFile(migPath)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	repo := New(pool)
	teardown := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}
	return repo, teardown
}

func mkPriceRequested(id int64, block uint64) *models.Event {
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	return &models.Event{
		Kind:            models.EventKindPriceRequested,
		ContractAddress: common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281"),
		TxHash:          common.BigToHash(big.NewInt(id)),
		BlockHash:       common.HexToHash("0xbeef"),
		BlockNumber:     block,
		LogIndex:        uint32(id),
		ObservedAt:      time.Now().UTC(),
		AssetID:         asset,
		ReqID:           big.NewInt(id),
		PriceRequested: &models.PriceRequestedPayload{
			ReqID:     big.NewInt(id),
			AssetID:   asset,
			Requester: common.HexToAddress("0xCEF4Fe1Ca9071f4ED4BAd6c1087CEb08838a983E"),
		},
	}
}

func mkPriceFulfilled(id int64, block uint64) *models.Event {
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	return &models.Event{
		Kind:            models.EventKindPriceFulfilled,
		ContractAddress: common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281"),
		TxHash:          common.BigToHash(big.NewInt(id + 1000)),
		BlockHash:       common.HexToHash("0xbeef"),
		BlockNumber:     block,
		LogIndex:        uint32(id),
		ObservedAt:      time.Now().UTC(),
		AssetID:         asset,
		ReqID:           big.NewInt(id),
		PriceFulfilled: &models.PriceFulfilledPayload{
			ReqID:     big.NewInt(id),
			AssetID:   asset,
			Price:     big.NewInt(345020000000),
			Timestamp: big.NewInt(1700000000),
		},
	}
}

func TestIntegration_InsertEventIdempotent(t *testing.T) {
	repo, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	evt := mkPriceRequested(1, 100)
	inserted, err := repo.InsertEvent(ctx, evt)
	if err != nil || !inserted {
		t.Fatalf("first insert: inserted=%v err=%v", inserted, err)
	}
	if evt.ID == 0 {
		t.Error("event ID not set after insert")
	}
	again, err := repo.InsertEvent(ctx, evt)
	if err != nil {
		t.Fatalf("second insert err: %v", err)
	}
	if again {
		t.Error("replay should have returned inserted=false")
	}
}

func TestIntegration_ListEventsVisibleImmediately(t *testing.T) {
	repo, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	// Emit-on-ingest: an event is visible the moment it's persisted —
	// no confirmation gate.
	evt := mkPriceRequested(2, 200)
	if _, err := repo.InsertEvent(ctx, evt); err != nil {
		t.Fatalf("insert: %v", err)
	}

	got, err := repo.ListEvents(ctx, ListEventsFilter{})
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected the event immediately, got %d", len(got))
	}
	if got[0].PriceRequested == nil || got[0].PriceRequested.ReqID.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("payload roundtrip mismatch: %+v", got[0].PriceRequested)
	}

	// Kind filter narrows.
	none, err := repo.ListEvents(ctx, ListEventsFilter{Kinds: []models.EventKind{models.EventKindAssetRegistered}})
	if err != nil {
		t.Fatalf("ListEvents kind filter: %v", err)
	}
	if len(none) != 0 {
		t.Errorf("kind filter leaked non-matching rows: %d", len(none))
	}
}

func TestIntegration_EventsForRequest(t *testing.T) {
	repo, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	req := mkPriceRequested(7, 100)
	ful := mkPriceFulfilled(7, 110)
	for _, e := range []*models.Event{req, ful} {
		if _, err := repo.InsertEvent(ctx, e); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	got, err := repo.EventsForRequest(ctx, big.NewInt(7))
	if err != nil {
		t.Fatalf("EventsForRequest: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d events, want 2", len(got))
	}
	// Should be ASC by (block, log_index).
	if got[0].BlockNumber > got[1].BlockNumber {
		t.Errorf("events not sorted ASC by block_number: %+v", got)
	}
}

func TestIntegration_AssetRegisteredSeedIdempotent(t *testing.T) {
	repo, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")

	// First seed inserts; re-seed is a no-op (deterministic synthetic
	// tx_hash + UNIQUE(tx_hash, log_index)).
	first, err := repo.InsertEvent(ctx, models.NewAssetRegisteredSeed(asset, agg))
	if err != nil || !first {
		t.Fatalf("first seed: inserted=%v err=%v", first, err)
	}
	again, err := repo.InsertEvent(ctx, models.NewAssetRegisteredSeed(asset, agg))
	if err != nil {
		t.Fatalf("re-seed: %v", err)
	}
	if again {
		t.Error("re-seed should be idempotent (inserted=false)")
	}

	got, err := repo.ListEvents(ctx, ListEventsFilter{Kinds: []models.EventKind{models.EventKindAssetRegistered}})
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(got) != 1 || got[0].AssetRegistered == nil || got[0].AssetRegistered.Aggregator != agg {
		t.Errorf("seeded AssetRegistered not readable as expected: %+v", got)
	}
}
