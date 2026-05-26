package models

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrPayloadMissing is returned when an Event has Kind set but no
// matching typed payload populated. Indicates a programming error
// upstream of persistence.
var ErrPayloadMissing = errors.New("event payload missing for declared kind")

// PriceRequestedPayload is the decoded shape of a PriceRequested log.
// `AssetID` is resolved from the emitting contract's address via the
// registry mapping — the on-chain event itself only carries reqId and
// requester.
type PriceRequestedPayload struct {
	ReqID     *big.Int       // uint256 reqId
	AssetID   common.Hash    // bytes32 — resolved from contract addr
	Requester common.Address // 20-byte requester
}

// PriceFulfilledPayload is the decoded shape of a PriceFulfilled log.
//
//	event PriceFulfilled(uint256 indexed reqId, int256 price, uint256 timestamp)
//
// `AssetID` is resolved from the emitting contract's address; `RoundID`
// is OPTIONAL and backfilled from latestRoundData() at the persist
// block by the chainsub when available, otherwise left as the zero
// big.Int (mapped to empty string in the proto).
type PriceFulfilledPayload struct {
	ReqID     *big.Int    // uint256 reqId
	AssetID   common.Hash // bytes32 — resolved from contract addr
	Price     *big.Int    // int256 in 8-decimal scale
	Timestamp *big.Int    // uint256 seconds-since-epoch
	RoundID   *big.Int    // OPTIONAL — chainlink-style uint80
}

// AssetRegisteredPayload is the decoded shape of an
// OracleRegistry.AssetRegistered log.
type AssetRegisteredPayload struct {
	AssetID    common.Hash    // bytes32 asset id
	Aggregator common.Address // newly registered aggregator
}

// Event is the indexer's single domain type for an observed log.
// Exactly one of the *Payload pointer fields is populated, matching
// Kind. Persistence + proto conversion both dispatch on Kind.
type Event struct {
	ID              int64
	Kind            EventKind
	ContractAddress common.Address
	TxHash          common.Hash
	BlockHash       common.Hash
	BlockNumber     uint64
	LogIndex        uint32
	ObservedAt      time.Time
	Confirmations   uint32
	Orphaned        bool

	// Denormalised lookup columns (also stored on the DB row for
	// indexed queries — see `events.asset_id` / `events.req_id`).
	AssetID common.Hash
	ReqID   *big.Int

	PriceRequested  *PriceRequestedPayload
	PriceFulfilled  *PriceFulfilledPayload
	AssetRegistered *AssetRegisteredPayload
}

// ChainCursor mirrors the single-row `chain_cursor` table.
type ChainCursor struct {
	LastProcessedBlock uint64
	UpdatedAt          time.Time
}

// ToProto converts the domain event to the wire shape served by
// ListEvents / StreamEvents. Returns an error if Kind has no matching
// payload populated.
func (e *Event) ToProto() (*indexerv1.Event, error) {
	meta := &indexerv1.EventMeta{
		ContractAddress: strings.ToLower(e.ContractAddress.Hex()),
		TxHash:          strings.ToLower(e.TxHash.Hex()),
		BlockHash:       strings.ToLower(e.BlockHash.Hex()),
		BlockNumber:     e.BlockNumber,
		LogIndex:        e.LogIndex,
		ObservedAt:      timestamppb.New(e.ObservedAt),
		Confirmations:   e.Confirmations,
	}

	out := &indexerv1.Event{
		Meta: meta,
		Kind: e.Kind.ToProto(),
	}

	switch e.Kind {
	case EventKindPriceRequested:
		if e.PriceRequested == nil {
			return nil, fmt.Errorf("%w: PriceRequested", ErrPayloadMissing)
		}
		out.Payload = &indexerv1.Event_PriceRequested{
			PriceRequested: &indexerv1.PriceRequestedEvent{
				ReqId:     bigIntDecimalString(e.PriceRequested.ReqID),
				AssetId:   strings.ToLower(e.PriceRequested.AssetID.Hex()),
				Requester: strings.ToLower(e.PriceRequested.Requester.Hex()),
			},
		}
	case EventKindPriceFulfilled:
		if e.PriceFulfilled == nil {
			return nil, fmt.Errorf("%w: PriceFulfilled", ErrPayloadMissing)
		}
		out.Payload = &indexerv1.Event_PriceFulfilled{
			PriceFulfilled: &indexerv1.PriceFulfilledEvent{
				ReqId:     bigIntDecimalString(e.PriceFulfilled.ReqID),
				AssetId:   strings.ToLower(e.PriceFulfilled.AssetID.Hex()),
				Price:     bigIntDecimalString(e.PriceFulfilled.Price),
				Timestamp: bigIntDecimalString(e.PriceFulfilled.Timestamp),
				RoundId:   bigIntDecimalStringOrEmpty(e.PriceFulfilled.RoundID),
			},
		}
	case EventKindAssetRegistered:
		if e.AssetRegistered == nil {
			return nil, fmt.Errorf("%w: AssetRegistered", ErrPayloadMissing)
		}
		out.Payload = &indexerv1.Event_AssetRegistered{
			AssetRegistered: &indexerv1.AssetRegisteredEvent{
				AssetId:    strings.ToLower(e.AssetRegistered.AssetID.Hex()),
				Aggregator: strings.ToLower(e.AssetRegistered.Aggregator.Hex()),
			},
		}
	default:
		return nil, fmt.Errorf("%w: kind=%s", ErrPayloadMissing, e.Kind)
	}

	return out, nil
}

// bigIntDecimalString returns "0" for nil. Always-set; matches the
// proto convention that uint256 fields are non-nullable strings.
func bigIntDecimalString(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

// bigIntDecimalStringOrEmpty returns "" for nil/zero values, used for
// the optional uint80 RoundID field where the contract event does NOT
// carry round_id natively.
func bigIntDecimalStringOrEmpty(n *big.Int) string {
	if n == nil || n.Sign() == 0 {
		return ""
	}
	return n.String()
}
