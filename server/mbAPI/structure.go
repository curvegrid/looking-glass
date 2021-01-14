// Copyright (c) 2021 Curvegrid Inc.

package mbAPI

import (
	"encoding/json"

	"github.com/curvegrid/looking-glass/server/blockchain"
)

// APICallResult represents the base structure of a Multibaas API call result
type APICallResult struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result,omitempty"`
}

// JSONPOSTMethodArgs is the arguments used to call a contract method
// with Multibaas API
type JSONPOSTMethodArgs struct {
	Signature        string        `json:"signature"`
	Args             []interface{} `json:"args"` // Delayed so that we type-check and convert once
	Preview          interface{}   `json:"preview,omitempty"`
	ContractOverride bool          `json:"contractOverride"`
	blockchain.TransactionArgs
}
