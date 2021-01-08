// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"net/url"

	"github.com/curvegrid/looking-glass/server/event"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
)

func Watch(u *url.URL) chan struct{} {
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
			var e event.JSONEvent
			c.ReadJSON(&e)
			if err != nil {
				logger.Fatalf("Cannot read websocket message:", err.Error())
				return
			}
			logger.Printf("%+v", e)
		}
	}()

	return done
}
