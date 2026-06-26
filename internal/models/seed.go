package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// seedBlock is the block_number stamped on bootstrap AssetRegistered
// rows. It is 1 (not 0) on purpose: StreamEvents replay filters
// `block_number >= from_block` with from_block >= 1, so a block-0 seed
// would be invisible to every replay path. Block 1 sorts before any
// real event and stays replay-reachable.
const seedBlock = 1

// NewAssetRegisteredSeed builds a bootstrap AssetRegistered event for a
// configured (assetID, aggregator) pair. The indexer is live-only and
// never observes the original on-chain AssetRegistered logs (they're
// historical), so on startup it seeds them into the events table from
// config — the API reads the asset set from these.
//
// The on-chain identity (tx_hash/block) is not known from config, so we
// synthesize a DETERMINISTIC tx_hash = keccak("asset-registered-seed" |
// aggregator | assetID) with log_index 0. Deterministic ⇒ the startup
// upsert is idempotent (the UNIQUE (tx_hash, log_index) constraint makes
// re-runs no-ops). These rows are bootstrap metadata, not real chain
// observations — documented as such. block_number is seedBlock.
func NewAssetRegisteredSeed(assetID common.Hash, aggregator common.Address) *Event {
	seed := append([]byte("asset-registered-seed"), aggregator.Bytes()...)
	seed = append(seed, assetID.Bytes()...)
	txHash := crypto.Keccak256Hash(seed)

	return &Event{
		Kind:            EventKindAssetRegistered,
		ContractAddress: aggregator,
		TxHash:          txHash,
		BlockHash:       common.Hash{},
		BlockNumber:     seedBlock,
		LogIndex:        0,
		ObservedAt:      time.Now().UTC(),
		AssetID:         assetID,
		AssetRegistered: &AssetRegisteredPayload{
			AssetID:    assetID,
			Aggregator: aggregator,
		},
	}
}
