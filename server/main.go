// Copyright (c) 2021 Curvegrid Inc.

package main

import (
	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/curvegrid/looking-glass/server/watcher"
)

func main() {
	blockchain.InitBlockchainsFromConfigFile("looking-glass")
	bridge.InitResourcesFromConfigFile("resources.json")

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
