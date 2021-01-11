// Copyright (c) 2021 Curvegrid Inc.

package mbAPI

import (
	"encoding/json"

	"github.com/curvegrid/looking-glass/server/blockchain"
)

type APICallResult struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type JSONPOSTMethodArgs struct {
	Signature        string            `json:"signature"`
	Args             []json.RawMessage `json:"args"` // Delayed so that we type-check and convert once
	Preview          json.RawMessage   `json:"preview,omitempty"`
	ContractOverride bool              `json:"contractOverride"`
	blockchain.TransactionArgs
}
