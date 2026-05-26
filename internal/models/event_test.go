package models

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
)

func sampleEvent(kind EventKind) *Event {
	e := &Event{
		Kind:            kind,
		ContractAddress: common.HexToAddress("0x075be31662c2548c4e940d7e769c328a34dcb281"),
		TxHash:          common.HexToHash("0xaa"),
		BlockHash:       common.HexToHash("0xbb"),
		BlockNumber:     42,
		LogIndex:        3,
		ObservedAt:      time.Unix(1700000000, 0).UTC(),
		Confirmations:   5,
		AssetID:         common.HexToHash("0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8"),
		ReqID:           big.NewInt(7),
	}
	switch kind {
	case EventKindPriceRequested:
		e.PriceRequested = &PriceRequestedPayload{
			ReqID:     big.NewInt(7),
			AssetID:   e.AssetID,
			Requester: common.HexToAddress("0xCEF4Fe1CA9071f4ED4BAd6C1087CEb08838a983E"),
		}
	case EventKindPriceFulfilled:
		e.PriceFulfilled = &PriceFulfilledPayload{
			ReqID:     big.NewInt(7),
			AssetID:   e.AssetID,
			Price:     big.NewInt(345020000000),
			Timestamp: big.NewInt(1700000000),
			RoundID:   nil, // optional
		}
	case EventKindAssetRegistered:
		e.AssetRegistered = &AssetRegisteredPayload{
			AssetID:    e.AssetID,
			Aggregator: e.ContractAddress,
		}
	}
	return e
}

func TestEventToProtoPriceRequested(t *testing.T) {
	e := sampleEvent(EventKindPriceRequested)
	p, err := e.ToProto()
	if err != nil {
		t.Fatalf("ToProto: %v", err)
	}
	if p.Kind != indexerv1.EventKind_EVENT_KIND_PRICE_REQUESTED {
		t.Errorf("kind = %v, want PRICE_REQUESTED", p.Kind)
	}
	pr := p.GetPriceRequested()
	if pr == nil {
		t.Fatal("PriceRequested payload missing in proto")
	}
	if pr.ReqId != "7" || pr.AssetId != "0x0f8a193ff464434486c0daf7db2a895884365d2bc84ba47a68fcf89c1b14b5b8" {
		t.Errorf("unexpected proto fields: %+v", pr)
	}
	if p.Meta.Confirmations != 5 || p.Meta.BlockNumber != 42 {
		t.Errorf("meta lost confirmations/block_number: %+v", p.Meta)
	}
}

func TestEventToProtoPriceFulfilledEmptyRoundID(t *testing.T) {
	e := sampleEvent(EventKindPriceFulfilled)
	p, err := e.ToProto()
	if err != nil {
		t.Fatalf("ToProto: %v", err)
	}
	pf := p.GetPriceFulfilled()
	if pf == nil {
		t.Fatal("PriceFulfilled payload missing")
	}
	if pf.RoundId != "" {
		t.Errorf("expected empty RoundId, got %q", pf.RoundId)
	}
	if pf.Price != "345020000000" {
		t.Errorf("Price mismatch: %q", pf.Price)
	}
}

func TestEventToProtoAssetRegistered(t *testing.T) {
	e := sampleEvent(EventKindAssetRegistered)
	p, err := e.ToProto()
	if err != nil {
		t.Fatalf("ToProto: %v", err)
	}
	ar := p.GetAssetRegistered()
	if ar == nil {
		t.Fatal("AssetRegistered payload missing")
	}
	if ar.AssetId == "" || ar.Aggregator == "" {
		t.Errorf("unexpected proto fields: %+v", ar)
	}
}

func TestEventToProtoMissingPayloadFails(t *testing.T) {
	e := &Event{Kind: EventKindPriceRequested}
	_, err := e.ToProto()
	if !errors.Is(err, ErrPayloadMissing) {
		t.Errorf("want ErrPayloadMissing, got %v", err)
	}
}
