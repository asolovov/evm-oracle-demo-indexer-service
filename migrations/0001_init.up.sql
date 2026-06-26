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

-- NOTE: this is the ONLY table. The indexer is live-only (no historical
-- catch-up → no chain_cursor) and the aggregator->asset mapping comes
-- from config (no aggregator_registry table). The AssetRegistered events
-- the API needs are bootstrapped into `events` on startup from config.
