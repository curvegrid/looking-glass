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
	Address      *Address `json:"address"`
	AddressLabel string   `json:"addressLabel"`
	Name         string   `json:"name"`
	Label        string   `json:"label"`
}

type TransactionInformation struct {
	TXFrom         *Address             `json:"from"`
	TXData         string               `json:"txData"`
	TXHash         common.Hash          `json:"txHash"`
	TXIndexInBlock int                  `json:"txIndexInBlock"`
	BlockHash      common.Hash          `json:"blockHash"`
	BlockNumber    int64                `json:"blockNumber"`
	Contract       *ContractInformation `json:"contract"`
	Method         *MethodInformation   `json:"method"`
}

// EventField holds a field in the event's data.
type EventField struct {
	Name   string      `json:"name"`
	Value  interface{} `json:"value"`
	Hashed bool        `json:"hashed"` // Determines whether the value has been hashed into a keccak256 string.
}

// JSONEvent is an event returned to an API call
type JSONEvent struct {
	TriggeredAt time.Time               `json:"triggeredAt"`
	Event       *EventInformation       `json:"event"`
	Transaction *TransactionInformation `json:"transaction"`
}

// Enum values for event_log_type
const (
	EventLogTypeCreate = "create"
	EventLogTypeRemove = "remove"
)

// ReturnedEvent is an event struct sent through websocket.
type ReturnedEvent struct {
	LogID int    `json:"log_id,omitempty"`
	Type  string `json:"type"` // "create" or "remove"
	*JSONEvent
}

// IsEqual checks if two given events are the same event
func (e *ReturnedEvent) IsEqual(other *ReturnedEvent) bool {
	return e.Transaction.BlockNumber == other.Transaction.BlockNumber &&
		e.Transaction.BlockHash == other.Transaction.BlockHash &&
		e.Transaction.TXIndexInBlock == other.Transaction.TXIndexInBlock &&
		e.Transaction.TXHash == other.Transaction.TXHash &&
		e.Event.EventIndexInLog == other.Event.EventIndexInLog
}
