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
