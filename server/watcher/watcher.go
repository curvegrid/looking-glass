// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"fmt"
	"net/url"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
)

// Watcher watches events emitted from a blockchain and handles those events
type Watcher struct {
	ChainID int
}

// getEventStreamURL gets the API for event streaming of a MB instance
// associated with the watched blockchain.
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

// handleDepositEvent reads a Deposit event. It then uses HSM auto-signing to
// vote/execute the (cross-chain) deposit proposal associated with the read event.
func (w *Watcher) handleDepositEvent(e *blockchain.JSONEvent, bc *blockchain.Blockchain) error {
	d, err := bridge.GetDeposit(e, bc)
	if err != nil {
		return err
	}
	d.OriginChainID = w.ChainID
	logger.Infof("Got a Deposit event from chain %d: %+v", w.ChainID, *d)

	err = bridge.VoteProposal(d)
	if err != nil {
		return err
	}
	logger.Infof("HSM: voted yes to a transfer proposal originated from the deposit %+v", *d)

	// for simplified version with only one relayer for each bridge contract,
	// we execute the proposal right after voting
	bridge.ExecuteProposal(d)
	logger.Infof("HSM: executed a transfer proposal originated from the deposit %+v", *d)
	return nil
}

// Watch starts a watcher. A watcher watches events from the blockchain
// using event streaming API and handles them using HSM.
func (w *Watcher) Watch() chan struct{} {
	bc, err := blockchain.GetBlockChainFromID(w.ChainID)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	u := w.getEventStreamURL(bc)
	logger.Infof("Connect to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatalf("Cannot connect to websocket dial:", err.Error())
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var e blockchain.JSONEvent
			c.ReadJSON(&e)
			if err != nil {
				logger.Errorf("Cannot read websocket message:", err.Error())
				continue
			}
			switch e.Event.Name {
			case "Deposit":
				err := w.handleDepositEvent(&e, bc)
				if err != nil {
					logger.Errorf("Cannot handle event %s: %s", e.Event.Name, err.Error())
				}
			}
		}
	}()

	return done
}
