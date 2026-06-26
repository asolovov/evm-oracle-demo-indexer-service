// Package chainsub is the live WebSocket subscription layer. It watches
// a single chain for the configured PriceAggregators + the
// OracleRegistry, decodes logs into domain events, persists them, and
// publishes each the moment it is persisted (emit-on-ingest — no
// confirmation gate). Live-only: the sole chain operation is the WS log
// subscription.
package chainsub

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/contracts/oracleregistry"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/contracts/priceaggregator"
)

// Event signature hashes — precomputed via abigen Filterer for cheap
// log discrimination. The abigen-generated `Parse*` methods would
// also work, but a topic[0] switch avoids reflective unpacking on the
// hot path.
var (
	priceRequestedTopic  = mustEventID(priceaggregator.PriceAggregatorMetaData.ABI, "PriceRequested")
	priceFulfilledTopic  = mustEventID(priceaggregator.PriceAggregatorMetaData.ABI, "PriceFulfilled")
	assetRegisteredTopic = mustEventID(oracleregistry.OracleRegistryMetaData.ABI, "AssetRegistered")
)

// AggregatorResolver maps a 20-byte aggregator address to the asset
// it serves. Backed by an in-memory cache that chainsub keeps in
// sync with the registry events + the on-startup enumeration.
type AggregatorResolver interface {
	AssetIDFor(addr common.Address) (common.Hash, bool)
}

// Parser decodes a raw log into a domain Event. Returns
// (nil, ErrUnknownEvent) when topic[0] doesn't match any known event;
// callers should log + skip those rather than treat as fatal.
type Parser struct {
	registry  common.Address
	resolver  AggregatorResolver
	aggFilter *priceaggregator.PriceAggregatorFilterer
	regFilter *oracleregistry.OracleRegistryFilterer
}

// ErrUnknownEvent is returned by Parser.Parse when the log's topic[0]
// doesn't match any event signature the indexer cares about.
var ErrUnknownEvent = errors.New("unknown event signature")

// ErrAssetIDUnknown is returned when the parser can't resolve the
// emitting aggregator to a known asset_id (it isn't in the configured
// asset set and hasn't been seen via a live AssetRegistered). The
// caller still persists the event with empty asset_id and logs a
// warning.
var ErrAssetIDUnknown = errors.New("aggregator address not in asset mapping")

// NewParser constructs a Parser. `registry` is the OracleRegistry
// address used to discriminate registry-emitted logs from aggregator-
// emitted ones.
func NewParser(registry common.Address, resolver AggregatorResolver) (*Parser, error) {
	aggFilter, err := priceaggregator.NewPriceAggregatorFilterer(common.Address{}, nil)
	if err != nil {
		return nil, fmt.Errorf("init aggregator filterer: %w", err)
	}
	regFilter, err := oracleregistry.NewOracleRegistryFilterer(common.Address{}, nil)
	if err != nil {
		return nil, fmt.Errorf("init registry filterer: %w", err)
	}
	return &Parser{
		registry:  registry,
		resolver:  resolver,
		aggFilter: aggFilter,
		regFilter: regFilter,
	}, nil
}

// Parse decodes log into a domain Event with confirmations=0. The
// returned Event's asset_id is denormalised (resolved from the emit
// address for PriceRequested/PriceFulfilled, taken from the payload
// for AssetRegistered).
func (p *Parser) Parse(log types.Log) (*models.Event, error) {
	if len(log.Topics) == 0 {
		return nil, ErrUnknownEvent
	}
	topic := log.Topics[0]

	switch topic {
	case priceRequestedTopic:
		return p.parsePriceRequested(log)
	case priceFulfilledTopic:
		return p.parsePriceFulfilled(log)
	case assetRegisteredTopic:
		return p.parseAssetRegistered(log)
	default:
		return nil, ErrUnknownEvent
	}
}

func (p *Parser) parsePriceRequested(log types.Log) (*models.Event, error) {
	raw, err := p.aggFilter.ParsePriceRequested(log)
	if err != nil {
		return nil, fmt.Errorf("decode PriceRequested: %w", err)
	}
	asset, ok := p.resolver.AssetIDFor(log.Address)
	if !ok {
		// Persist with empty asset_id; surface to the caller as a
		// warning (the aggregator isn't in the configured asset set).
		return p.buildPriceRequested(log, raw, common.Hash{}, ErrAssetIDUnknown)
	}
	return p.buildPriceRequested(log, raw, asset, nil)
}

func (p *Parser) buildPriceRequested(log types.Log, raw *priceaggregator.PriceAggregatorPriceRequested, asset common.Hash, warn error) (*models.Event, error) {
	e := &models.Event{
		Kind:            models.EventKindPriceRequested,
		ContractAddress: log.Address,
		TxHash:          log.TxHash,
		BlockHash:       log.BlockHash,
		BlockNumber:     log.BlockNumber,
		LogIndex:        uint32(log.Index), //nolint:gosec // log index per block is bounded by gas limit; never exceeds uint32.
		ObservedAt:      time.Now().UTC(),
		AssetID:         asset,
		ReqID:           new(big.Int).Set(raw.ReqId),
		PriceRequested: &models.PriceRequestedPayload{
			ReqID:     new(big.Int).Set(raw.ReqId),
			AssetID:   asset,
			Requester: raw.Requester,
		},
	}
	return e, warn
}

func (p *Parser) parsePriceFulfilled(log types.Log) (*models.Event, error) {
	raw, err := p.aggFilter.ParsePriceFulfilled(log)
	if err != nil {
		return nil, fmt.Errorf("decode PriceFulfilled: %w", err)
	}
	asset, ok := p.resolver.AssetIDFor(log.Address)
	var warn error
	if !ok {
		warn = ErrAssetIDUnknown
	}
	e := &models.Event{
		Kind:            models.EventKindPriceFulfilled,
		ContractAddress: log.Address,
		TxHash:          log.TxHash,
		BlockHash:       log.BlockHash,
		BlockNumber:     log.BlockNumber,
		LogIndex:        uint32(log.Index), //nolint:gosec // log index per block is bounded by gas limit; never exceeds uint32.
		ObservedAt:      time.Now().UTC(),
		AssetID:         asset,
		ReqID:           new(big.Int).Set(raw.ReqId),
		PriceFulfilled: &models.PriceFulfilledPayload{
			ReqID:     new(big.Int).Set(raw.ReqId),
			AssetID:   asset,
			Price:     new(big.Int).Set(raw.Price),
			Timestamp: new(big.Int).Set(raw.Timestamp),
			RoundID:   nil, // backfilled by chainsub via latestRoundData when available.
		},
	}
	return e, warn
}

func (p *Parser) parseAssetRegistered(log types.Log) (*models.Event, error) {
	if log.Address != p.registry {
		// The registry-event topic is unique to OracleRegistry, but
		// be defensive: refuse to decode AssetRegistered emitted by
		// any other contract.
		return nil, fmt.Errorf("AssetRegistered from non-registry address %s", log.Address.Hex())
	}
	raw, err := p.regFilter.ParseAssetRegistered(log)
	if err != nil {
		return nil, fmt.Errorf("decode AssetRegistered: %w", err)
	}
	asset := common.BytesToHash(raw.AssetId[:])
	return &models.Event{
		Kind:            models.EventKindAssetRegistered,
		ContractAddress: log.Address,
		TxHash:          log.TxHash,
		BlockHash:       log.BlockHash,
		BlockNumber:     log.BlockNumber,
		LogIndex:        uint32(log.Index), //nolint:gosec // log index per block is bounded by gas limit; never exceeds uint32.
		ObservedAt:      time.Now().UTC(),
		AssetID:         asset,
		AssetRegistered: &models.AssetRegisteredPayload{
			AssetID:    asset,
			Aggregator: raw.Aggregator,
		},
	}, nil
}

// mustEventID computes an event signature hash from a JSON ABI string
// at init time. Panics if the event is missing — caught at init, not
// per-log.
func mustEventID(abiJSON, name string) common.Hash {
	parsed, err := abiFromJSON(abiJSON)
	if err != nil {
		panic(fmt.Errorf("chainsub: parse ABI for %s: %w", name, err))
	}
	ev, ok := parsed.Events[name]
	if !ok {
		panic(fmt.Errorf("chainsub: ABI missing event %s", name))
	}
	return ev.ID
}
