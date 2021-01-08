// Copyright (c) 2021 Curvegrid Inc.

package blockchain

import "github.com/ethereum/go-ethereum/common"

type Blockchain struct {
	MbEndpoint    string         `desc:"MultiBaas endpoint URL"`
	BearerToken   string         `desc:"Mulibaas DApp API key"`
	Confirmations uint           `desc:"number of block confirmations to wait"`
	BridgeAddress common.Address `desc:"looking glass address"`
}
