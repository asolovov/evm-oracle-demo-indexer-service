package models

import (
	"testing"

	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
)

func TestEventKindStringRoundTrip(t *testing.T) {
	for _, k := range []EventKind{EventKindPriceRequested, EventKindPriceFulfilled, EventKindAssetRegistered} {
		parsed, err := ParseEventKind(k.String())
		if err != nil {
			t.Fatalf("ParseEventKind(%q): %v", k.String(), err)
		}
		if parsed != k {
			t.Errorf("round-trip mismatch: want %v, got %v", k, parsed)
		}
		if !k.IsValid() {
			t.Errorf("IsValid() returned false for valid kind %v", k)
		}
	}
}

func TestEventKindUnknown(t *testing.T) {
	if EventKindUnknown.IsValid() {
		t.Error("EventKindUnknown should not be valid")
	}
	if EventKindUnknown.String() != "UNKNOWN" {
		t.Errorf("EventKindUnknown.String() = %q, want %q", EventKindUnknown.String(), "UNKNOWN")
	}
}

func TestParseEventKindUnknownRejected(t *testing.T) {
	if _, err := ParseEventKind("NOT_A_KIND"); err == nil {
		t.Error("expected error for unrecognized input")
	}
}

func TestEventKindProtoRoundTrip(t *testing.T) {
	pairs := []struct {
		domain EventKind
		proto  indexerv1.EventKind
	}{
		{EventKindPriceRequested, indexerv1.EventKind_EVENT_KIND_PRICE_REQUESTED},
		{EventKindPriceFulfilled, indexerv1.EventKind_EVENT_KIND_PRICE_FULFILLED},
		{EventKindAssetRegistered, indexerv1.EventKind_EVENT_KIND_ASSET_REGISTERED},
	}
	for _, p := range pairs {
		if got := p.domain.ToProto(); got != p.proto {
			t.Errorf("%v.ToProto() = %v, want %v", p.domain, got, p.proto)
		}
		if got := EventKindFromProto(p.proto); got != p.domain {
			t.Errorf("EventKindFromProto(%v) = %v, want %v", p.proto, got, p.domain)
		}
	}

	if got := EventKindFromProto(indexerv1.EventKind_EVENT_KIND_UNSPECIFIED); got != EventKindUnknown {
		t.Errorf("EventKindFromProto(UNSPECIFIED) = %v, want %v", got, EventKindUnknown)
	}
}
