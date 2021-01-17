// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import (
	"fmt"
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

var blockchainMapping *BlockchainMapping

// Initialize blockchain mapping
func InitBlockchainMapping(m *BlockchainMapping) { blockchainMapping = m }

// GetBlockChainFromID returns a blockchain based on its chain ID
func GetBlockChainFromID(chainID int) (*Blockchain, error) {
	blockchain, exists := blockchainMapping.ChainIdToBlockchain[chainID]
	if !exists {
		return nil, fmt.Errorf("unknown chainID: %d", chainID)
	}
	return blockchain, nil
}
