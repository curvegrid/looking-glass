package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

func InitAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/api/deposit", Deposit).Methods("POST")

	logger.Info("Starting Looking-Glass server on port 8082")
	http.ListenAndServe(":8082", r)
}

// parseJSONBody reads the body from an http request and unmarshals
// into the provided object, returning the appropriate response in case of error
func parseJSONBody(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	var deposit DepositReqBodyJSON
	if err := parseJSONBody(r.Body, &deposit); err != nil {
		logger.Error(err.Error())
		return
	}
	originIDs := bridge.GetResourceIDsFromTokenAddress(deposit.OriginTokenAddress, deposit.OriginChainID)
	destinationIDs := bridge.GetResourceIDsFromTokenAddress(deposit.DestinationTokenAddress, deposit.DestinationChainID)

	idsInOrigin := make(map[string]bool)
	for _, id := range originIDs {
		idsInOrigin[id] = true
	}
	var selectedID string
	for _, id := range destinationIDs {
		if idsInOrigin[id] {
			selectedID = id
			break
		}
	}
	// cannot find any resourceID associated to
	// both two token addresses specified in the
	// cross-chain deposit
	if len(selectedID) == 0 {
		return
	}

	d := &bridge.Deposit{
		ResourceID:         selectedID,
		Amount:             deposit.Amount,
		Recipient:          deposit.Recipient,
		OriginChainID:      deposit.OriginChainID,
		DestinationChainID: deposit.DestinationChainID,
	}
	bc, err := blockchain.GetBlockChainFromID(deposit.OriginChainID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if err := bridge.CreateDeposit(d, bc); err != nil {
		logger.Error(err.Error())
	}
}
