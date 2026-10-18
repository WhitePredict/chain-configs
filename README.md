# chain-configs

The WhitePredict chain registry. Each chain lives in `<chainId>/config.json`
and is the single place a chain is added: the Starbase API serves it to the
frontend through `GET /api/v1/chains`, and one Starbase indexer process
indexes every enabled chain.

| Field | Meaning |
| --- | --- |
| `chainId`, `name`, `testnet` | Chain identity. |
| `enabled` | Served and indexed. An enabled chain needs `predictionMarket`, `deploymentBlock` and `confirmations`. |
| `nativeCurrency` | Gas coin, as a wallet adds the network. |
| `rpcUrl`, `explorerUrl` | Public endpoints. A private RPC is a deployment secret and stays out of this repository. |
| `deploymentBlock` | Inclusive block of the PredictionMarket proxy deployment. |
| `confirmations` | Blocks an event waits before it is indexed. |
| `predictionMarket`, `predictionMarketImplementation`, `diaResolver`, `usdw`, `faucet` | Contract addresses; `null` until deployed there. |

Everything here is public: never add RPC credentials or keys.

## Use from Go

The JSON files are embedded into the `chainconfigs` package, so a consumer
pins the addresses by module version:

```sh
go get github.com/WhitePredict/chain-configs@v0.2.0
```

```go
import chainconfigs "github.com/WhitePredict/chain-configs"

config, err := chainconfigs.Get(11155420)
enabled, err := chainconfigs.Enabled()
```

## Releasing a change

1. Edit or add `<chainId>/config.json` and run `go test ./...`.
2. Merge to `main`, then tag a new version: `git tag v0.1.1 && git push origin v0.1.1`.
3. In each consumer run `go get github.com/WhitePredict/chain-configs@v0.1.1`.

A tag is immutable once the Go module proxy has seen it; publish a new tag
instead of moving an existing one.
