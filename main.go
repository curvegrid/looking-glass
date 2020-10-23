// Copyright (c) 2020 Curvegrid Inc.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/curvegrid/gofig"
	"github.com/ethereum/go-ethereum/common"
)

type Blockchain struct {
	Endpoint      string         `desc:"MultiBaas endpoint URL"`
	Token         string         `desc:"MultiBaas token"`
	Confirmations uint           `desc:"number of block confirmations to wait"`
	Address       common.Address `desc:"looking glass address"`
	Assets        common.Address `desc:"asset addresses to monitor for new transactions"`
}

type Config struct {
	A Blockchain
	B Blockchain
}

type Input struct {
	Name  string
	Value string
}

type ContractStruct struct {
	Address common.Address
	Label   string
}

type EventStruct struct {
	Name      string
	Signature string
	Inputs    []Input
	Contract  ContractStruct
}

type EventTransaction struct {
	TriggeredAt string
	Event       EventStruct
}

func GetEvents(endpoint string, token string) []EventTransaction {
	endpoint += "/api/v0/events"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Printf("body: %s\n", string(body))

	// unmarshal body
	et := []EventTransaction{}
	err = json.Unmarshal(body, &et)
	if err != nil {
		panic(err)
	}

	return et
}

func main() {
	// parse config
	cfg := Config{}
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile("looking-glass") // gofig will try to load default.json, default.toml and default.yaml
	gofig.Parse(&cfg)

	blockchains := []Blockchain{cfg.A, cfg.B}

	wg := sync.WaitGroup{}

	for _, blockchain := range blockchains {
		wg.Add(1)
		go func(b Blockchain) {
			// store last event

			// poll for events
			GetEvents(b.Endpoint, b.Token)

			// pause
			time.Sleep(time.Second * 5)

		}(blockchain)
	}

	wg.Wait()

	// poll for transactions on A and B
	// todo: listen for transactions via websockets

	// when a transaction is received, wait for a number of blocks before processing

	// burn on one side

	// mint on the other
}
