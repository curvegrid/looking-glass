package config

import (
	"github.com/curvegrid/gofig"
	"github.com/curvegrid/looking-glass/server/blockchain"
)

type TwoChainsConfig struct {
	A blockchain.Blockchain
	B blockchain.Blockchain
}

func InitTwoChainsConfig() TwoChainsConfig {
	cfg := TwoChainsConfig{}
	gofig.SetEnvPrefix("LG")
	gofig.SetConfigFileFlag("c", "config file")
	gofig.AddConfigFile("config/looking-glass") // gofig will try to load looking-glass.json, looking-glass.toml and looking-glass.yaml
	gofig.Parse(&cfg)
	return cfg
}
