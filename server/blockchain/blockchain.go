// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import (
	"fmt"

	"github.com/curvegrid/gofig"
	"github.com/curvegrid/multibaas/server/app/sqltypes"
)

type Blockchain struct {
	MbEndpoint    string           `desc:"MultiBaas endpoint URL"`
	BearerToken   string           `desc:"Mulibaas API key"`
	Confirmations uint             `desc:"number of block confirmations to wait"`
	BridgeAddress sqltypes.Address `desc:"bridge contract address"`
	HSMAddress    sqltypes.Address `desc:"HSM address to automate the signing process"`
}

type BlockchainMapping struct {
	ChainIdToBlockchain map[int]*Blockchain
}

var blockchainMapping BlockchainMapping

func GetBlockChainFromID(chainID int) (*Blockchain, error) {
	blockchain, exists := blockchainMapping.ChainIdToBlockchain[chainID]
	if !exists {
		return nil, fmt.Errorf("unknown chainID: %d", chainID)
	}
	return blockchain, nil
}

func InitBlockchainsFromConfigFile(filepath string) {
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile(filepath) // gofig will try to load looking-glass.json, looking-glass.toml and looking-glass.yaml
	gofig.Parse(&blockchainMapping)
}
