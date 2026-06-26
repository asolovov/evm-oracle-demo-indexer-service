# evm-oracle-demo-indexer-service

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev/)

Single-chain event indexer for the [EVM Oracle Demo](https://github.com/asolovov). It is the **single chain-observer** for the system: it WS-subscribes to the deployed `PriceAggregator`s plus the `OracleRegistry`, persists every observed event into `evm_indexer`, and exposes a gRPC stream that the rest of the stack (oracle-service, api) consumes.

**Live-only + emit-on-ingest.** The only chain operation is a WebSocket log subscription — **no historical `eth_getLogs`, no `eth_call`, no archive queries.** That's deliberate: free RPC tiers don't serve historical logs (they 403/429), and a demo only needs events as they happen. Each event is published to `StreamEvents` subscribers the moment it is persisted, for the lowest possible latency.

The aggregator→asset mapping comes from **config** (the deployed asset set — see `INDEXER_ASSETS`), not from reading the registry on chain. On startup the indexer **bootstraps one `AssetRegistered` event per configured asset** into `evm_indexer` (idempotent) so the API can read the asset set without any historical backfill. Live `AssetRegistered` logs extend the set at runtime.

> **Trade-offs (read this).**
> - **No backfill.** Events emitted while the indexer is disconnected, or before it starts, are not recovered — only live events from subscription onward. (A live demo triggers events and watches them stream, so this is the right call.)
> - **No reorg handling.** Events fire at 0 confirmations; a reorg can roll back an already-streamed event. **Finality is a downstream contract:** any consumer that acts irreversibly (oracle-service signing) applies its own finality guard. The indexer's job is speed.

The indexer has **no outbound gRPC clients**. Consumers pull from `indexer.v1.IndexerService/StreamEvents`.

Public-RPC WebSockets drop often, so **reconnect is the steady state, not an error path**: the run loop reconnects with exponential backoff + jitter. The WS client is owned end-to-end by one goroutine (no shared state, no data race).

---

## Architecture

```
   config assets ──► bootstrap: seed 1 AssetRegistered / asset (idempotent) ──► evm_indexer
   (INDEXER_ASSETS)                                                                  ▲
                          chain (single, e.g. Ethereum Sepolia)                      │
                              │  wss:// log subscription ONLY                        │
                              ▼                                                      │
            ┌─────────────────────────────────────┐                                 │
            │     internal/chainsub  (1 goroutine) │  live-only; owns the WS         │
            │   reconnect loop (exp backoff+jitter)│  client end-to-end (no          │
            │   SubscribeFilterLogs(registry +     │  shared state, no race).        │
            │     configured aggregators)          │  Decodes via abigen;            │
            │   ingest = parse → persist ──────────┼──► InsertEvent ─────────────────┘
            │                  └────────► PUBLISH   │  (emit-on-ingest, 0-conf)
            └──────────────────────────┬───────────┘
                                       │ Publish
                                       ▼
                            ┌──────────────────┐
                            │ internal/streamhub │ bounded per-sub buffer;
                            │  live pub/sub      │ slow + over-cap subs dropped.
                            └──────┬─────────────┘
                                   │
                                   ▼
       ┌─────────────────────────────────────────┐
       │ internal/grpcsrv — IndexerService:       │
       │   - ListEvents(filter)        (reads DB) │
       │   - GetRequest(req_id)        (reads DB) │
       │   - StreamEvents(filter)  ← oracle & API │
       │       replay (DB) then live (hub);       │
       │       log-granular (block,log_index)     │
       │       dedup at the boundary.             │
       └─────────────────────────────────────────┘

       No eth_getLogs / eth_call, no confirmer, no cursor, no reorg
       tracking. The only chain call is the WS log subscription.
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
# With docker-compose (brings up Postgres + migrate + indexer).
# CHAIN_WS_URL is the only required chain setting — e.g. an Alchemy
# free-tier wss endpoint. No RPC URL: the indexer is live-only.
export CHAIN_WS_URL=wss://eth-sepolia.g.alchemy.com/v2/<your-key>
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
| `CHAIN_WS_URL` | _empty (required)_ | WebSocket endpoint — the only chain connection (e.g. Alchemy free-tier wss) |
| `CHAIN_REGISTRY_ADDRESS` | deployed registry | 0x-prefixed 20-byte hex; subscribed to for live `AssetRegistered` |
| `INDEXER_STREAM_SUBSCRIBER_BUFFER` | `256` | per-subscriber outbound buffer (overflow → drop) |
| `INDEXER_ASSETS` | _baked 10-asset default_ | optional JSON array `[{"symbol","asset_id","aggregator"}]` to override the deployed asset set |
| `TELEMETRY_LOG_LEVEL` | `info` | |
| `TELEMETRY_LOG_FORMAT` | `json` | |

---

## Project layout

```
config/                   config scheme + viper defaults + Validate()
internal/
  application.go          SINGLE wiring point (architecture rules 1+2)
  models/                 domain types + proto<->model conversions (rule 3)
  repository/             pgx repository for the events table
  chainsub/               live WS subscriber + abigen log parser; owns the
                          WS client, publishes on ingest; mapping from config
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
8. Bootstrap is env-var-driven (rule 8): the configured asset set seeds `AssetRegistered` events into the DB on startup via an idempotent upsert — no separate seed CLIs / Job containers.
9. `internal/genproto/*.go` is gitignored and regenerated by `make proto-gen` at build time; `protocols/` holds proto sources only. Abigen bindings are the documented exception.

---

## Known gaps / v1 limitations

- **No backfill.** Live-only: events emitted while disconnected, or before startup, are not recovered. Deliberate — free RPC tiers don't serve historical logs, and a live demo doesn't need it.
- **No reorg handling.** Emit at 0 confirmations; finality is the consumer's responsibility (see Trade-offs at the top).
- **PriceFulfilled `round_id`** is always empty. The on-chain event doesn't carry it, and the indexer makes no `eth_call` to recover it (that would need archive access). The proto field is emitted as an empty string.
- **Hot-add of aggregators.** A live `AssetRegistered` updates the in-memory mapping immediately, but the new aggregator's own price logs reach the WS filter only on the next reconnect. Fine for the demo (the asset set is fixed in config).
- **Asset set is config-trusted.** The aggregator→asset mapping is the configured `INDEXER_ASSETS` (defaults to the deployed 10); the indexer does not verify it against the on-chain registry.
- **Single-chain.** Multi-chain would require one chainsub per chain — out of scope.

---

## Author

Built by **Andrei Solovov** for the EVM Oracle Demo portfolio.

[LinkedIn](https://www.linkedin.com/in/asolovov) · [GitHub](https://github.com/asolovov) · [Live demo](https://github.com/asolovov)
