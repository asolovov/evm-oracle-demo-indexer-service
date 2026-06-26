package models

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestNewAssetRegisteredSeed(t *testing.T) {
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	agg := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")

	e := NewAssetRegisteredSeed(asset, agg)
	if e.Kind != EventKindAssetRegistered {
		t.Fatalf("kind = %v, want ASSET_REGISTERED", e.Kind)
	}
	// Must be block >= 1 so StreamEvents replay (from_block >= 1) reaches it.
	if e.BlockNumber < 1 {
		t.Errorf("seed block_number = %d, want >= 1 (replay-reachable)", e.BlockNumber)
	}
	if e.AssetRegistered == nil || e.AssetRegistered.AssetID != asset || e.AssetRegistered.Aggregator != agg {
		t.Errorf("payload mismatch: %+v", e.AssetRegistered)
	}
	// Deterministic: same inputs -> same synthetic tx_hash (idempotent seed).
	if NewAssetRegisteredSeed(asset, agg).TxHash != e.TxHash {
		t.Error("tx_hash not deterministic across calls")
	}
	// Distinct assets -> distinct tx_hash (no UNIQUE collision).
	other := common.HexToHash("0x98da2c5e4c6b1db946694570273b859a6e4083ccc8faa155edfc4c54eb3cfd73")
	if NewAssetRegisteredSeed(other, agg).TxHash == e.TxHash {
		t.Error("different assets produced the same tx_hash")
	}
}
