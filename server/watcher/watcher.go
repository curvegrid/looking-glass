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
	Blockchain *blockchain.Blockchain
}

func (w *Watcher) getEventStreamURL() *url.URL {
	params := url.Values{}
	params.Add("token", w.Blockchain.BearerToken)
	u := url.URL{
		Scheme: "ws",
		Host:   w.Blockchain.MbEndpoint,
		Path: fmt.Sprintf("api/v0/chains/ethereum/addresses/%s/events/stream",
			w.Blockchain.BridgeAddress.String()),
	}
	u.RawQuery = params.Encode()
	return &u
}

func (w *Watcher) Watch() chan struct{} {
	u := w.getEventStreamURL()
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
			d := w.getDepositData(&e)
			if d != nil {
				logger.Printf("%+v", *d)
			}
		}
	}()

	return done
}
