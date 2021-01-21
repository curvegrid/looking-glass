// Copyright (c) 2021 Curvegrid Inc.

package blockchain

// TransactionArgs holds the special/reserved word values for posting transactions to a chain
type TransactionArgs struct {
	From          *Address `json:"from,omitempty"`
	To            *Address `json:"to,omitempty"`
	GasLimit      uint64   `json:"gasLimit,omitempty"`
	GasPrice      *Number  `json:"gasPrice,omitempty"`
	Value         *Number  `json:"value,omitempty"`
	Nonce         *Number  `json:"nonce,omitempty"`
	SignAndSubmit bool     `json:"signAndSubmit,omitempty"`
}
