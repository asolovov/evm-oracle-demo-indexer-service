package chainsub

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// abiFromJSON is a thin wrapper around abi.JSON so the parser file
// stays focused on event-decoding logic.
func abiFromJSON(s string) (*abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// shortHash trims a 0x-prefixed hex string to "0xabcd…ef01" so log
// lines stay scannable without losing identity. Inputs shorter than
// 14 characters pass through unchanged.
func shortHash(h string) string {
	if len(h) < 14 {
		return h
	}
	return h[:6] + "…" + h[len(h)-4:]
}
