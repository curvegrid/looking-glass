package mbAPI

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionArgs struct {
	From          *common.Address `json:"from,omitempty"`
	To            *common.Address `json:"to,omitempty"`
	GasLimit      uint64          `json:"gasLimit,omitempty"`
	GasPrice      *Number         `json:"gasPrice,omitempty"`
	Value         *Number         `json:"value,omitempty"`
	Nonce         *Number         `json:"nonce,omitempty"`
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

// Number represents an input number that might be inputted as either int or string.
type Number struct {
	*big.Int
}

var (
	_ json.Unmarshaler = (*Number)(nil)
)

// UnmarshalInt tries to unmarshal an integer value either from
// a JSON string or number.
func unmarshalInt(b []byte) (n *big.Int, err error) {
	if err = json.Unmarshal(b, &n); err == nil {
		return n, nil
	}
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	n, ok := n.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("cannot parse \"%s\" as int", s)
	}
	return
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Number) UnmarshalJSON(b []byte) error {
	v, err := unmarshalInt(b)
	n.Int = v
	return err
}

// GetInt returns the int. Covers the null case.
func (n *Number) GetInt() *big.Int {
	if n == nil {
		return nil
	}
	return n.Int
}
