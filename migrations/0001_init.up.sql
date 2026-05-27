-- evm_indexer initial schema.
--
-- Single-chain — no chain_id column. Per-kind typed JSONB payload
-- columns (one is populated per row, the others NULL) keep queries
-- ergonomic and avoid the single opaque `bytes payload` blob the
-- earlier spec draft proposed.

CREATE TABLE IF NOT EXISTS events (
    id                          BIGSERIAL PRIMARY KEY,
    kind                        TEXT        NOT NULL,
    contract_address            TEXT        NOT NULL,
    tx_hash                     TEXT        NOT NULL,
    block_hash                  TEXT        NOT NULL,
    block_number                BIGINT      NOT NULL,
    log_index                   INTEGER     NOT NULL,

    -- Denormalised lookup columns (also present inside the payload
    -- JSONB blobs). Pulled out so the indices below can serve
    -- ListEvents / GetRequest queries without touching JSONB.
    asset_id                    TEXT,
    req_id                      TEXT,

    price_requested_payload     JSONB,
    price_fulfilled_payload     JSONB,
    asset_registered_payload    JSONB,

    observed_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    confirmations               INTEGER     NOT NULL DEFAULT 0,
    orphaned                    BOOLEAN     NOT NULL DEFAULT FALSE,

    CONSTRAINT events_kind_check
        CHECK (kind IN ('PRICE_REQUESTED', 'PRICE_FULFILLED', 'ASSET_REGISTERED')),

    -- A single on-chain log is uniquely identified by (tx_hash, log_index).
    CONSTRAINT events_log_unique UNIQUE (tx_hash, log_index)
);

-- Recent-events feed for ListEvents (kinds filter + block range).
CREATE INDEX IF NOT EXISTS events_kind_block_desc
    ON events (kind, block_number DESC);

-- Asset-scoped queries (dashboard drill-down).
CREATE INDEX IF NOT EXISTS events_asset_kind_block_desc
    ON events (asset_id, kind, block_number DESC)
    WHERE asset_id IS NOT NULL;

-- GetRequest joins PriceRequested + PriceFulfilled on req_id.
CREATE INDEX IF NOT EXISTS events_req_id
    ON events (req_id)
    WHERE req_id IS NOT NULL;

-- Confirmer scan: walks through not-yet-final, not-yet-orphaned events.
-- The threshold (N=5 default) is NOT baked into the index — the query
-- adds `confirmations < N` itself, so changes to Confirmations don't
-- invalidate the index.
CREATE INDEX IF NOT EXISTS events_pending
    ON events (block_number)
    WHERE orphaned = FALSE;


-- Single-row table holding the last block successfully drained from
-- the chain. PK is forced to 1 by an INSERT in this migration so
-- upserts always target the same row.
CREATE TABLE IF NOT EXISTS chain_cursor (
    id                      INTEGER PRIMARY KEY,
    last_processed_block    BIGINT      NOT NULL DEFAULT 0,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chain_cursor_singleton CHECK (id = 1)
);

INSERT INTO chain_cursor (id, last_processed_block) VALUES (1, 0)
    ON CONFLICT (id) DO NOTHING;


-- Aggregator address -> assetId mapping populated from
-- OracleRegistry.AssetRegistered events + the startup listAssets()
-- enumeration. Used by chainsub to attach asset_id to PriceRequested
-- / PriceFulfilled events that don't carry it natively.
CREATE TABLE IF NOT EXISTS aggregator_registry (
    aggregator      TEXT PRIMARY KEY,           -- 20-byte lowercase 0x hex
    asset_id        TEXT NOT NULL,              -- bytes32 lowercase 0x hex
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
