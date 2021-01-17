package bridge

import (
	"github.com/curvegrid/looking-glass/server/blockchain"
)

type Resource struct {
	ChainID            int                `json:"chainID"`
	TokenAddress       blockchain.Address `json:"tokenAddress"`
	ERC20HandlerAddres blockchain.Address `json:"erc20HandlerAddres"`
}

// ResourceMapping maps from a resource id to a list of corresponding
// resources across different chains.
type ResourceMapping struct {
	ResourceIDToResource map[string][]Resource `json:"resourceIDToResource"`
}

var resourceMapping *ResourceMapping

// InitResourceMapping initializes resource mapping
func InitResourceMapping(m *ResourceMapping) { resourceMapping = m }

// GetResourceMapping returns the resource mapping
func GetResourceMapping() *ResourceMapping { return resourceMapping }

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
