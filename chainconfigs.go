// Package chainconfigs is the WhitePredict chain registry: every chain the
// product knows, its public endpoints, its deployed contracts and whether it
// is enabled. The JSON files in this repository are the source of truth; they
// are embedded at build time, so a consumer pins them through its module
// version. Everything here is public; private RPC credentials and signer keys
// never belong in this repository.
package chainconfigs

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
)

//go:embed */config.json
var files embed.FS

// Config is one chain's deployment. An address is nil while the contract is
// not deployed on that chain. An enabled chain is served and indexed, so it
// must name its prediction market, deployment block and confirmations.
type Config struct {
	ChainID        uint64         `json:"chainId"`
	Name           string         `json:"name"`
	Enabled        bool           `json:"enabled"`
	Testnet        bool           `json:"testnet"`
	NativeCurrency NativeCurrency `json:"nativeCurrency"`
	// RPCURL is a public endpoint; a consumer may replace it with a private one.
	RPCURL      string `json:"rpcUrl"`
	ExplorerURL string `json:"explorerUrl"`
	// DeploymentBlock is the inclusive block the prediction market proxy was
	// deployed in: the first block an empty index scans.
	DeploymentBlock uint64 `json:"deploymentBlock"`
	Confirmations   uint64 `json:"confirmations"`

	PredictionMarket               *string `json:"predictionMarket"`
	PredictionMarketImplementation *string `json:"predictionMarketImplementation"`
	DIAResolver                    *string `json:"diaResolver"`
	USDW                           *string `json:"usdw"`
	Faucet                         *string `json:"faucet"`
}

// NativeCurrency is the chain's gas coin, as a wallet adds the network.
type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

func (config Config) validate() error {
	if config.Name == "" || config.RPCURL == "" || config.ExplorerURL == "" || config.NativeCurrency.Symbol == "" {
		return fmt.Errorf("name, rpcUrl, explorerUrl and nativeCurrency are required")
	}
	for _, address := range []*string{config.PredictionMarket, config.PredictionMarketImplementation, config.DIAResolver, config.USDW, config.Faucet} {
		if address != nil && !addressPattern.MatchString(*address) {
			return fmt.Errorf("address %q is not a 0x-prefixed 20-byte hex address", *address)
		}
	}
	if !config.Enabled {
		return nil
	}
	if config.PredictionMarket == nil || config.DeploymentBlock == 0 || config.Confirmations == 0 {
		return fmt.Errorf("an enabled chain needs predictionMarket, deploymentBlock and confirmations")
	}
	return nil
}

var addressPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

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
		if err := config.validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
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

// Enabled returns the enabled chain configs ordered by chain id.
func Enabled() ([]Config, error) {
	configs, err := All()
	if err != nil {
		return nil, err
	}
	enabled := configs[:0]
	for _, config := range configs {
		if config.Enabled {
			enabled = append(enabled, config)
		}
	}
	return enabled, nil
}
