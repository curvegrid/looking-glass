// Copyright (c) 2021 Curvegrid Inc.

package bridge

import (
	"encoding/hex"
	"fmt"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/customError"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/ethereum/go-ethereum/crypto"
)

// Proposal is used to represent a proposal of a cross-chain deposit.
// Proposal needs to be passed by relayers from the destination chain
// to execute the transaction.
type Proposal struct {
	OriginChainID      int
	DestinationChainID int
	Status             int
	DepositNonce       int64
	ResourceID         string
	DataHash           string
}

// getProposalDataHash returns the data hash used to call voteProposal method
// of a Bridge contract.
func getProposalDataHash(d *Deposit, handlerAddress *blockchain.Address) string {
	return crypto.Keccak256Hash(append(handlerAddress.Bytes(), getDepositData(d)...)).String()
}

// ExecuteProposal executes a Deposit proposal by calling Bridge contract's ExecuteProposal method
func ExecuteProposal(d *Deposit) error {
	bc, err := blockchain.GetBlockChainFromID(d.DestinationChainID)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/executeProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []interface{}{
			d.OriginChainID,
			d.DepositNonce,
			"0x" + hex.EncodeToString(getDepositData(d)),
			d.ResourceID,
		},
		TransactionArgs: blockchain.TransactionArgs{
			From:          &bc.HSMAddress,
			SignAndSubmit: true,
		},
	}

	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return err
	}
	if result.Status != 200 {
		return customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	return nil
}

// ExecuteProposal votes to a Deposit proposal by calling Bridge contract's VoteProposal method
func VoteProposal(d *Deposit) error {
	bc, err := blockchain.GetBlockChainFromID(d.DestinationChainID)
	if err != nil {
		return err
	}
	handlerAddress, err := getHandlerAddress(d.ResourceID, bc)
	if err != nil {
		return err
	}
	dataHash := getProposalDataHash(d, handlerAddress)

	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/voteProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []interface{}{
			d.OriginChainID,
			d.DepositNonce,
			d.ResourceID,
			dataHash,
		},
		TransactionArgs: blockchain.TransactionArgs{
			From:          &bc.HSMAddress,
			SignAndSubmit: true,
		},
	}

	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		return err
	}
	if result.Status != 200 {
		return customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	return nil
}
