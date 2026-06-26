# evm-oracle-demo-indexer-service

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev/)

Single-chain event indexer for the [EVM Oracle Demo](https://github.com/asolovov). It is the **single chain-observer** for the system: it WS-subscribes to every registered `PriceAggregator` plus the `OracleRegistry`, persists every observed event into `evm_indexer`, and exposes a gRPC stream that the rest of the stack (oracle-service, rest-api) consumes.

**Emit-on-ingest — no confirmation gate.** An event is published to `StreamEvents` subscribers the moment it is persisted, for the lowest possible latency ("fire ASAP"). There is no confirmation threshold and no reorg/orphan tracking inside the indexer.

> **Reorg trade-off (read this).** Because events fire at 0 confirmations, a chain reorg can roll back an event that was already streamed. Re-orged-out events are NOT retracted; a replacement log at a different `(tx_hash, log_index)` is ingested as an additional event. This is a deliberate demo choice. **Finality is a downstream contract:** any consumer that takes an irreversible action (e.g. oracle-service signing a submission) must apply its own finality guard before acting. The indexer's job is speed, not finality.

The indexer has **no outbound gRPC clients**. It does not call price-service, oracle-service, or anything else. Consumers pull from `indexer.v1.IndexerService/StreamEvents`.

Public-RPC WebSockets drop constantly, so **reconnect is the steady state, not an error path**: on every (re)connect the subscriber subscribes live first, then catches up the gap since the persisted `chain_cursor` via `eth_getLogs`, so no event is lost across a disconnect. Reconnect uses exponential backoff with jitter. The chain client is owned end-to-end by one goroutine (no shared state, no data race).

---

## Architecture

```
                          chain (single, e.g. Ethereum Sepolia)
                              │
                  ┌───────────┼───────────────────────┐
                  │           │                       │
       OracleRegistry   PriceAggregator(s)     wss:// WS subscription
                              │
                              ▼
            ┌─────────────────────────────────────┐
            │     internal/chainsub  (1 goroutine) │  ← owns the chain
            │   per (re)connect, with backoff:     │    client end-to-end
            │   1. SubscribeFilterLogs (live)      │    (no shared state,
            │   2. catch-up [cursor+1..head] via   │    no data race).
            │      eth_getLogs (per-address chunk) │    Decodes via abigen,
            │   3. drain live; advance chain_cursor│    resolves asset_id
            │   ingest = parse → persist → PUBLISH  │    from the registry.
            └───────┬───────────────────┬──────────┘
                    │ InsertEvent        │ Publish (on insert, 0-conf)
                    ▼                    ▼
            ┌──────────────┐    ┌──────────────────┐
            │  evm_indexer │    │ internal/streamhub │ bounded per-sub
            │  Postgres    │    │  live pub/sub      │ buffer; slow + over-
            │  pgx/v5      │    │  (subscriber cap)  │ cap subs dropped.
            └──────┬───────┘    └──────┬─────────────┘
                   │                   │
                   ▼                   ▼
       ┌─────────────────────────────────────────┐
       │ internal/grpcsrv — IndexerService:       │
       │   - ListEvents(filter)        (reads DB) │
       │   - GetRequest(req_id)        (reads DB) │
       │   - StreamEvents(filter)  ← oracle & BFF │
       │       replay (DB) then live (hub);       │
       │       log-granular (block,log_index)     │
       │       dedup at the boundary.             │
       └─────────────────────────────────────────┘

       No confirmation gate, no confirmer, no reorg/orphan tracking.
       Catch-up runs on EVERY connect, so a disconnect gap is refilled.
```

---

## Run locally

```bash
make build              # codegen + compile, binary at ./evm-oracle-demo-indexer-service
make test               # unit tests
make test-integration   # docker-backed integration suite (testcontainers)
make lint               # golangci-lint
```

```bash
# With docker-compose (brings up Postgres + migrate + indexer)
export CHAIN_WS_URL=wss://...
export CHAIN_RPC_URL=https://...
docker compose up --build
```

### gRPC surface

| RPC | Notes |
| --- | --- |
| `IndexerService/ListEvents` | Paginated historical. Filter by `kinds`, `asset_id`, `from_block`, `to_block`. |
| `IndexerService/GetRequest` | Joins `PriceRequested` + matching `PriceFulfilled` by `req_id`. NotFound if the request hasn't been observed. |
| `IndexerService/StreamEvents` | Long-lived server stream, emit-on-ingest (no confirmation gate). `from_block > 0` → replay history (ASC) then attach live, with log-granular `(block, log_index)` dedup at the boundary; `from_block = 0` → live-only. Per-subscriber buffer + subscriber cap; slow/over-cap consumers dropped. |
| `grpc.health.v1.Health/Check` | Standard. Reflection is on by default for `grpcurl`. |

### HTTP surface

| Path | Notes |
| --- | --- |
| `/healthz` | Liveness — 200 once the listener is up. |
| `/readyz` | Walks every module's `HealthCheck`; 503 on any failure. |
| `/metrics` | Prometheus. Counters / gauges defined in `internal/metrics`. |

---

## Configuration

Every value is set via env (`viper.SetDefault` registers a default for every key per architecture rule 6). Required-but-secret keys default to empty strings; `config.Validate` fails fast on missing values.

| Env var | Default | Notes |
| --- | --- | --- |
| `ENV` | `prod` | |
| `DATABASE_HOST` | `localhost` | |
| `DATABASE_PORT` | `5432` | |
| `DATABASE_USER` | `indexer_user` | |
| `DATABASE_PASSWORD` | _empty (required)_ | |
| `DATABASE_NAME` | `evm_indexer` | dedicated DB, never shared |
| `DATABASE_SSL_MODE` | `disable` | |
| `DATABASE_MAX_OPEN_CONNS` | `25` | |
| `DATABASE_MAX_IDLE_CONNS` | `5` | |
| `DATABASE_CONN_MAX_LIFETIME` | `300` (seconds) | |
| `GRPC_HOST` | `0.0.0.0` | |
| `GRPC_PORT` | `9090` | |
| `GRPC_MAX_RECV_MSG_SIZE` | `62914560` (60MB) | |
| `GRPC_MAX_SEND_MSG_SIZE` | `62914560` (60MB) | |
| `GRPC_REFLECTION` | `true` | |
| `HEALTHZ_HOST` | `0.0.0.0` | |
| `HEALTHZ_PORT` | `8080` | |
| `CHAIN_NAME` | `ethereum-sepolia` | |
| `CHAIN_CHAIN_ID` | `11155111` | |
| `CHAIN_WS_URL` | _empty (required)_ | |
| `CHAIN_RPC_URL` | _empty (required)_ | |
| `CHAIN_REGISTRY_ADDRESS` | _empty (required)_ | 0x-prefixed 20-byte hex |
| `CHAIN_BACKFILL_FROM_BLOCK` | `0` | catch-up seed when `chain_cursor.last_processed_block = 0`. Set to a recent block to avoid a huge cold-start replay. |
| `INDEXER_BACKFILL_CHUNK_SIZE` | `1000` | blocks per `eth_getLogs` call (validated to be in `[1, 10000]`) |
| `INDEXER_STREAM_SUBSCRIBER_BUFFER` | `256` | per-subscriber outbound buffer (overflow → drop) |
| `TELEMETRY_LOG_LEVEL` | `info` | |
| `TELEMETRY_LOG_FORMAT` | `json` | |

---

## Project layout

```
config/                   config scheme + viper defaults + Validate()
internal/
  application.go          SINGLE wiring point (architecture rules 1+2)
  models/                 domain types + proto<->model conversions (rule 3)
  repository/             pgx repository for events / chain_cursor / aggregator_registry
  chainsub/               WS subscriber + per-connect catch-up + log parser
                          (abigen-driven); owns the chain client, publishes
                          on ingest
  streamhub/              in-memory live pub/sub (per-sub buffer + cap)
  grpcsrv/                IndexerService server (replay-then-live stream)
  metrics/                Prometheus registry + adapters
  healthz/                /healthz, /readyz, /metrics listener
  module/                 lifecycle manager (template-provided, retained)
pkg/
  contracts/              abigen bindings (committed — rule 5 exception)
    oracleregistry/
    priceaggregator/
    reporterset/
  logger/                 logrus accessor
  version/                build-time release metadata
protocols/                git subtree of evm-oracle-demo-protocols
                          (proto-source-only, language-agnostic)
migrations/               golang-migrate scripts for evm_indexer
```

---

## Architecture rules honoured

1. `cmd/` does CLI + config init only.
2. `internal/application.go` is the single wiring point; nothing self-wires.
3. Every domain type + every proto/JSONB conversion lives in `internal/models/`.
4. Internal modules are repositories, servers (gRPC + healthz), and in-project handlers. `chainsub` + `streamhub` are plain packages — not template modules with Init/Start/Stop.
5. `pkg/contracts/` holds abigen-generated bindings (rule 5 exception) — abigen bindings ARE committed; protobuf stubs are not.
6. All config is in `/config/`. `viper.SetDefault` for every nested key.
7. Owns DB `evm_indexer` exclusively.
8. Bootstrap is env-var-driven (no separate seed CLIs / Job containers).
9. `internal/genproto/*.go` is gitignored and regenerated by `make proto-gen` at build time; `protocols/` holds proto sources only. Abigen bindings are the documented exception.

---

## Known gaps / v1 limitations

- **Hot-add of aggregators.** When `AssetRegistered` arrives the persistent mapping + in-memory cache update immediately, but the new aggregator's own logs won't reach the active WS filter — picked up on the next reconnect. Acceptable for the demo where aggregators are seeded at deploy time.
- **PriceFulfilled `round_id`.** The on-chain event doesn't carry `round_id`. chainsub backfills it best-effort via a bounded `latestRoundData()` call at the fulfilling block — **only on the live tail**, skipped during catch-up (the archival call would hammer the RPC). Failures leave the field empty and the proto emits an empty string.
- **`/metrics` cardinality.** Per-kind labels are bounded (3 values); the registry doesn't expose per-asset or per-address counters yet.
- **Single-chain.** The whole rig assumes one chain — multi-chain would require one chainsub per chain, which is out of scope.

---

## Author

Built by **Andrei Solovov** for the EVM Oracle Demo portfolio.

[LinkedIn](https://www.linkedin.com/in/asolovov) · [GitHub](https://github.com/asolovov) · [Live demo](https://github.com/asolovov)
