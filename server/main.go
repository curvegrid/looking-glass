// Copyright (c) 2021 Curvegrid Inc.

package main

import (
	"log"
	"net/url"

	"github.com/curvegrid/gofig"
	"github.com/curvegrid/looking-glass/server/watcher"
	"github.com/ethereum/go-ethereum/common"
)

type Blockchain struct {
	MbEndpoint    string         `desc:"MultiBaas endpoint URL"`
	BearerToken   string         `desc:"Mulibaas DApp API key"`
	Confirmations uint           `desc:"number of block confirmations to wait"`
	BridgeAddress common.Address `desc:"looking glass address"`
	TokenAddress  common.Address `desc:"asset addresses to monitor for new transactions"`
}

type Config struct {
	A Blockchain
	B Blockchain
}

func getEventStreamURL(bc *Blockchain) *url.URL {
	params := url.Values{}
	params.Add("token", bc.BearerToken)
	u := url.URL{Scheme: "ws", Host: bc.MbEndpoint, Path: "api/v0/chains/ethereum/addresses/autotoken/events/stream"}
	u.RawQuery = params.Encode()
	return &u
}

func main() {
	// parse config
	cfg := Config{}
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile("looking-glass") // gofig will try to load default.json, default.toml and default.yaml
	gofig.Parse(&cfg)

	log.Printf("%+v", cfg)

	watcher.Watch(getEventStreamURL(&cfg.A))
}
