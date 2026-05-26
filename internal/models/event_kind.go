// Package models holds the indexer-service domain types. Per
// architecture rule 3 every conversion method (proto <-> model, DB
// row <-> model) lives here so transport and storage layers stay free
// of conversion logic.
package models

import (
	"fmt"

	indexerv1 "github.com/asolovov/evm-oracle-demo-indexer-service/internal/genproto/indexer/v1"
)

// EventKind is the int-backed discriminator that mirrors
// `indexer.v1.EventKind` (proto) and the `kind` column in `events`
// (Postgres). The shape mirrors the template's `user_status.go`
// (architecture rule 3 reference).
type EventKind int

const (
	// EventKindUnknown is the zero value; never persisted, only used
	// as a sentinel for parse failures.
	EventKindUnknown EventKind = iota
	// EventKindPriceRequested = PriceAggregator.PriceRequested.
	EventKindPriceRequested
	// EventKindPriceFulfilled = PriceAggregator.PriceFulfilled.
	EventKindPriceFulfilled
	// EventKindAssetRegistered = OracleRegistry.AssetRegistered.
	EventKindAssetRegistered
)

// eventKindNames is the canonical lookup table. The string form is
// also what gets persisted in the `events.kind` column.
var eventKindNames = map[EventKind]string{
	EventKindUnknown:         "UNKNOWN",
	EventKindPriceRequested:  "PRICE_REQUESTED",
	EventKindPriceFulfilled:  "PRICE_FULFILLED",
	EventKindAssetRegistered: "ASSET_REGISTERED",
}

var eventKindFromString = map[string]EventKind{
	"PRICE_REQUESTED":  EventKindPriceRequested,
	"PRICE_FULFILLED":  EventKindPriceFulfilled,
	"ASSET_REGISTERED": EventKindAssetRegistered,
}

// String returns the canonical string form. Used in DB storage and
// log/metric labels.
func (k EventKind) String() string {
	if name, ok := eventKindNames[k]; ok {
		return name
	}
	return "UNKNOWN"
}

// IsValid reports whether k is one of the three real kinds.
func (k EventKind) IsValid() bool {
	return k == EventKindPriceRequested ||
		k == EventKindPriceFulfilled ||
		k == EventKindAssetRegistered
}

// ParseEventKind converts the canonical string form back to a kind.
// Returns EventKindUnknown + a non-nil error for unrecognised input.
func ParseEventKind(s string) (EventKind, error) {
	if k, ok := eventKindFromString[s]; ok {
		return k, nil
	}
	return EventKindUnknown, fmt.Errorf("unknown EventKind %q", s)
}

// ToProto converts the domain kind to the wire enum.
func (k EventKind) ToProto() indexerv1.EventKind {
	switch k {
	case EventKindPriceRequested:
		return indexerv1.EventKind_EVENT_KIND_PRICE_REQUESTED
	case EventKindPriceFulfilled:
		return indexerv1.EventKind_EVENT_KIND_PRICE_FULFILLED
	case EventKindAssetRegistered:
		return indexerv1.EventKind_EVENT_KIND_ASSET_REGISTERED
	default:
		return indexerv1.EventKind_EVENT_KIND_UNSPECIFIED
	}
}

// EventKindFromProto converts a proto enum back to the domain kind.
// Unspecified maps to EventKindUnknown — callers should treat that as
// an unfiltered request (per the proto's "Empty = all kinds" rule).
func EventKindFromProto(p indexerv1.EventKind) EventKind {
	switch p {
	case indexerv1.EventKind_EVENT_KIND_PRICE_REQUESTED:
		return EventKindPriceRequested
	case indexerv1.EventKind_EVENT_KIND_PRICE_FULFILLED:
		return EventKindPriceFulfilled
	case indexerv1.EventKind_EVENT_KIND_ASSET_REGISTERED:
		return EventKindAssetRegistered
	default:
		return EventKindUnknown
	}
}
