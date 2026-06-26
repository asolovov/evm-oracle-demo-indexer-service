// Package config defines indexer-service configuration defaults and schema.
//
// Only the modules this service actually uses are described here
// (architecture rule 4: indexer is a server-only, observer-only rig).
// Notably absent:
//
//   - http        — indexer serves over gRPC; healthz lives under
//     its own dedicated http listener (see HealthzConfig).
//   - grpc_client — indexer does NOT call out to any other service;
//     consumers pull from StreamEvents instead.
//   - websocket   — not exposed by this service.
package config

// Scheme is the indexer-service configuration scheme.
type Scheme struct {
	Database  *DatabaseConfig  `mapstructure:"database"`
	GRPC      *GRPCConfig      `mapstructure:"grpc"`
	Healthz   *HealthzConfig   `mapstructure:"healthz"`
	Chain     *ChainConfig     `mapstructure:"chain"`
	Indexer   *IndexerConfig   `mapstructure:"indexer"`
	Telemetry *TelemetryConfig `mapstructure:"telemetry"`
	Env       string           `mapstructure:"env"`
}

// DatabaseConfig holds the connection to the dedicated `evm_indexer` Postgres database.
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"` // always "postgres"
	Host            string `mapstructure:"host"`
	Name            string `mapstructure:"name"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	SSLMode         string `mapstructure:"ssl_mode"`
	Port            int    `mapstructure:"port"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // seconds
}

// GRPCConfig holds gRPC server settings.
type GRPCConfig struct {
	Host             string `mapstructure:"host"`
	Timeout          string `mapstructure:"timeout"`
	MaxSendMsgSize   int    `mapstructure:"max_send_msg_size"`
	MaxRecvMsgSize   int    `mapstructure:"max_recv_msg_size"`
	Port             int    `mapstructure:"port"`
	NumStreamWorkers uint32 `mapstructure:"num_stream_workers"`
	Reflection       bool   `mapstructure:"reflection"`
}

// HealthzConfig holds the dedicated HTTP listener for /healthz, /readyz, /metrics.
type HealthzConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ChainConfig describes the single target EVM chain the indexer observes.
//
// Single-chain by design — running multiple chains simultaneously is
// out of scope (spec §1 / §3.2). To swap chains, change the values
// here and replay from BackfillFromBlock.
type ChainConfig struct {
	// Friendly chain name (e.g. "ethereum-sepolia"). Surfaced in logs and metrics.
	Name string `mapstructure:"name"`

	// EIP-155 chain id (e.g. 11155111 for Ethereum Sepolia).
	ChainID uint64 `mapstructure:"chain_id"`

	// WebSocket endpoint used for live `SubscribeFilterLogs`.
	WSURL string `mapstructure:"ws_url"`

	// JSON-RPC endpoint used for `eth_getLogs`, `latestRoundData`, etc.
	RPCURL string `mapstructure:"rpc_url"`

	// 20-byte 0x-prefixed lowercase hex of the deployed OracleRegistry.
	RegistryAddress string `mapstructure:"registry_address"`

	// Block height to start backfilling from on a cold start. Pinned to
	// the block at/just before the contracts were deployed so backfill
	// is bounded.
	BackfillFromBlock uint64 `mapstructure:"backfill_from_block"`
}

// IndexerConfig holds chainsub catch-up + stream-hub knobs. There is
// no confirmation gate — events emit on ingest — so the old
// Confirmations / ReorgCheckIntervalSec knobs are gone.
type IndexerConfig struct {
	// Block-chunk size for `eth_getLogs` catch-up calls. Bounded in
	// Validate so a misconfigured value can't ask a public RPC for an
	// enormous range in one call.
	BackfillChunkSize uint64 `mapstructure:"backfill_chunk_size"`

	// Per-subscriber outbound queue depth on the stream hub. When a
	// subscriber lags past this, it gets dropped (backpressure).
	StreamSubscriberBuffer int `mapstructure:"stream_subscriber_buffer"`
}

// TelemetryConfig holds logging + metrics knobs.
type TelemetryConfig struct {
	LogLevel  string `mapstructure:"log_level"`  // debug|info|warn|error
	LogFormat string `mapstructure:"log_format"` // json|text
}
