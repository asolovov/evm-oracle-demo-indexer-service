package repository

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
)

func TestEncodePayloadsPriceRequested(t *testing.T) {
	e := &models.Event{
		Kind: models.EventKindPriceRequested,
		PriceRequested: &models.PriceRequestedPayload{
			ReqID:     big.NewInt(7),
			AssetID:   common.HexToHash("0xabc"),
			Requester: common.HexToAddress("0xCEF4FE1Ca9071f4ED4BAd6c1087CeB08838a983E"),
		},
	}
	out, err := encodePayloads(e)
	if err != nil {
		t.Fatalf("encodePayloads: %v", err)
	}
	if len(out.priceFulfilled) != 0 || len(out.assetRegistered) != 0 {
		t.Error("non-target payloads should be nil")
	}
	var decoded jsonPriceRequested
	if err := json.Unmarshal(out.priceRequested, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.ReqID != "7" {
		t.Errorf("ReqID = %q want %q", decoded.ReqID, "7")
	}
	if decoded.Requester != "0xcef4fe1ca9071f4ed4bad6c1087ceb08838a983e" {
		t.Errorf("requester not lower-cased: %q", decoded.Requester)
	}
}

func TestEncodePayloadsRoundIDOmittedWhenEmpty(t *testing.T) {
	e := &models.Event{
		Kind: models.EventKindPriceFulfilled,
		PriceFulfilled: &models.PriceFulfilledPayload{
			ReqID:     big.NewInt(1),
			AssetID:   common.HexToHash("0x01"),
			Price:     big.NewInt(100),
			Timestamp: big.NewInt(1),
			RoundID:   nil,
		},
	}
	out, err := encodePayloads(e)
	if err != nil {
		t.Fatalf("encodePayloads: %v", err)
	}
	if string(out.priceFulfilled) == "" || !jsonContains(out.priceFulfilled, `"round_id"`) {
		t.Logf("payload: %s", string(out.priceFulfilled))
	}
	// Verify RoundID is omitted (omitempty) when zero.
	if jsonContains(out.priceFulfilled, `"round_id":""`) {
		t.Errorf("RoundID should be omitted entirely, not emitted as empty string")
	}
}

func TestEncodePayloadsNilSubpayloadRejected(t *testing.T) {
	e := &models.Event{Kind: models.EventKindPriceRequested}
	if _, err := encodePayloads(e); err == nil {
		t.Error("expected error for nil subpayload")
	}
}

func TestEncodePayloadsUnknownKindRejected(t *testing.T) {
	e := &models.Event{Kind: models.EventKindUnknown}
	if _, err := encodePayloads(e); err == nil {
		t.Error("expected error for unknown kind")
	}
}

func TestBigIntJSONHelpers(t *testing.T) {
	if bigIntJSON(nil) != "0" {
		t.Errorf("bigIntJSON(nil) should be '0'")
	}
	if bigIntJSONOptional(nil) != "" {
		t.Errorf("bigIntJSONOptional(nil) should be empty")
	}
	if bigIntJSONOptional(big.NewInt(0)) != "" {
		t.Errorf("bigIntJSONOptional(0) should be empty")
	}
	if bigIntJSONOptional(big.NewInt(42)) != "42" {
		t.Errorf("bigIntJSONOptional(42) should be '42'")
	}
}

func TestParseBigInt(t *testing.T) {
	if got := parseBigInt("123"); got == nil || got.Int64() != 123 {
		t.Errorf("parseBigInt('123') = %v", got)
	}
	if got := parseBigIntOptional(""); got != nil {
		t.Errorf("parseBigIntOptional('') = %v, want nil", got)
	}
}

func jsonContains(b []byte, needle string) bool {
	return len(b) > 0 && string(b) != "null" && containsString(string(b), needle)
}

func containsString(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
