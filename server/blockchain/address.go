// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Address wraps a geth's common.Address.
// It uses SQL interfaces and JSON Unmarshaller from common.Address
// but implements different JSON Marshaller.
type Address struct {
	common.Address
}

// Format implements (and overrides common.Address's implementation of) fmt.Formatter.
// Prints the address as if it was a hex address.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a.String())
}

// HexToAddress creates an Address from an address string and returns its pointer.
func HexToAddress(s string) *Address {
	return &Address{common.HexToAddress(s)}
}

// MarshalText returns the hex representation in EIP-55 format of an address.
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}
