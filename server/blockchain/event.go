// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type EventInformation struct {
	Name            string               `json:"name"`
	Signature       string               `json:"signature"`
	Inputs          []EventField         `json:"inputs"`
	Raw             string               `json:"rawFields"`
	Contract        *ContractInformation `json:"contract"`
	EventIndexInLog int                  `json:"indexInLog"`
}

type MethodInformation struct {
	Name            string       `json:"name"`
	Signature       string       `json:"signature"`
	Inputs          []EventField `json:"inputs"`
	FromConstructor bool         `json:"fromConstructor"`
	IsFallback      bool         `json:"isFallback"`
}

type ContractInformation struct {
	Address      *common.Address `json:"address"`
	AddressLabel string          `json:"addressLabel"`
	Name         string          `json:"name"`
	Label        string          `json:"label"`
}

type TransactionInformation struct {
	TXFrom         *common.Address      `json:"from"`
	TXData         string               `json:"txData"`
	TXHash         common.Hash          `json:"txHash"`
	TXIndexInBlock int                  `json:"txIndexInBlock"`
	BlockHash      common.Hash          `json:"blockHash"`
	BlockNumber    int64                `json:"blockNumber"`
	Contract       *ContractInformation `json:"contract"`
	Method         *MethodInformation   `json:"method"`
}

type EventField struct {
	Name   string      `json:"name"`
	Value  interface{} `json:"value"`
	Hashed bool        `json:"hashed"` // Determines whether the value has been hashed into a keccak256 string.
}

type JSONEvent struct {
	TriggeredAt time.Time               `json:"triggeredAt"`
	Event       *EventInformation       `json:"event"`
	Transaction *TransactionInformation `json:"transaction"`
}
