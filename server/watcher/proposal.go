package watcher

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/mbAPI"
	"github.com/curvegrid/multibaas/server/app/sqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func getProposalDataHash(d *Deposit, handlerAddress *sqltypes.Address) common.Hash {
	log.Printf("data: %s", hex.EncodeToString(append(handlerAddress.Bytes(), getDepositData(d)...)))
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
			json.RawMessage(`"` + d.ResourceID.String() + `"`),
		},
		TransactionArgs: mbAPI.TransactionArgs{
			From:          &bc.HSMAddress,
			SignAndSubmit: true,
		},
	}

	b, _ := json.Marshal(payload)
	log.Printf("payload: %s", string(b))

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
	handlerAddress := getHandlerAddress(d.ResourceID.String(), bc)
	dataHash := getProposalDataHash(d, handlerAddress)

	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/voteProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{
			json.RawMessage(`"` + fmt.Sprintf("%d", d.OriginChainID) + `"`),
			json.RawMessage(`"` + fmt.Sprintf("%d", d.DepositNonce) + `"`),
			json.RawMessage(`"` + d.ResourceID.String() + `"`),
			json.RawMessage(`"` + dataHash.String() + `"`),
		},
		TransactionArgs: mbAPI.TransactionArgs{
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
