package blockchain

type TransactionArgs struct {
	From          *Address `json:"from,omitempty"`
	To            *Address `json:"to,omitempty"`
	GasLimit      uint64   `json:"gasLimit,omitempty"`
	GasPrice      *Number  `json:"gasPrice,omitempty"`
	Value         *Number  `json:"value,omitempty"`
	Nonce         *Number  `json:"nonce,omitempty"`
	SignAndSubmit bool     `json:"signAndSubmit,omitempty"`
}
