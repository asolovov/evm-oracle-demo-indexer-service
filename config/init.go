// Package config defines indexer-service configuration defaults and schema.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// defaultRegistryAddress is the deployed OracleRegistry on Ethereum
// Sepolia. defaultAssets() are its 10 aggregators (asset_id +
// aggregator), mirrored from
// evm-oracle-demo-contracts/deployments/ethereum-sepolia/addresses.json.
// These are the indexer's source of truth for the deployment — it no
// longer reads the registry on chain.
const defaultRegistryAddress = "0x89a6c12a403733c6a817472cec46a530581cb7ef"

// assetsEnvVar, when set to a JSON array of {symbol,asset_id,aggregator}
// objects, replaces the baked default asset set (for a different deploy).
const assetsEnvVar = "INDEXER_ASSETS"

func defaultAssets() []AssetConfig {
	return []AssetConfig{
		{Symbol: "WETH", AssetID: "0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8", Aggregator: "0x075be31662c2548c4e940d7e769c328a34dcb281"},
		{Symbol: "WBTC", AssetID: "0x98da2c5e4c6b1db946694570273b859a6e4083ccc8faa155edfc4c54eb3cfd73", Aggregator: "0xf8ad3a2505eece7ad276db038c7c56930bd436e4"},
		{Symbol: "LINK", AssetID: "0x921a3539bcb764c889432630877414523e7fbca00c211bc787aeae69e2e3a779", Aggregator: "0xecc43e6ec38ce135b81ae8042df96eef55915d14"},
		{Symbol: "UNI", AssetID: "0xfba01d52a7cd84480d0573725899486a0b5e55c20ff45d6628874349375d1650", Aggregator: "0x69d16087172f404925ffc61c0ac25c608ff215b4"},
		{Symbol: "AAVE", AssetID: "0xde46fbfa339d54cd65b79d8320a7a53c78177565c2aaf4c8b13eed7865e7cfc8", Aggregator: "0xa011fa0757b5d2a9a4c73cfb4647c29d96da7a2f"},
		{Symbol: "XAU", AssetID: "0x7c687a3207cd9c05b4b11d8dd7ac337919c2200102d72989a597ebc5afcf180b", Aggregator: "0x61125ef037305e4b81c5e5a864225860f455d318"},
		{Symbol: "XAG", AssetID: "0x5ccc5c04130d272bf07d6e066f4cae40cfc0313643d815db3e17af00e6ebf601", Aggregator: "0x4e05cc443cbcd5425b5b7c7df124101ad70b8b02"},
		{Symbol: "SPX", AssetID: "0x1308465f1da3a6702b88abc29db16011bdb6f6a7cb404fee1daa81f8da9d9972", Aggregator: "0x3fa9e3fd3e5e70f26ccf4b67825489276f9cbb27"},
		{Symbol: "WTI", AssetID: "0x1f29567db4e0c1628fa0f8675c031b615246dd0dd3de399fdf8b5aec1829181d", Aggregator: "0x70131a2612682f7d56a2a30010075e8f0e9d8eca"},
		{Symbol: "HG", AssetID: "0x7f1edccb34ff65dc749f950e76926ca09253b4f2e87cc2a946d4ecaa2716decf", Aggregator: "0x87249f3aeb58c46be3f5edd1d5071ee76d816900"},
	}
}

//nolint:gochecknoinits // configuration defaults are registered at package load.
func init() {
	setDefaults()
}

// ApplyEnvOverrides applies overrides that viper can't express through
// SetDefault — currently the JSON asset-set override. Call after
// viper.Unmarshal, before Validate.
func ApplyEnvOverrides(s *Scheme) error {
	if s.Indexer == nil {
		return nil
	}
	if len(s.Indexer.Assets) == 0 {
		s.Indexer.Assets = defaultAssets()
	}
	raw := strings.TrimSpace(os.Getenv(assetsEnvVar))
	if raw == "" {
		return nil
	}
	var assets []AssetConfig
	if err := json.Unmarshal([]byte(raw), &assets); err != nil {
		return fmt.Errorf("parse %s: %w", assetsEnvVar, err)
	}
	s.Indexer.Assets = assets
	return nil
}

func isHexAddress(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "0x") && len(s) == 42
}

func isHexHash(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "0x") && len(s) == 66
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

	// Chain — ws_url must be set via env at deploy time.
	viper.SetDefault("chain.name", "ethereum-sepolia")
	viper.SetDefault("chain.chain_id", uint64(11155111))
	viper.SetDefault("chain.ws_url", "")
	viper.SetDefault("chain.registry_address", defaultRegistryAddress)

	// Indexer knobs.
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
		if !isHexAddress(s.Chain.RegistryAddress) {
			errs = append(errs, "chain.registry_address must be a 0x-prefixed 20-byte hex address")
		}
		if s.Chain.ChainID == 0 {
			errs = append(errs, "chain.chain_id is required")
		}
	}

	if s.Indexer == nil {
		errs = append(errs, "indexer block is required")
	} else {
		if s.Indexer.StreamSubscriberBuffer <= 0 {
			errs = append(errs, "indexer.stream_subscriber_buffer must be >= 1")
		}
		if len(s.Indexer.Assets) == 0 {
			errs = append(errs, "indexer.assets must list at least one asset")
		}
		seenAgg := make(map[string]struct{}, len(s.Indexer.Assets))
		for i, a := range s.Indexer.Assets {
			if !isHexHash(a.AssetID) {
				errs = append(errs, fmt.Sprintf("indexer.assets[%d].asset_id must be a 0x-prefixed 32-byte hex", i))
			}
			if !isHexAddress(a.Aggregator) {
				errs = append(errs, fmt.Sprintf("indexer.assets[%d].aggregator must be a 0x-prefixed 20-byte hex address", i))
				continue
			}
			// One aggregator -> one asset. A duplicate would silently
			// overwrite in the in-memory mapping (last wins), mislabeling
			// that aggregator's price logs.
			key := strings.ToLower(a.Aggregator)
			if _, dup := seenAgg[key]; dup {
				errs = append(errs, fmt.Sprintf("indexer.assets[%d].aggregator %s is listed more than once", i, a.Aggregator))
			}
			seenAgg[key] = struct{}{}
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
