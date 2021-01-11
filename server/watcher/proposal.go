package watcher

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func getProposalDataHash(d *Deposit, handlerAddress *blockchain.Address) common.Hash {
	return crypto.Keccak256Hash(append(handlerAddress.Bytes(), getDepositData(d)...))
}

func executeProposal(d *Deposit) error {
	bc, err := blockchain.GetBlockChainFromID(d.DestinationChainID)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/executeProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{
			json.RawMessage(`"` + fmt.Sprintf("%d", d.OriginChainID) + `"`),
			json.RawMessage(`"` + fmt.Sprintf("%d", d.DepositNonce) + `"`),
			json.RawMessage(`"0x` + hex.EncodeToString(getDepositData(d)) + `"`),
			json.RawMessage(`"` + d.ResourceID + `"`),
		},
		TransactionArgs: blockchain.TransactionArgs{
			From:          &bc.HSMAddress,
			SignAndSubmit: true,
		},
	}

	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	return nil
}

func voteProposal(d *Deposit) error {
	bc, err := blockchain.GetBlockChainFromID(d.DestinationChainID)
	if err != nil {
		return err
	}
	handlerAddress := getHandlerAddress(d.ResourceID, bc)
	dataHash := getProposalDataHash(d, handlerAddress)

	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/voteProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{
			json.RawMessage(`"` + fmt.Sprintf("%d", d.OriginChainID) + `"`),
			json.RawMessage(`"` + fmt.Sprintf("%d", d.DepositNonce) + `"`),
			json.RawMessage(`"` + d.ResourceID + `"`),
			json.RawMessage(`"` + dataHash.String() + `"`),
		},
		TransactionArgs: blockchain.TransactionArgs{
			From:          &bc.HSMAddress,
			SignAndSubmit: true,
		},
	}

	result, err := mbAPI.Post(endpoint, bc.BearerToken, payload)
	if err != nil {
		panic(err)
	}
	if result.Status != 200 {
		panic(result.Message)
	}
	return nil
}
