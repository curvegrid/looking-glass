// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/curvegrid/multibaas/server/app/sqltypes"
	"github.com/ethereum/go-ethereum/common"
)

type Deposit struct {
	OriginChainID      int
	DestinationChainID int
	DepositNonce       int64
	ResourceID         string
	Amount             mbAPI.Number
	Recipient          sqltypes.Address
	TokenAddress       sqltypes.Address
}

// getHandlerAddress gets the handler address from the resourceID recevied from Deposit event
func getHandlerAddress(resourceID string, bc *blockchain.Blockchain) *sqltypes.Address {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/_resourceIDToHandlerAddress",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{json.RawMessage(`"` + resourceID + `"`)},
		TransactionArgs: mbAPI.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	var data struct {
		Ouput sqltypes.Address `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		panic(err)
	}
	return &data.Ouput
}

func getDepositFee(bc *blockchain.Blockchain) *mbAPI.Number {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/_fee",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		TransactionArgs: mbAPI.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	var data struct {
		Ouput mbAPI.Number `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		panic(err)
	}
	return &data.Ouput
}

func getDepositData(d *Deposit) []byte {
	var data []byte
	data = append(data, common.LeftPadBytes(d.Amount.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes([]byte{20}, 32)...)
	data = append(data, d.Recipient.Bytes()...)
	return data
}

// getDeposit creates a Deposit struct from a Deposit event emitted by the Bridge contract
func getDeposit(e *blockchain.JSONEvent, bc *blockchain.Blockchain) *Deposit {
	// event received from the Bridge contract only stores
	// resource ID of the token handler contract,
	// destination chain ID and the deposit nonce.
	chainID := fmt.Sprintf("%v", e.Event.Inputs[0].Value)
	resourceID := fmt.Sprintf("%v", e.Event.Inputs[1].Value)
	depositNonce := fmt.Sprintf("%v", e.Event.Inputs[2].Value)
	// we need to use resourceID to find the token handler contract's address
	handlerAddress := getHandlerAddress(resourceID, bc)

	// get the deposit data by calling depositRecords method of the handler contract
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/erc20handler/methods/_depositRecords",
		bc.MbEndpoint, handlerAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{json.RawMessage(chainID), json.RawMessage(depositNonce)},
		TransactionArgs: mbAPI.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
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

	// parse DepositData from known variables
	var d Deposit
	if d.DepositNonce, err = strconv.ParseInt(depositNonce, 10, 64); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data.Ouput[0], &d.TokenAddress); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data.Ouput[1], &d.DestinationChainID); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data.Ouput[2], &d.ResourceID); err != nil {
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
