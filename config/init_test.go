package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func resetViper(t *testing.T) {
	t.Helper()
	viper.Reset()
	setDefaults()
	t.Cleanup(viper.Reset)
}

func TestDefaultsRegistered(t *testing.T) {
	resetViper(t)

	cases := []struct {
		key  string
		want any
	}{
		{"env", "prod"},
		{"database.driver", "postgres"},
		{"database.name", "evm_indexer"},
		{"database.user", "indexer_user"},
		{"grpc.port", 9090},
		{"grpc.reflection", true},
		{"healthz.port", 8080},
		{"chain.name", "ethereum-sepolia"},
		{"chain.chain_id", uint64(11155111)},
		{"chain.registry_address", defaultRegistryAddress},
		{"indexer.stream_subscriber_buffer", 256},
		{"telemetry.log_level", "info"},
		{"telemetry.log_format", "json"},
	}

	for _, c := range cases {
		got := viper.Get(c.key)
		// Compare via fmt so int vs uint typing artifacts in viper don't trip the test.
		if fmt.Sprint(got) != fmt.Sprint(c.want) {
			t.Errorf("key %q: want %v, got %v", c.key, c.want, got)
		}
	}
}

func TestValidate_AcceptsCompleteScheme(t *testing.T) {
	s := &Scheme{
		Database: &DatabaseConfig{Name: "evm_indexer", Password: "x"},
		Chain: &ChainConfig{
			Name:            "ethereum-sepolia",
			ChainID:         11155111,
			WSURL:           "ws://node",
			RegistryAddress: "0x89a6c12a403733c6a817472cec46a530581cb7ef",
		},
		Indexer: &IndexerConfig{
			StreamSubscriberBuffer: 256,
			Assets:                 defaultAssets(),
		},
	}
	if err := Validate(s); err != nil {
		t.Fatalf("Validate returned %v, want nil", err)
	}
}

func TestApplyEnvOverrides_DefaultsAssets(t *testing.T) {
	t.Setenv("INDEXER_ASSETS", "")
	s := &Scheme{Indexer: &IndexerConfig{StreamSubscriberBuffer: 256}}
	if err := ApplyEnvOverrides(s); err != nil {
		t.Fatalf("ApplyEnvOverrides: %v", err)
	}
	if len(s.Indexer.Assets) != 10 {
		t.Errorf("expected 10 default assets, got %d", len(s.Indexer.Assets))
	}
}

func TestApplyEnvOverrides_JSON(t *testing.T) {
	t.Setenv("INDEXER_ASSETS", `[{"symbol":"FOO","asset_id":"0xa1","aggregator":"0xb2"}]`)
	s := &Scheme{Indexer: &IndexerConfig{StreamSubscriberBuffer: 256}}
	if err := ApplyEnvOverrides(s); err != nil {
		t.Fatalf("ApplyEnvOverrides: %v", err)
	}
	if len(s.Indexer.Assets) != 1 {
		t.Fatalf("env override not applied: %+v", s.Indexer.Assets)
	}
	// Assert ALL fields decode — asset_id/aggregator need json tags
	// (encoding/json ignores mapstructure tags). A missing tag silently
	// leaves them empty.
	got := s.Indexer.Assets[0]
	if got.Symbol != "FOO" || got.AssetID != "0xa1" || got.Aggregator != "0xb2" {
		t.Errorf("env override decoded incomplete: %+v", got)
	}
}

func TestValidate_RejectsMissingFields(t *testing.T) {
	s := &Scheme{}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected an error for an empty scheme")
	}
	for _, expect := range []string{
		"database.password is required",
		"chain block is required",
		"indexer block is required",
	} {
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("missing %q in error: %v", expect, err)
		}
	}
}

func TestValidate_RejectsBadRegistryAddress(t *testing.T) {
	s := &Scheme{
		Database: &DatabaseConfig{Name: "evm_indexer", Password: "x"},
		Chain: &ChainConfig{
			Name:            "x",
			ChainID:         1,
			WSURL:           "ws://x",
			RegistryAddress: "not-an-address",
		},
		Indexer: &IndexerConfig{
			StreamSubscriberBuffer: 256,
			Assets:                 defaultAssets(),
		},
	}
	err := Validate(s)
	if err == nil || !strings.Contains(err.Error(), "registry_address") {
		t.Fatalf("expected registry_address error, got %v", err)
	}
}
