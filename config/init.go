// Package config defines indexer-service configuration defaults and schema.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

//nolint:gochecknoinits // configuration defaults are registered at package load.
func init() {
	setDefaults()
}

// setDefaults registers a viper default for every key the service
// reads — required by architecture rule 6 (viper.AutomaticEnv alone
// does not populate nested keys on Unmarshal). Required-but-secret
// values (DSN passwords, RPC URLs) default to empty strings and are
// rejected by Validate at startup if still empty post-load.
//
//nolint:funlen // a single linear list of viper defaults is more readable than fragmented helpers.
func setDefaults() {
	// Core
	viper.SetDefault("env", "prod")

	// Database — owns DB `evm_indexer`, user `indexer_user`.
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 300)
	viper.SetDefault("database.user", "indexer_user")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "evm_indexer")

	// gRPC server (always enabled — indexer is server-only).
	viper.SetDefault("grpc.host", "0.0.0.0")
	viper.SetDefault("grpc.port", 9090)
	viper.SetDefault("grpc.timeout", "30s")
	viper.SetDefault("grpc.max_send_msg_size", 60*1024*1024)
	viper.SetDefault("grpc.max_recv_msg_size", 60*1024*1024)
	viper.SetDefault("grpc.num_stream_workers", 0)
	viper.SetDefault("grpc.reflection", true)

	// Healthz / readyz / metrics listener.
	viper.SetDefault("healthz.host", "0.0.0.0")
	viper.SetDefault("healthz.port", 8080)

	// Chain — must be overridden via env at deploy time.
	viper.SetDefault("chain.name", "ethereum-sepolia")
	viper.SetDefault("chain.chain_id", uint64(11155111))
	viper.SetDefault("chain.ws_url", "")
	viper.SetDefault("chain.rpc_url", "")
	viper.SetDefault("chain.registry_address", "")
	viper.SetDefault("chain.backfill_from_block", uint64(0))

	// Indexer knobs.
	viper.SetDefault("indexer.backfill_chunk_size", uint64(1000))
	viper.SetDefault("indexer.stream_subscriber_buffer", 256)

	// Telemetry.
	viper.SetDefault("telemetry.log_level", "info")
	viper.SetDefault("telemetry.log_format", "json")
}

// Validate fails fast on misconfiguration so an orchestrator's
// crash-loop surfaces the problem instead of the service running with
// a half-broken setup. Returns a multi-line aggregated error so
// operators see every problem in one pass.
func Validate(s *Scheme) error {
	var errs []string

	if s.Database == nil || s.Database.Password == "" {
		errs = append(errs, "database.password is required")
	}
	if s.Database != nil && s.Database.Name == "" {
		errs = append(errs, "database.name is required")
	}

	if s.Chain == nil {
		errs = append(errs, "chain block is required")
	} else {
		if s.Chain.WSURL == "" {
			errs = append(errs, "chain.ws_url is required")
		}
		if s.Chain.RPCURL == "" {
			errs = append(errs, "chain.rpc_url is required")
		}
		if !strings.HasPrefix(strings.ToLower(s.Chain.RegistryAddress), "0x") || len(s.Chain.RegistryAddress) != 42 {
			errs = append(errs, "chain.registry_address must be a 0x-prefixed 20-byte hex address")
		}
		if s.Chain.ChainID == 0 {
			errs = append(errs, "chain.chain_id is required")
		}
	}

	if s.Indexer == nil {
		errs = append(errs, "indexer block is required")
	} else {
		if s.Indexer.BackfillChunkSize == 0 || s.Indexer.BackfillChunkSize > 10_000 {
			errs = append(errs, "indexer.backfill_chunk_size must be in [1, 10000]")
		}
		if s.Indexer.StreamSubscriberBuffer <= 0 {
			errs = append(errs, "indexer.stream_subscriber_buffer must be >= 1")
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("invalid configuration:\n  - %s", strings.Join(errs, "\n  - "))
}

// ErrNotConfigured is returned when a required config block is missing
// from a struct that is supposed to be non-nil.
var ErrNotConfigured = errors.New("config block not configured")
