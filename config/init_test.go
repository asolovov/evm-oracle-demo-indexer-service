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
		{"indexer.confirmations", uint32(5)},
		{"indexer.reorg_check_interval_sec", uint32(10)},
		{"indexer.backfill_chunk_size", uint64(1000)},
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
			Name:              "ethereum-sepolia",
			ChainID:           11155111,
			WSURL:             "ws://node",
			RPCURL:            "http://node",
			RegistryAddress:   "0x89a6c12a403733c6a817472cec46a530581cb7ef",
			BackfillFromBlock: 1,
		},
		Indexer: &IndexerConfig{
			Confirmations:          5,
			ReorgCheckIntervalSec:  10,
			BackfillChunkSize:      1000,
			StreamSubscriberBuffer: 256,
		},
	}
	if err := Validate(s); err != nil {
		t.Fatalf("Validate returned %v, want nil", err)
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
			RPCURL:          "http://x",
			RegistryAddress: "not-an-address",
		},
		Indexer: &IndexerConfig{
			Confirmations:          5,
			ReorgCheckIntervalSec:  10,
			BackfillChunkSize:      1000,
			StreamSubscriberBuffer: 256,
		},
	}
	err := Validate(s)
	if err == nil || !strings.Contains(err.Error(), "registry_address") {
		t.Fatalf("expected registry_address error, got %v", err)
	}
}
