package watcher

import (
	"math/big"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/ethereum/go-ethereum/common"
	logger "github.com/sirupsen/logrus"
)

type DepositData struct {
	Amount       *big.Int
	Recipient    *common.Address
	TokenAddress *common.Address
}

func getDepositData(e *blockchain.JSONEvent) *DepositData {
	if e.Event.Name != "Deposit" {
		return nil
	}
	chainID := e.Event.Inputs[0].Value
	resourceID := e.Event.Inputs[1].Value
	depositNonce := e.Event.Inputs[2].Value
	logger.Printf("event inputs :%v %v %v", chainID, resourceID, depositNonce)
	return &DepositData{}
}
