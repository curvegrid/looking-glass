// Copyright (c) 2021 Curvegrid Inc.

package bridge

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/customError"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/ethereum/go-ethereum/common"
)

// Deposit is used to represent a cross-chain Deposit transaction initiated
// by the Bridge contract.
type Deposit struct {
	OriginChainID      int
	DestinationChainID int
	DepositNonce       int64
	ResourceID         string
	Amount             blockchain.Number
	Recipient          blockchain.Address
	TokenAddress       blockchain.Address
}

// getHandlerAddress gets the handler address from the resourceID recevied from Deposit event
func getHandlerAddress(resourceID string, bc *blockchain.Blockchain) (*blockchain.Address, error) {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/_resourceIDToHandlerAddress",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []interface{}{resourceID},
		TransactionArgs: blockchain.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return nil, err
	}
	if result.Status != 200 {
		return nil, customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	var data struct {
		Ouput blockchain.Address `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		return nil, err
	}
	return &data.Ouput, nil
}

// getDepositFee returns the deposit fee required by the Bridge contract
// to initiate a deposit transaction
func getDepositFee(bc *blockchain.Blockchain) (*blockchain.Number, error) {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/_fee",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		TransactionArgs: blockchain.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return nil, err
	}
	if result.Status != 200 {
		return nil, customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	var data struct {
		Ouput blockchain.Number `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		return nil, err
	}
	return &data.Ouput, nil
}

// getDepositData returns the packed deposit data with the format:
// amount(bytes32)recipient_address_length(bytes32)recipient_address(bytes)
func getDepositData(d *Deposit) []byte {
	var data []byte
	data = append(data, common.LeftPadBytes(d.Amount.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes([]byte{20}, 32)...)
	data = append(data, d.Recipient.Bytes()...)
	return data
}

// GetDeposit parses a cross-chain deposit transaction from a Deposit event emitted from the Bridge contract
func GetDeposit(e *blockchain.JSONEvent, bc *blockchain.Blockchain) (*Deposit, error) {
	// event received from the Bridge contract only stores
	// resource ID of the token handler contract,
	// destination chain ID and the deposit nonce.
	chainID := fmt.Sprint(e.Event.Inputs[0].Value)
	resourceID := fmt.Sprint(e.Event.Inputs[1].Value)
	depositNonce := fmt.Sprint(e.Event.Inputs[2].Value)
	// we need to use resourceID to find the token handler contract's address
	handlerAddress, err := getHandlerAddress(resourceID, bc)
	if err != nil {

		return nil, err
	}

	// get the deposit data by calling depositRecords method of the handler contract
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/erc20handler/methods/_depositRecords",
		bc.MbEndpoint, handlerAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []interface{}{chainID, depositNonce},
		TransactionArgs: blockchain.TransactionArgs{
			From: &bc.HSMAddress,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return nil, err
	}
	if result.Status != 200 {
		return nil, customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	var data struct {
		Ouput []json.RawMessage `json:"output"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		return nil, err
	}

	// parse DepositData from known variables
	var d Deposit
	if d.DepositNonce, err = strconv.ParseInt(depositNonce, 10, 64); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data.Ouput[0], &d.TokenAddress); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data.Ouput[1], &d.DestinationChainID); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data.Ouput[2], &d.ResourceID); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data.Ouput[3], &d.Recipient); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data.Ouput[5], &d.Amount); err != nil {
		return nil, err
	}
	return &d, nil
}

// CreateDeposit initiates a cross-chain transfer by calling Bridge contract's Deposit method
func CreateDeposit(d *Deposit, bc *blockchain.Blockchain, submit bool) (*mbAPI.APICallResult, error) {
	fee, err := getDepositFee(bc)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/deposit",
		bc.MbEndpoint, bc.BridgeAddress.String())

	// send from the recipient address in the original chain to the same address
	// in the destination chain or use HSM to automate the signing process
	from := &d.Recipient
	if submit {
		from = &bc.HSMAddress
	}
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []interface{}{
			d.DestinationChainID,
			d.ResourceID,
			"0x" + hex.EncodeToString(getDepositData(d)),
		},
		TransactionArgs: blockchain.TransactionArgs{
			From:          from,
			Value:         fee,
			SignAndSubmit: submit,
		},
	}
	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return nil, err
	}
	if result.Status != 200 {
		return nil, customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	return result, nil
}
