// Copyright (c) 2021 Curvegrid Inc.

package watcher

import (
	"net/url"

	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
)

func Watch(u *url.URL) {
	logger.Infof("Connect to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatalf("Cannot connect to websocket dial:", err.Error())
		return
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Fatalf("Cannot read websocket message:", err.Error())
				return
			}
			logger.Println(string(message))
		}
	}()

	<-done
}
