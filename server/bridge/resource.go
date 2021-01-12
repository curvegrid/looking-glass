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

// GetResourceIDsFromTokenAddress returns a list of resourceID
// associated with a given tokenAddress, chainID pair.
func GetResourceIDsFromTokenAddress(tokenAddress blockchain.Address, chainID int) []string {
	var ids []string
	for resourceID, resources := range resourceMapping.ResourceIDToResource {
		for _, resource := range resources {
			if resource.TokenAddress.String() == tokenAddress.String() && chainID == resource.ChainID {
				ids = append(ids, resourceID)
				break
			}
		}
	}
	return ids
}
