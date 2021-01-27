// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/curvegrid/looking-glass/server/customError"
	"github.com/curvegrid/looking-glass/server/mbAPI"
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
func (w *Watcher) handleDepositEvent(e *blockchain.ReturnedEvent) error {
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
	logger.Infof("Handle a Deposit event from chain %d: %+v", w.ChainID, *d)

	err = bridge.VoteProposal(d)
	if err != nil {
		return err
	}
	logger.Infof("HSM: voted yes to a transfer proposal originated from the deposit %+v", *d)
	return nil
}

func getLatestBlockNumber(bc *blockchain.Blockchain) (int64, error) {
	endpoint := fmt.Sprintf("http://%s/api/v0/chains/ethereum/status", bc.MbEndpoint)
	result, err := mbAPI.Get(endpoint, bc.BearerToken)
	if err != nil {
		return 0, err
	}
	if result.Status != 200 {
		return 0, customError.NewAPICallError(endpoint, result.Status, result.Message)
	}
	var data struct {
		BlockNumber int64 `json:"blockNumber"`
	}
	if err := json.Unmarshal(result.Result, &data); err != nil {
		return 0, err
	}
	return data.BlockNumber, nil
}

func (w *Watcher) parseProposalFromEvent(e *blockchain.ReturnedEvent) (*bridge.Proposal, error) {
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
func (w *Watcher) handleProposalEvent(e *blockchain.ReturnedEvent) error {
	p, err := w.parseProposalFromEvent(e)
	if err != nil {
		return err
	}

	logger.Infof("Handle ProposalEvent event from chain %d with proposal %+v", w.ChainID, *p)
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

func (w *Watcher) handleProposalVote(e *blockchain.ReturnedEvent) error {
	p, err := w.parseProposalFromEvent(e)
	if err != nil {
		return err
	}
	logger.Infof("Handle ProposalVote event from chain %d with proposal %+v", w.ChainID, *p)
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
		logger.Fatalf("Cannot connect to websocket dial: %s", err.Error())
	}

	done := make(chan struct{})
	go func() {
		/*
		   To implement block confirmations for the watcher,
		   we store a list of events received from the monitor.
		   We only handle an event of a given block if after that
		   block, a specific number of blocks has been confirmed.
		*/

		defer close(done)

		eventCh := make(chan blockchain.ReturnedEvent)

		// read events from websocket and pass them to event Channel
		go func() {
			for {
				var e blockchain.ReturnedEvent
				c.ReadJSON(&e)
				if err != nil {
					logger.Errorf("Cannot read websocket message: %s", err.Error())
					continue
				}

				eventCh <- e
				logger.Printf("Got event %s (type: %s) from chain %d (block %d)",
					e.Event.Name, e.Type, w.ChainID, e.Transaction.BlockNumber)
			}
		}()

		var events []blockchain.ReturnedEvent
		for {
			select {
			case e := <-eventCh:
				{
					if e.Type == blockchain.EventLogTypeCreate {
						events = append(events, e)
					} else { // = EventLogTypeRemove
						// if the event received from the monitor is deleted from the blockchain,
						// we remove it from the list of events we need to consider.
						for i, v := range events {
							if v.IsEqual(&e) {
								events = append(events[:i], events[i+1:]...)
								break
							}
						}
					}
				}
			default:
				{
					latestBlock, err := getLatestBlockNumber(bc)
					if err != nil {
						logger.Errorf("Cannot retrieve the latest block of chain %d: %s", w.ChainID, err.Error())
						continue
					}

					// we only handle the first event if "bc.Confirmations" blocks have been confirmed since the event
					for len(events) > 0 && events[0].Transaction.BlockNumber+int64(bc.Confirmations) <= latestBlock {
						e := events[0]
						events = events[1:]
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

					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return done
}
