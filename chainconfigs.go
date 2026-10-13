// Package chainconfigs exposes the deployed WhitePredict contract addresses of
// every supported chain. The JSON files in this repository are the source of
// truth; they are embedded at build time, so a consumer pins them through its
// module version.
package chainconfigs

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed */config.json
var files embed.FS

// Config is one chain's deployment. An address is nil while the contract is
// not deployed on that chain.
type Config struct {
	ChainID                        uint64  `json:"chainId"`
	Name                           string  `json:"name"`
	PredictionMarket               *string `json:"predictionMarket"`
	PredictionMarketImplementation *string `json:"predictionMarketImplementation"`
	DIAResolver                    *string `json:"diaResolver"`
	USDW                           *string `json:"usdw"`
	Faucet                         *string `json:"faucet"`
}

// All returns every chain config ordered by chain id.
func All() ([]Config, error) {
	paths, err := fs.Glob(files, "*/config.json")
	if err != nil {
		return nil, err
	}
	configs := make([]Config, 0, len(paths))
	for _, path := range paths {
		raw, err := files.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var config Config
		if err := json.Unmarshal(raw, &config); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if want := fmt.Sprintf("%d/config.json", config.ChainID); path != want {
			return nil, fmt.Errorf("%s: chainId %d does not match its directory", path, config.ChainID)
		}
		configs = append(configs, config)
	}
	sort.Slice(configs, func(i, j int) bool { return configs[i].ChainID < configs[j].ChainID })
	return configs, nil
}

// Get returns the config of one chain.
func Get(chainID uint64) (Config, error) {
	configs, err := All()
	if err != nil {
		return Config{}, err
	}
	for _, config := range configs {
		if config.ChainID == chainID {
			return config, nil
		}
	}
	return Config{}, fmt.Errorf("chain %d is not configured", chainID)
}
