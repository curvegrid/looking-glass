package mbAPI

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionArgs struct {
	From          *common.Address `json:"from,omitempty"`
	To            *common.Address `json:"to,omitempty"`
	GasLimit      uint64          `json:"gasLimit,omitempty"`
	GasPrice      *big.Int        `json:"gasPrice,omitempty"`
	Value         *big.Int        `json:"value,omitempty"`
	Nonce         *big.Int        `json:"nonce,omitempty"`
	SignAndSubmit bool            `json:"signAndSubmit,omitempty"`
	Signer        *common.Address
}

type JSONPOSTMethodArgs struct {
	Signature        string            `json:"signature"`
	Args             []json.RawMessage `json:"args"` // Delayed so that we type-check and convert once
	Preview          json.RawMessage   `json:"preview,omitempty"`
	ContractOverride bool              `json:"contractOverride"`
	TransactionArgs
}

type APICallResult struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result,omitempty"`
}
