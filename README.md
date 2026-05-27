# evm-oracle-demo-indexer-service

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev/)

Single-chain event indexer for the [EVM Oracle Demo](https://github.com/asolovov). It is the **single chain-observer** for the system: it WS-subscribes to every registered `PriceAggregator` plus the `OracleRegistry`, persists every observed event into `evm_indexer` past the configured confirmation depth, and exposes a confirmation-gated gRPC stream that the rest of the stack (oracle-service, rest-api) consumes.

**The `StreamEvents` server is the confirmation gate.** Subscribers (oracle-service filtering on `PRICE_REQUESTED`, rest-api on all kinds) get a clean post-confirmation feed and never see in-flight or re-orged logs — they don't need to gate on `meta.confirmations` themselves.

The indexer has **no outbound gRPC clients**. It does not call price-service, oracle-service, or anything else. Consumers pull from `indexer.v1.IndexerService/StreamEvents`.

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
            ┌─────────────────────────────────┐
            │     internal/chainsub           │  ← decodes via abigen,
            │   - WS SubscribeFilterLogs       │    resolves asset_id from
            │   - latestRoundData backfill     │    contract address via
            │   - keeps aggregator → asset map │    the registry mapping.
            └────────────┬────────────────────┘
                         │ InsertEvent(confirmations=0)
                         ▼
                   ┌──────────────┐
                   │  evm_indexer │   pgx/v5, typed JSONB payload
                   │  Postgres    │   columns + denormalised asset_id/req_id
                   └──────┬───────┘
                          │ PendingEvents
                          ▼
            ┌─────────────────────────────────┐
            │    internal/confirmer            │
            │   - HeaderByNumber(blockNumber)  │
            │   - bumps confirmations          │
            │   - marks orphaned on reorg      │
            │   - on threshold-cross: Publish  │
            └────────────┬────────────────────┘
                         │
                         ▼
                ┌──────────────────┐
                │ internal/streamhub │   bounded per-subscriber buffer;
                │   (THE GATE)      │   slow consumers dropped.
                └──────┬────────────┘
                       │
                       ▼
       ┌─────────────────────────────────────────┐
       │ internal/grpcsrv — IndexerService:       │
       │   - ListEvents(filter)                    │
       │   - GetRequest(req_id)                    │
       │   - StreamEvents(filter)  ← oracle & BFF  │
       └─────────────────────────────────────────┘

       Cold start fills the gap [cursor+1, head] via
       internal/backfill — same Parser + InsertEvent path
       so confirmer + hub see one ordered firehose.
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
| `IndexerService/ListEvents` | Paginated historical. Filter by `kinds`, `asset_id`, `from_block`, `to_block`. Excludes orphaned + sub-threshold rows. |
| `IndexerService/GetRequest` | Joins `PriceRequested` + matching `PriceFulfilled` by `req_id`. NotFound if no observed request past the threshold. |
| `IndexerService/StreamEvents` | **Confirmation gate.** Long-lived server stream. `from_block > 0` → replay history (ASC), then attach live. `from_block = 0` → live-only. Per-subscriber buffer; slow consumers dropped. |
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
| `CHAIN_BACKFILL_FROM_BLOCK` | `0` | used only when `chain_cursor.last_processed_block = 0` |
| `INDEXER_CONFIRMATIONS` | `5` | threshold for StreamEvents emission |
| `INDEXER_REORG_CHECK_INTERVAL_SEC` | `10` | |
| `INDEXER_BACKFILL_CHUNK_SIZE` | `1000` | blocks per `eth_getLogs` call |
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
  chainsub/               WS subscriber + log parser (abigen-driven)
  confirmer/              ticking confirmation gate, reorg handling
  streamhub/              in-memory pub/sub — the confirmation gate
  backfill/               cold-start gap-fill via eth_getLogs
  grpcsrv/                IndexerService server
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
4. Internal modules are repositories, servers (gRPC + healthz), and in-project handlers. `chainsub`, `confirmer`, `streamhub`, `backfill` are plain packages — not template modules with Init/Start/Stop.
5. `pkg/contracts/` holds abigen-generated bindings (rule 5 exception) — abigen bindings ARE committed; protobuf stubs are not.
6. All config is in `/config/`. `viper.SetDefault` for every nested key.
7. Owns DB `evm_indexer` exclusively.
8. Bootstrap is env-var-driven (no separate seed CLIs / Job containers).
9. `internal/genproto/*.go` is gitignored and regenerated by `make proto-gen` at build time; `protocols/` holds proto sources only. Abigen bindings are the documented exception.

---

## Known gaps / v1 limitations

- **Hot-add of aggregators.** When `AssetRegistered` arrives the persistent mapping + in-memory cache update immediately, but the new aggregator's own logs won't reach the active WS filter — picked up on the next reconnect. Acceptable for the demo where aggregators are seeded at deploy time.
- **PriceFulfilled `round_id`.** The on-chain event doesn't carry `round_id`. The chainsub backfills it best-effort via `latestRoundData()` at the fulfilling block; failures (revert, RPC hiccup) leave the field empty and the proto emits an empty string.
- **`/metrics` cardinality.** Per-kind labels are bounded (3 values); the registry doesn't expose per-asset or per-address counters yet.
- **Single-chain.** The whole rig assumes one chain — multi-chain would require duplicating chainsub + backfill + confirmer per chain, which is out of scope.

---

## Author

Built by **Andrei Solovov** for the EVM Oracle Demo portfolio.

[LinkedIn](https://www.linkedin.com/in/asolovov) · [GitHub](https://github.com/asolovov) · [Live demo](https://github.com/asolovov)
