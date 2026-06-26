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
// Single-chain by design (spec §1 / §3.2). The indexer is LIVE-ONLY:
// it WS-subscribes and never queries history, so there is no RPC
// endpoint and no backfill block — only the WS URL. (Free RPC tiers
// don't serve historical eth_getLogs anyway.)
type ChainConfig struct {
	// Friendly chain name (e.g. "ethereum-sepolia"). Surfaced in logs.
	Name string `mapstructure:"name"`

	// EIP-155 chain id (e.g. 11155111 for Ethereum Sepolia).
	ChainID uint64 `mapstructure:"chain_id"`

	// WebSocket endpoint used for live `SubscribeFilterLogs`. The only
	// chain connection the indexer makes.
	WSURL string `mapstructure:"ws_url"`

	// 20-byte 0x-prefixed lowercase hex of the deployed OracleRegistry.
	// Subscribed to (for live AssetRegistered logs); never enumerated
	// on chain — the asset set comes from Assets below.
	RegistryAddress string `mapstructure:"registry_address"`
}

// AssetConfig is one entry of the deployed asset set. It is the single
// source of truth the indexer needs about the deployment (the indexer
// no longer reads the OracleRegistry on chain): it seeds the
// aggregator->asset mapping for the live subscription AND the
// AssetRegistered events bootstrapped into the DB on startup.
// json tags matter: the INDEXER_ASSETS override is decoded with
// encoding/json (see ApplyEnvOverrides), which ignores mapstructure
// tags — without json tags the fields would silently stay empty.
type AssetConfig struct {
	Symbol     string `mapstructure:"symbol"     json:"symbol"`
	AssetID    string `mapstructure:"asset_id"   json:"asset_id"`   // bytes32 0x hex
	Aggregator string `mapstructure:"aggregator" json:"aggregator"` // 20-byte 0x hex
}

// IndexerConfig holds the stream-hub knob + the deployed asset set.
type IndexerConfig struct {
	// Per-subscriber outbound queue depth on the stream hub. When a
	// subscriber lags past this, it gets dropped (backpressure).
	StreamSubscriberBuffer int `mapstructure:"stream_subscriber_buffer"`

	// Assets is the deployed asset set (aggregator + asset_id per
	// asset). Defaults to the 10 Ethereum-Sepolia aggregators; override
	// via the INDEXER_ASSETS env (JSON array) for a different deploy.
	Assets []AssetConfig `mapstructure:"assets"`
}

// TelemetryConfig holds logging + metrics knobs.
type TelemetryConfig struct {
	LogLevel  string `mapstructure:"log_level"`  // debug|info|warn|error
	LogFormat string `mapstructure:"log_format"` // json|text
}
