// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"fmt"
	"net/url"
	"strconv"

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
func (w *Watcher) handleDepositEvent(e *blockchain.JSONEvent) error {
	destinationChainID, err := strconv.Atoi(fmt.Sprint(e.Event.Inputs[0].Value))
	if err != nil {
		return err
	}
	resourceID := fmt.Sprint(e.Event.Inputs[1].Value)
	depositNonce, err := strconv.ParseInt(fmt.Sprint(e.Event.Inputs[2].Value), 10, 64)
	if err != nil {
		return err
	}

	// retrieve the blockchain data from the watcher (current chain)
	// to query the deposit data.
	bc, err := blockchain.GetBlockChainFromID(w.ChainID)
	if err != nil {
		return err
	}

	d, err := bridge.GetDeposit(bc, destinationChainID, resourceID, depositNonce)
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
	return nil
}

func (w *Watcher) parseProposalFromEvent(e *blockchain.JSONEvent) (*bridge.Proposal, error) {
	var p bridge.Proposal
	var err error

	p.DestinationChainID = w.ChainID
	p.OriginChainID, err = strconv.Atoi(fmt.Sprint(e.Event.Inputs[0].Value))
	if err != nil {
		return nil, err
	}
	p.ResourceID = fmt.Sprint(e.Event.Inputs[1].Value)
	p.DepositNonce, err = strconv.ParseInt(fmt.Sprint(e.Event.Inputs[2].Value), 10, 64)
	if err != nil {
		return nil, err
	}
	p.Status, err = strconv.Atoi(fmt.Sprint(e.Event.Inputs[3].Value))
	if err != nil {
		return nil, err
	}
	p.DataHash = fmt.Sprint(e.Event.Inputs[4].Value)

	return &p, nil
}

// handleProposalEvent reads a ProposalEvent. Depends on the proposal's status, the
// watcher then uses HSM to execute the proposal.
func (w *Watcher) handleProposalEvent(e *blockchain.JSONEvent) error {
	p, err := w.parseProposalFromEvent(e)
	if err != nil {
		return err
	}

	logger.Infof("Got ProposalEvent event from chain %d with proposal %+v", w.ChainID, *p)
	if p.Status != 2 {
		return nil
	}

	// retrieve blockchain from the origin chain to
	// query the deposit data
	bc, err := blockchain.GetBlockChainFromID(p.OriginChainID)
	if err != nil {
		return err
	}

	d, err := bridge.GetDeposit(bc, p.DestinationChainID, p.ResourceID, p.DepositNonce)
	if err != nil {
		return err
	}
	if err := bridge.ExecuteProposal(d); err != nil {
		return err
	}
	logger.Infof("HSM: executed a transfer proposal originated from the deposit %+v", *d)
	return nil
}

func (w *Watcher) handleProposalVote(e *blockchain.JSONEvent) error {
	p, err := w.parseProposalFromEvent(e)
	if err != nil {
		return err
	}
	logger.Infof("Got ProposalVote event from chain %d with proposal %+v", w.ChainID, *p)
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
				if err := w.handleDepositEvent(&e); err != nil {
					logger.Errorf("Cannot handle event %s: %s", e.Event.Name, err.Error())
				}
			case "ProposalEvent":
				{
					if err := w.handleProposalEvent(&e); err != nil {
						logger.Errorf("Cannot handle event %s: %s", e.Event.Name, err.Error())
					}
				}
			case "ProposalVote":
				{
					if err := w.handleProposalVote(&e); err != nil {
						logger.Errorf("Cannot handle event %s: %s", e.Event.Name, err.Error())
					}
				}
			}
		}
	}()

	return done
}
