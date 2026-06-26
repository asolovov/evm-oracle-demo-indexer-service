package chainsub

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

type fakeResolver struct {
	mapping map[common.Address]common.Hash
}

func (f *fakeResolver) AssetIDFor(addr common.Address) (common.Hash, bool) {
	h, ok := f.mapping[addr]
	return h, ok
}

func TestEventTopicsMatchSolidityHashes(t *testing.T) {
	// Hashes pulled from the abigen comments / Solidity sources.
	wantPriceRequested := common.HexToHash("0x33b362cfc336cd3829a3a7896832b11a70bdbc0194a75fb7f1003919b23c8600")
	wantPriceFulfilled := common.HexToHash("0x82c08aa7285d5667568ba3b8821c82fa50deef99dc0c6b75d46fb5c7455ec22a")

	if priceRequestedTopic != wantPriceRequested {
		t.Errorf("PriceRequested topic = %s, want %s", priceRequestedTopic.Hex(), wantPriceRequested.Hex())
	}
	if priceFulfilledTopic != wantPriceFulfilled {
		t.Errorf("PriceFulfilled topic = %s, want %s", priceFulfilledTopic.Hex(), wantPriceFulfilled.Hex())
	}
	if (assetRegisteredTopic == common.Hash{}) {
		t.Error("AssetRegistered topic uninitialized")
	}
}

func TestParseUnknownTopicRejected(t *testing.T) {
	p, err := NewParser(common.HexToAddress("0x1"), &fakeResolver{})
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	_, err = p.Parse(types.Log{Topics: []common.Hash{common.HexToHash("0xdeadbeef")}})
	if !errors.Is(err, ErrUnknownEvent) {
		t.Errorf("got %v, want ErrUnknownEvent", err)
	}
}

func TestParseEmptyTopicsRejected(t *testing.T) {
	p, err := NewParser(common.HexToAddress("0x1"), &fakeResolver{})
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	if _, err := p.Parse(types.Log{}); !errors.Is(err, ErrUnknownEvent) {
		t.Errorf("got %v, want ErrUnknownEvent", err)
	}
}

// buildPriceRequestedLog constructs a real abigen-shaped Log for the
// `PriceRequested(uint256 indexed reqId, address indexed requester)`
// signature.
func buildPriceRequestedLog(contract common.Address, reqID *big.Int, requester common.Address, block uint64) types.Log {
	reqIDBytes := common.LeftPadBytes(reqID.Bytes(), 32)
	requesterBytes := common.LeftPadBytes(requester.Bytes(), 32)
	return types.Log{
		Address:     contract,
		Topics:      []common.Hash{priceRequestedTopic, common.BytesToHash(reqIDBytes), common.BytesToHash(requesterBytes)},
		Data:        nil,
		BlockNumber: block,
		BlockHash:   common.HexToHash("0xbeef"),
		TxHash:      common.HexToHash("0xfeed"),
		Index:       7,
	}
}

func TestParsePriceRequested_WithKnownAssetMapping(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")
	reqID := big.NewInt(7)

	p, err := NewParser(common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef"), &fakeResolver{
		mapping: map[common.Address]common.Hash{aggregator: asset},
	})
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	log := buildPriceRequestedLog(aggregator, reqID, requester, 42)

	evt, err := p.Parse(log)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if evt.Kind != models.EventKindPriceRequested {
		t.Errorf("kind = %v, want PRICE_REQUESTED", evt.Kind)
	}
	if evt.AssetID != asset {
		t.Errorf("AssetID = %s, want %s", evt.AssetID.Hex(), asset.Hex())
	}
	if evt.PriceRequested == nil || evt.PriceRequested.ReqID.Cmp(reqID) != 0 {
		t.Errorf("PriceRequested payload mismatch: %+v", evt.PriceRequested)
	}
	if evt.PriceRequested.Requester != requester {
		t.Errorf("Requester = %s, want %s", evt.PriceRequested.Requester.Hex(), requester.Hex())
	}
}

func TestParsePriceRequested_UnknownAggregatorWarn(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	requester := common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6c1087CEb08838a983E")
	reqID := big.NewInt(1)

	p, _ := NewParser(common.HexToAddress("0x1"), &fakeResolver{mapping: map[common.Address]common.Hash{}})
	log := buildPriceRequestedLog(aggregator, reqID, requester, 42)

	evt, err := p.Parse(log)
	if !errors.Is(err, ErrAssetIDUnknown) {
		t.Fatalf("want ErrAssetIDUnknown, got %v", err)
	}
	if evt == nil {
		t.Fatal("event should still be returned (caller persists with empty asset_id)")
	}
	if (evt.AssetID != common.Hash{}) {
		t.Errorf("AssetID should be zero when unresolved, got %s", evt.AssetID.Hex())
	}
}

// buildPriceFulfilledLog packs `PriceFulfilled(uint256 indexed reqId,
// int256 price, uint256 timestamp)`.
func buildPriceFulfilledLog(contract common.Address, reqID, price, ts *big.Int, block uint64) types.Log {
	data := append(common.LeftPadBytes(price.Bytes(), 32), common.LeftPadBytes(ts.Bytes(), 32)...)
	return types.Log{
		Address:     contract,
		Topics:      []common.Hash{priceFulfilledTopic, common.BigToHash(reqID)},
		Data:        data,
		BlockNumber: block,
		BlockHash:   common.HexToHash("0xbeef"),
		TxHash:      common.HexToHash("0xfeed"),
		Index:       3,
	}
}

func TestParsePriceFulfilled(t *testing.T) {
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")

	p, _ := NewParser(common.HexToAddress("0x1"), &fakeResolver{mapping: map[common.Address]common.Hash{aggregator: asset}})

	reqID := big.NewInt(7)
	price := big.NewInt(345020000000)
	ts := big.NewInt(1700000000)

	evt, err := p.Parse(buildPriceFulfilledLog(aggregator, reqID, price, ts, 42))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if evt.Kind != models.EventKindPriceFulfilled {
		t.Errorf("kind = %v, want PRICE_FULFILLED", evt.Kind)
	}
	if evt.PriceFulfilled.Price.Cmp(price) != 0 || evt.PriceFulfilled.Timestamp.Cmp(ts) != 0 {
		t.Errorf("price/ts mismatch: %+v", evt.PriceFulfilled)
	}
	if evt.PriceFulfilled.RoundID != nil {
		t.Errorf("RoundID should default to nil (set by latestRoundData fetch); got %v", evt.PriceFulfilled.RoundID)
	}
	if evt.AssetID != asset {
		t.Errorf("AssetID denormalisation failed: %s", evt.AssetID.Hex())
	}
}

// buildAssetRegisteredLog packs `AssetRegistered(bytes32 indexed assetId, address indexed aggregator)`.
func buildAssetRegisteredLog(registry common.Address, asset common.Hash, aggregator common.Address, block uint64) types.Log {
	return types.Log{
		Address:     registry,
		Topics:      []common.Hash{assetRegisteredTopic, asset, common.BytesToHash(common.LeftPadBytes(aggregator.Bytes(), 32))},
		Data:        nil,
		BlockNumber: block,
		BlockHash:   common.HexToHash("0xbeef"),
		TxHash:      common.HexToHash("0xfeed"),
		Index:       0,
	}
}

func TestParseAssetRegistered(t *testing.T) {
	registry := common.HexToAddress("0x89a6c12a403733c6a817472cec46a530581cb7ef")
	asset := common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8")
	aggregator := common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281")

	p, _ := NewParser(registry, &fakeResolver{mapping: map[common.Address]common.Hash{}})
	evt, err := p.Parse(buildAssetRegisteredLog(registry, asset, aggregator, 1))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if evt.Kind != models.EventKindAssetRegistered {
		t.Errorf("kind = %v", evt.Kind)
	}
	if evt.AssetRegistered == nil ||
		evt.AssetRegistered.AssetID != asset ||
		evt.AssetRegistered.Aggregator != aggregator {
		t.Errorf("payload mismatch: %+v", evt.AssetRegistered)
	}
}

func TestParseAssetRegisteredFromNonRegistryRejected(t *testing.T) {
	registry := common.HexToAddress("0x01")
	asset := common.HexToHash("0xaa")
	aggregator := common.HexToAddress("0x02")

	p, _ := NewParser(registry, &fakeResolver{})
	imposter := common.HexToAddress("0x03")
	_, err := p.Parse(buildAssetRegisteredLog(imposter, asset, aggregator, 1))
	if err == nil {
		t.Error("AssetRegistered from non-registry address should be rejected")
	}
}
