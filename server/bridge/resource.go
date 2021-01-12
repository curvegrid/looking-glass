package bridge

import (
	"encoding/json"
	"io/ioutil"

	"github.com/curvegrid/looking-glass/server/blockchain"
	logger "github.com/sirupsen/logrus"
)

type Resource struct {
	ChainID            int                `json:"chainID"`
	TokenAddress       blockchain.Address `json:"tokenAddress"`
	ERC20HandlerAddres blockchain.Address `json:"erc20HandlerAddres"`
}

type ResourceMapping struct {
	ResourceIDToResource map[string][]Resource `json:"resourceIDToResource"`
}

var resourceMapping ResourceMapping

// InitResourcesFromConfigFile initializes resourceMapping from a config file
func InitResourcesFromConfigFile(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Fatalf("cannot read resources config file: %s", err.Error())
	}
	if err := json.Unmarshal(data, &resourceMapping); err != nil {
		logger.Fatalf("cannot unmarshal resources config file: %s", err.Error())
	}
}

func GetResourceMapping() *ResourceMapping {
	return &resourceMapping
}
