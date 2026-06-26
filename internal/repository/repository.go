// Package repository provides the pgx-backed persistence layer for
// indexer-service. The only table is `events` — no ORM, raw SQL.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

// ErrNotFound is returned when a single-row lookup matches no rows.
var ErrNotFound = errors.New("not found")

// Repository is the persistence boundary. Implementations are
// goroutine-safe and stateless across calls.
type Repository struct {
	pool *pgxpool.Pool
}

// New constructs a Repository around an already-open pgxpool.
func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Pool exposes the underlying pool for tests + healthchecks. Treat as
// read-only — callers must NOT issue writes outside the repository
// surface.
func (r *Repository) Pool() *pgxpool.Pool { return r.pool }

// Ping forwards to the pool's health probe.
func (r *Repository) Ping(ctx context.Context) error { return r.pool.Ping(ctx) }

// ---------------------------------------------------------------------
// events
// ---------------------------------------------------------------------

// InsertEvent persists a freshly-parsed log. Idempotent: replays of the
// same (tx_hash, log_index) are a no-op.
//
// Returns (true, nil) if a new row was inserted, (false, nil) if the
// row already existed.
func (r *Repository) InsertEvent(ctx context.Context, e *models.Event) (bool, error) {
	payloads, err := encodePayloads(e)
	if err != nil {
		return false, err
	}

	const q = `
INSERT INTO events (
    kind, contract_address, tx_hash, block_hash, block_number, log_index,
    asset_id, req_id,
    price_requested_payload, price_fulfilled_payload, asset_registered_payload,
    observed_at
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8,
    $9, $10, $11,
    $12
)
ON CONFLICT (tx_hash, log_index) DO NOTHING
RETURNING id`

	row := r.pool.QueryRow(ctx, q,
		e.Kind.String(),
		strings.ToLower(e.ContractAddress.Hex()),
		strings.ToLower(e.TxHash.Hex()),
		strings.ToLower(e.BlockHash.Hex()),
		int64(e.BlockNumber), //nolint:gosec // block heights are bounded well under int64 max.
		int32(e.LogIndex),    //nolint:gosec // log indices fit in int32.
		hashOrNil(e.AssetID),
		bigIntOrNil(e.ReqID),
		payloads.priceRequested,
		payloads.priceFulfilled,
		payloads.assetRegistered,
		e.ObservedAt,
	)
	var id int64
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return false, fmt.Errorf("insert event (pg %s): %w", pgErr.Code, err)
		}
		return false, fmt.Errorf("insert event: %w", err)
	}
	e.ID = id
	return true, nil
}

// ListEventsFilter narrows ListEvents queries.
type ListEventsFilter struct {
	Kinds     []models.EventKind
	AssetID   *common.Hash
	FromBlock uint64
	ToBlock   uint64
	Limit     int
	Offset    int
}

// ListEvents returns events matching the filter, sorted descending by
// (block_number, log_index). There is no confirmation gate — events
// are visible the moment they are persisted (emit-on-ingest).
func (r *Repository) ListEvents(ctx context.Context, f ListEventsFilter) ([]*models.Event, error) {
	var (
		args  []any
		where = []string{"TRUE"}
	)

	if len(f.Kinds) > 0 {
		kinds := make([]string, 0, len(f.Kinds))
		for _, k := range f.Kinds {
			if k.IsValid() {
				kinds = append(kinds, k.String())
			}
		}
		if len(kinds) > 0 {
			args = append(args, kinds)
			where = append(where, fmt.Sprintf("kind = ANY($%d)", len(args)))
		}
	}
	if f.AssetID != nil {
		args = append(args, strings.ToLower(f.AssetID.Hex()))
		where = append(where, fmt.Sprintf("asset_id = $%d", len(args)))
	}
	if f.FromBlock > 0 {
		args = append(args, int64(f.FromBlock)) //nolint:gosec // bounded.
		where = append(where, fmt.Sprintf("block_number >= $%d", len(args)))
	}
	if f.ToBlock > 0 {
		args = append(args, int64(f.ToBlock)) //nolint:gosec // bounded.
		where = append(where, fmt.Sprintf("block_number <= $%d", len(args)))
	}

	limit := f.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	args = append(args, limit, offset)

	q := fmt.Sprintf(`
SELECT id, kind, contract_address, tx_hash, block_hash, block_number, log_index,
       asset_id, req_id,
       price_requested_payload, price_fulfilled_payload, asset_registered_payload,
       observed_at
FROM events
WHERE %s
ORDER BY block_number DESC, log_index DESC
LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "),
		len(args)-1, len(args),
	)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list events query: %w", err)
	}
	defer rows.Close()

	return scanEvents(rows)
}

// EventsForRequest returns every event with the supplied req_id,
// ordered ascending by (block_number, log_index). No confirmation gate.
func (r *Repository) EventsForRequest(ctx context.Context, reqID *big.Int) ([]*models.Event, error) {
	if reqID == nil {
		return nil, fmt.Errorf("req_id is required")
	}
	const q = `
SELECT id, kind, contract_address, tx_hash, block_hash, block_number, log_index,
       asset_id, req_id,
       price_requested_payload, price_fulfilled_payload, asset_registered_payload,
       observed_at
FROM events
WHERE req_id = $1
ORDER BY block_number ASC, log_index ASC`

	rows, err := r.pool.Query(ctx, q, reqID.String())
	if err != nil {
		return nil, fmt.Errorf("events for request query: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

// ---------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------

type encodedPayloads struct {
	priceRequested, priceFulfilled, assetRegistered []byte
}

func encodePayloads(e *models.Event) (encodedPayloads, error) {
	var out encodedPayloads
	switch e.Kind {
	case models.EventKindPriceRequested:
		if e.PriceRequested == nil {
			return out, fmt.Errorf("nil PriceRequested payload")
		}
		b, err := json.Marshal(jsonPriceRequested{
			ReqID:     bigIntJSON(e.PriceRequested.ReqID),
			AssetID:   strings.ToLower(e.PriceRequested.AssetID.Hex()),
			Requester: strings.ToLower(e.PriceRequested.Requester.Hex()),
		})
		if err != nil {
			return out, fmt.Errorf("marshal PriceRequested: %w", err)
		}
		out.priceRequested = b
	case models.EventKindPriceFulfilled:
		if e.PriceFulfilled == nil {
			return out, fmt.Errorf("nil PriceFulfilled payload")
		}
		b, err := json.Marshal(jsonPriceFulfilled{
			ReqID:     bigIntJSON(e.PriceFulfilled.ReqID),
			AssetID:   strings.ToLower(e.PriceFulfilled.AssetID.Hex()),
			Price:     bigIntJSON(e.PriceFulfilled.Price),
			Timestamp: bigIntJSON(e.PriceFulfilled.Timestamp),
			RoundID:   bigIntJSONOptional(e.PriceFulfilled.RoundID),
		})
		if err != nil {
			return out, fmt.Errorf("marshal PriceFulfilled: %w", err)
		}
		out.priceFulfilled = b
	case models.EventKindAssetRegistered:
		if e.AssetRegistered == nil {
			return out, fmt.Errorf("nil AssetRegistered payload")
		}
		b, err := json.Marshal(jsonAssetRegistered{
			AssetID:    strings.ToLower(e.AssetRegistered.AssetID.Hex()),
			Aggregator: strings.ToLower(e.AssetRegistered.Aggregator.Hex()),
		})
		if err != nil {
			return out, fmt.Errorf("marshal AssetRegistered: %w", err)
		}
		out.assetRegistered = b
	case models.EventKindUnknown:
		return out, fmt.Errorf("cannot encode payload for EventKindUnknown")
	default:
		return out, fmt.Errorf("unsupported event kind: %s", e.Kind)
	}
	return out, nil
}

type jsonPriceRequested struct {
	ReqID     string `json:"req_id"`
	AssetID   string `json:"asset_id"`
	Requester string `json:"requester"`
}

type jsonPriceFulfilled struct {
	ReqID     string `json:"req_id"`
	AssetID   string `json:"asset_id"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
	RoundID   string `json:"round_id,omitempty"`
}

type jsonAssetRegistered struct {
	AssetID    string `json:"asset_id"`
	Aggregator string `json:"aggregator"`
}

func bigIntJSON(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

func bigIntJSONOptional(n *big.Int) string {
	if n == nil || n.Sign() == 0 {
		return ""
	}
	return n.String()
}

func bigIntOrNil(n *big.Int) any {
	if n == nil {
		return nil
	}
	return n.String()
}

func hashOrNil(h common.Hash) any {
	if (h == common.Hash{}) {
		return nil
	}
	return strings.ToLower(h.Hex())
}

func scanEvents(rows pgx.Rows) ([]*models.Event, error) {
	out := []*models.Event{}
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func scanEvent(rows pgx.Rows) (*models.Event, error) {
	var (
		id, blockNumber                   int64
		kindStr, contractAddr, txHash, bh string
		logIndex                          int32
		assetIDStr, reqIDStr              *string
		prJSON, pfJSON, arJSON            []byte
		observedAt                        time.Time
	)
	if err := rows.Scan(
		&id, &kindStr, &contractAddr, &txHash, &bh, &blockNumber, &logIndex,
		&assetIDStr, &reqIDStr,
		&prJSON, &pfJSON, &arJSON,
		&observedAt,
	); err != nil {
		return nil, fmt.Errorf("scan event row: %w", err)
	}

	kind, err := models.ParseEventKind(kindStr)
	if err != nil {
		return nil, fmt.Errorf("unknown kind in db: %w", err)
	}

	e := &models.Event{
		ID:              id,
		Kind:            kind,
		ContractAddress: common.HexToAddress(contractAddr),
		TxHash:          common.HexToHash(txHash),
		BlockHash:       common.HexToHash(bh),
		BlockNumber:     uint64(blockNumber), //nolint:gosec // already validated >= 0 by DB type.
		LogIndex:        uint32(logIndex),    //nolint:gosec // bounded.
		ObservedAt:      observedAt,
	}

	if assetIDStr != nil && *assetIDStr != "" {
		e.AssetID = common.HexToHash(*assetIDStr)
	}
	if reqIDStr != nil && *reqIDStr != "" {
		n, ok := new(big.Int).SetString(*reqIDStr, 10)
		if !ok {
			return nil, fmt.Errorf("invalid req_id in db: %q", *reqIDStr)
		}
		e.ReqID = n
	}

	switch kind {
	case models.EventKindPriceRequested:
		var p jsonPriceRequested
		if err := json.Unmarshal(prJSON, &p); err != nil {
			return nil, fmt.Errorf("unmarshal PriceRequested payload: %w", err)
		}
		e.PriceRequested = &models.PriceRequestedPayload{
			ReqID:     parseBigInt(p.ReqID),
			AssetID:   common.HexToHash(p.AssetID),
			Requester: common.HexToAddress(p.Requester),
		}
	case models.EventKindPriceFulfilled:
		var p jsonPriceFulfilled
		if err := json.Unmarshal(pfJSON, &p); err != nil {
			return nil, fmt.Errorf("unmarshal PriceFulfilled payload: %w", err)
		}
		e.PriceFulfilled = &models.PriceFulfilledPayload{
			ReqID:     parseBigInt(p.ReqID),
			AssetID:   common.HexToHash(p.AssetID),
			Price:     parseBigInt(p.Price),
			Timestamp: parseBigInt(p.Timestamp),
			RoundID:   parseBigIntOptional(p.RoundID),
		}
	case models.EventKindAssetRegistered:
		var p jsonAssetRegistered
		if err := json.Unmarshal(arJSON, &p); err != nil {
			return nil, fmt.Errorf("unmarshal AssetRegistered payload: %w", err)
		}
		e.AssetRegistered = &models.AssetRegisteredPayload{
			AssetID:    common.HexToHash(p.AssetID),
			Aggregator: common.HexToAddress(p.Aggregator),
		}
	case models.EventKindUnknown:
		return nil, fmt.Errorf("EventKindUnknown is not persistable: row id=%d", e.ID)
	default:
		return nil, fmt.Errorf("unsupported event kind in db row: %s", kind)
	}
	return e, nil
}

func parseBigInt(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

func parseBigIntOptional(s string) *big.Int {
	if s == "" {
		return nil
	}
	return parseBigInt(s)
}
