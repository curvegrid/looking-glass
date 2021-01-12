package api

import "github.com/curvegrid/looking-glass/server/blockchain"

type DepositReqBodyJSON struct {
	Amount                  blockchain.Number  `json:"amount"`
	Recipient               blockchain.Address `json:"recipient"`
	OriginChainID           int                `json:"originChainID"`
	OriginTokenAddress      blockchain.Address `json:"originTokenAddress"`
	DestinationChainID      int                `json:"destinationChainID"`
	DestinationTokenAddress blockchain.Address `json:"destinationTokenAddress"`
}
