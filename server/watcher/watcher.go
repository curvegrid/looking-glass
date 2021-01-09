// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"fmt"
	"net/url"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
)

type Watcher struct {
	ChainID int
}

func (w *Watcher) getEventStreamURL(bc *blockchain.Blockchain) *url.URL {
	params := url.Values{}
	params.Add("token", bc.BearerToken)
	u := url.URL{
		Scheme: "ws",
		Host:   bc.MbEndpoint,
		Path: fmt.Sprintf("api/v0/chains/ethereum/addresses/%s/events/stream",
			bc.BridgeAddress.String()),
	}
	u.RawQuery = params.Encode()
	return &u
}

func (w *Watcher) handleDepositEvent(e *blockchain.JSONEvent, bc *blockchain.Blockchain) {
	d := getDeposit(e, bc)
	d.OriginChainID = w.ChainID
	logger.Infof("Got a Deposit event from chain %d: %+v", w.ChainID, *d)

	err := voteProposal(d)
	logger.Infof("HSM: vote yes to a transfer proposal originated from the deposit %+v", *d)
	if err != nil {
		logger.Error(err)
	}

	// for simplified version with only one relayer for each bridge contract,
	// we execute the proposal right after voting
	executeProposal(d)
	logger.Infof("HSM: execute the a transfer proposal originated from the deposit %+v", *d)
}

func (w *Watcher) Watch() chan struct{} {
	bc, err := blockchain.GetBlockChainFromID(w.ChainID)
	if err != nil {
		panic(err)
	}

	u := w.getEventStreamURL(bc)
	logger.Infof("Connect to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatalf("Cannot connect to websocket dial:", err.Error())
		return nil
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var e blockchain.JSONEvent
			c.ReadJSON(&e)
			if err != nil {
				logger.Fatalf("Cannot read websocket message:", err.Error())
				return
			}
			switch e.Event.Name {
			case "Deposit":
				w.handleDepositEvent(&e, bc)
			}
		}
	}()

	return done
}
