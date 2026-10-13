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
	if _, err := Get(1); err == nil {
		t.Fatal("unknown chain must fail")
	}
}
