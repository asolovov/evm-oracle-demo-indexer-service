package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewAssetRegisteredSeed builds a bootstrap AssetRegistered event for a
// configured (assetID, aggregator) pair. The indexer is live-only and
// never observes the original on-chain AssetRegistered logs (they're
// historical), so on startup it seeds them into the events table from
// config — the API reads the asset set from these.
//
// The on-chain identity (tx_hash/block) is not known from config, so we
// synthesize a DETERMINISTIC tx_hash = keccak("asset-registered-seed" |
// aggregator | assetID) with log_index 0 and block 0. Deterministic ⇒
// the startup upsert is idempotent (the UNIQUE (tx_hash, log_index)
// constraint makes re-runs no-ops). These rows are bootstrap metadata,
// not real chain observations — documented as such.
func NewAssetRegisteredSeed(assetID common.Hash, aggregator common.Address) *Event {
	seed := append([]byte("asset-registered-seed"), aggregator.Bytes()...)
	seed = append(seed, assetID.Bytes()...)
	txHash := crypto.Keccak256Hash(seed)

	return &Event{
		Kind:            EventKindAssetRegistered,
		ContractAddress: aggregator,
		TxHash:          txHash,
		BlockHash:       common.Hash{},
		BlockNumber:     0,
		LogIndex:        0,
		ObservedAt:      time.Now().UTC(),
		AssetID:         assetID,
		AssetRegistered: &AssetRegisteredPayload{
			AssetID:    assetID,
			Aggregator: aggregator,
		},
	}
}
