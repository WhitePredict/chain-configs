package chainconfigs

import "testing"

func TestEveryConfigLoads(t *testing.T) {
	configs, err := All()
	if err != nil {
		t.Fatal(err)
	}
	if len(configs) == 0 {
		t.Fatal("no chain configs embedded")
	}
	config, err := Get(1874)
	if err != nil || config.PredictionMarket == nil {
		t.Fatalf("Whitechain config: %+v, %v", config, err)
	}
	enabled, err := Enabled()
	if err != nil || len(enabled) == 0 {
		t.Fatalf("enabled chains: %+v, %v", enabled, err)
	}
	for _, config := range enabled {
		if config.PredictionMarket == nil || config.DeploymentBlock == 0 || config.Confirmations == 0 {
			t.Fatalf("enabled chain %d is incomplete: %+v", config.ChainID, config)
		}
	}
	if _, err := Get(999999); err == nil {
		t.Fatal("unknown chain must fail")
	}
}

func TestValidateRejectsAnIncompleteEnabledChain(t *testing.T) {
	market := "0x47b42c74A5965a06E835bdaFaaC70bEbb48524bc"
	complete := Config{
		ChainID: 1, Name: "Chain", Enabled: true, RPCURL: "https://rpc", ExplorerURL: "https://explorer",
		NativeCurrency:  NativeCurrency{Name: "Ether", Symbol: "ETH", Decimals: 18},
		DeploymentBlock: 1, Confirmations: 1, PredictionMarket: &market,
	}
	if err := complete.validate(); err != nil {
		t.Fatal(err)
	}
	bad := "0x1234"
	for name, mutate := range map[string]func(*Config){
		"no market":        func(c *Config) { c.PredictionMarket = nil },
		"no block":         func(c *Config) { c.DeploymentBlock = 0 },
		"no confirmations": func(c *Config) { c.Confirmations = 0 },
		"no rpc":           func(c *Config) { c.RPCURL = "" },
		"bad address":      func(c *Config) { c.USDW = &bad },
	} {
		config := complete
		mutate(&config)
		if config.validate() == nil {
			t.Fatalf("%s: incomplete config accepted", name)
		}
	}
}
