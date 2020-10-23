// Copyright (c) 2020 Curvegrid Inc.

package main

import (
	"github.com/curvegrid/gofig"
	"github.com/ethereum/go-ethereum/common"
)

type Blockchain struct {
	Endpoint string `desc:"MultiBaas endpoint URL"`
	Token string `desc:"MultiBaas token"`
	Confirmations uint `desc:"number of block confirmations to wait"`
	Address common.Address `desc:"looking glass address"`
	Assets common.Address `desc:"asset addresses to monitor for new transactions"`
}

type Config struct {
	A Blockchain
	B Blockchain
}

func main() {
	// parse config

	// listen for transactions on A and B
	
}
