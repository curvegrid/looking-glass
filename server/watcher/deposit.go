// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"encoding/json"
	"fmt"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/ethereum/go-ethereum/common"
)

type DepositData struct {
	Amount       mbAPI.Number
	Recipient    common.Address
	TokenAddress common.Address
}

// getHandlerAddress gets the handler address from the resourceID recevied from Deposit event
func (w *Watcher) getHandlerAddress(resourceID string) *common.Address {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/_resourceIDToHandlerAddress",
		w.Blockchain.MbEndpoint, w.Blockchain.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{json.RawMessage(`"` + resourceID + `"`)},
		TransactionArgs: mbAPI.TransactionArgs{
			From: &w.Blockchain.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, w.Blockchain.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	var data struct {
		Ouput common.Address `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		panic(err)
	}
	return &data.Ouput
}

// getDepositData gets DepositData from a Deposit event emitted by the Bridge contract
func (w *Watcher) getDepositData(e *blockchain.JSONEvent) *DepositData {
	if e.Event.Name != "Deposit" {
		return nil
	}
	// event received from the Bridge contract only stores
	// resource ID of the token handler contract,
	// destination chain ID and the deposit nonce.
	chainID := fmt.Sprintf("%v", e.Event.Inputs[0].Value)
	resourceID := fmt.Sprintf("%v", e.Event.Inputs[1].Value)
	depositNonce := fmt.Sprintf("%v", e.Event.Inputs[2].Value)
	// we need to use resourceID to find the token handler contract's address
	handlerAddress := w.getHandlerAddress(resourceID)

	// get the deposit data by calling depositRecords method of the handler contract
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/erc20handler/methods/_depositRecords",
		w.Blockchain.MbEndpoint, handlerAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{json.RawMessage(chainID), json.RawMessage(depositNonce)},
		TransactionArgs: mbAPI.TransactionArgs{
			From: &w.Blockchain.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, w.Blockchain.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	var data struct {
		Ouput []json.RawMessage `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		panic(err)
	}

	var d DepositData
	if err := json.Unmarshal(data.Ouput[0], &d.TokenAddress); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data.Ouput[3], &d.Recipient); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data.Ouput[5], &d.Amount); err != nil {
		panic(err)
	}
	return &d
}
