// Copyright (c) 2021 Curvegrid Inc.

package main

import (
	"github.com/curvegrid/gofig"
	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/watcher"
)

type Config struct {
	A blockchain.Blockchain
	B blockchain.Blockchain
}

func main() {
	// parse config
	cfg := Config{}
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile("looking-glass") // gofig will try to load looking-glass.json, looking-glass.toml and looking-glass.yaml
	gofig.Parse(&cfg)

	watcherA := &watcher.Watcher{Blockchain: &cfg.A}
	watcherB := &watcher.Watcher{Blockchain: &cfg.B}

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
