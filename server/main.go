// Copyright (c) 2021 Curvegrid Inc.

package main

import (
	"fmt"
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
	u := url.URL{
		Scheme: "ws",
		Host:   bc.MbEndpoint,
		Path:   fmt.Sprintf("api/v0/chains/ethereum/addresses/%s/events/stream", bc.BridgeAddress.String()),
	}
	u.RawQuery = params.Encode()
	return &u
}

func main() {
	// parse config
	cfg := Config{}
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile("looking-glass") // gofig will try to load looking-glass.json, looking-glass.toml and looking-glass.yaml
	gofig.Parse(&cfg)

	doneA := watcher.Watch(getEventStreamURL(&cfg.A))
	doneB := watcher.Watch(getEventStreamURL(&cfg.B))
doneLoop:
	for {
		select {
		case <-doneA:
			break doneLoop
		case <-doneB:
			break doneLoop
		}
	}
}
