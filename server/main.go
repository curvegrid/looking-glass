// Copyright (c) 2021 Curvegrid Inc.

package main

import (
	"github.com/curvegrid/gofig"
	"github.com/curvegrid/looking-glass/server/api"
	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/curvegrid/looking-glass/server/watcher"
)

// Config is the top-level configuration structure
type Config struct {
	Bind              string                       `desc:"HTTP binding address"`
	Web               string                       `desc:"enable web server with the provided root directory"`
	BlockchainMapping blockchain.BlockchainMapping `desc:"a mapping from chain id to the corresponding blockchain"`
	ResourceMapping   bridge.ResourceMapping       `desc:"a mapping from resource id to the corresponding resources"`
	CorsDomains       []string                     `desc:"a list of allowed cors domains"`
}

func main() {
	// Config
	cfg := &Config{}
	cfg.Bind = "localhost:8082"
	gofig.SetEnvPrefix("LG")
	gofig.AddConfigFile("looking-glass")
	gofig.SetConfigFileFlag("c", "configuration file")
	gofig.Parse(cfg)

	// Init mappings
	blockchain.InitBlockchainMapping(&cfg.BlockchainMapping)
	bridge.InitResourceMapping(&cfg.ResourceMapping)
	api.InitCorsDomains(cfg.CorsDomains)

	// API
	go api.InitAPI(cfg.Bind, cfg.Web)

	watcherA := &watcher.Watcher{ChainID: 0}
	watcherB := &watcher.Watcher{ChainID: 1}

	doneA := watcherA.Watch()
	doneB := watcherB.Watch()
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
