// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import (
	"fmt"

	"github.com/curvegrid/gofig"
)

// Blockchain represents a blockchain deployed in a Multibaas instance
type Blockchain struct {
	MbEndpoint    string  `desc:"MultiBaas endpoint URL"`
	BearerToken   string  `desc:"Mulibaas API key"`
	Confirmations uint    `desc:"number of block confirmations to wait"`
	BridgeAddress Address `desc:"bridge contract address"`
	HSMAddress    Address `desc:"HSM address to automate the signing process"`
}

// BlockchainMapping maps an unique ID to a blockchain
type BlockchainMapping struct {
	ChainIdToBlockchain map[int]*Blockchain
}

var blockchainMapping BlockchainMapping

// GetBlockChainFromID returns a blockchain based on its chain ID
func GetBlockChainFromID(chainID int) (*Blockchain, error) {
	blockchain, exists := blockchainMapping.ChainIdToBlockchain[chainID]
	if !exists {
		return nil, fmt.Errorf("unknown chainID: %d", chainID)
	}
	return blockchain, nil
}

// InitBlockchainsFromConfigFile initializes blockchainMapping from a config file
func InitBlockchainsFromConfigFile(filepath string) {
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile(filepath) // gofig will try to load looking-glass.json, looking-glass.toml and looking-glass.yaml
	gofig.Parse(&blockchainMapping)
}
