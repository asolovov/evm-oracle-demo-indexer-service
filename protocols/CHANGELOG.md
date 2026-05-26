# Changelog

All notable changes to `evm-oracle-demo-protocols` are documented here.
The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — TBD

Initial release. Cut by the human after task 05 (`evm-oracle-demo-price-service`) confirms downstream subtree integration end-to-end.

### Added

- `price/v1/price.proto` — `PriceService.GetPrice` + `Subscribe`. Prices flow as IEEE-754 `double` end-to-end inside the off-chain pipeline; conversion to on-chain int256 happens once, in oracle-service.
- `oracle/v1/oracle.proto` — `OracleService.SetHeartbeat` (admin), `GetSubmissionStatus`, `ListSubmissions`. No `TriggerUpdate` RPC: oracle-service subscribes to `indexer.StreamEvents` for trigger events; heartbeats are internal. `SubmissionStatus.Status` covers pending / confirmed / failed / dropped.
- `indexer/v1/indexer.proto` — `IndexerService.ListEvents`, `StreamEvents`, `GetRequest`. Single-chain. Typed `oneof` over `PriceRequestedEvent` / `PriceFulfilledEvent` / `AssetRegisteredEvent`, discriminated by `EventKind`. `StreamEvents` emits only past-confirmation events.
- Repo-wide scalar conventions (addresses, hashes, uint256 / int256 as strings; off-chain prices as `double`; native proto integers for non-amount counters); documented in README.
- `buf breaking` baseline against `origin/<base_ref>`; first PR after this tag locks the wire format.
- All `.github/workflows/` actions SHA-pinned (NFR-07).

### Removed

- Template example `user/` service.
- Template `common/v1/status.proto` and `common/v1/types.proto` (UUID, Money, Address, CommonStatus) — none of which the oracle pipeline uses.

[Unreleased]: https://github.com/asolovov/evm-oracle-demo-protocols/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/asolovov/evm-oracle-demo-protocols/releases/tag/v0.1.0
