package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/curvegrid/looking-glass/server/blockchain"
	"github.com/curvegrid/looking-glass/server/bridge"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/sirupsen/logrus"
)

// InitAPI initializes looking-glass APIs
func InitAPI(host string, web string) {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/api/deposit", Deposit).Methods("POST")
	muxRouter.HandleFunc("/api/resources", GetResources).Methods("GET")

	echoRouter := echo.New()
	if web != "" {
		// static file serving
		echoRouter.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   web,
			Browse: false,
			HTML5:  true,
			Index:  "index.html",
		}))
	}
	apiRouter := echoRouter.Group("/api",
		CORSMiddleware(), // CORS support
	)
	apiRouter.Any("/*", echo.WrapHandler(muxRouter))

	logger.Infof("Starting Looking-Glass server on %s", host)
	logger.Fatal(echoRouter.Start(host))
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

// writeJSON writes a JSON data into the http response
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseBytes, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("Unable to send HTTP response, marshall JSON: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		logger.Errorf("Unable to send HTTP response, response writer returned: %v", err)
	}
}

// GetResources returns a resource mapping that the application knows
func GetResources(w http.ResponseWriter, r *http.Request) {
	resourceMapping := bridge.GetResourceMapping()
	writeJSON(w, http.StatusOK, resourceMapping)
}

// Deposit receives cross-chain deposit data and initiates the
// transaction using the Bridge contract.
func Deposit(w http.ResponseWriter, r *http.Request) {
	var deposit DepositReqBodyJSON
	if err := parseJSONBody(r.Body, &deposit); err != nil {
		logger.Error(err.Error())
		return
	}

	// first, we need to find the shared resourceID across two different chains
	// for two given tokens
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
	submit := false
	result, err := bridge.CreateDeposit(d, bc, submit)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	writeJSON(w, result.Status, result.Result)
}
