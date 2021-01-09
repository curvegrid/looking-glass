package watcher

import (
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
	return crypto.Keccak256Hash(append(handlerAddress.Bytes(), getDepositData(d)...))
}

func voteProposal(d *Deposit) error {
	bc, err := blockchain.GetBlockChainFromID(d.DestinationChainID)
	if err != nil {

	}
	handlerAddress := getHandlerAddress(d.ResourceID.String(), bc)
	dataHash := getProposalDataHash(d, handlerAddress)

	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/addresses/%s/contracts/bridge/methods/voteProposal",
		bc.MbEndpoint, bc.BridgeAddress.String())
	payload := mbAPI.JSONPOSTMethodArgs{
		Args: []json.RawMessage{
			json.RawMessage(`"` + fmt.Sprintf("%d", d.DestinationChainID) + `"`),
			json.RawMessage(`"` + fmt.Sprintf("%d", d.DepositNonce) + `"`),
			json.RawMessage(`"` + d.ResourceID.String() + `"`),
			json.RawMessage(`"` + dataHash.String() + `"`),
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
