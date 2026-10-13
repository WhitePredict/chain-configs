# chain-configs

Deployed WhitePredict contract addresses per chain. Each chain lives in
`<chainId>/config.json`; an address is `null` until the contract is deployed
there.

## Use from Go

The JSON files are embedded into the `chainconfigs` package, so a consumer
pins the addresses by module version:

```sh
go get github.com/WhitePredict/chain-configs@v0.1.0
```

```go
import chainconfigs "github.com/WhitePredict/chain-configs"

config, err := chainconfigs.Get(11155420)
```

## Releasing a change

1. Edit or add `<chainId>/config.json` and run `go test ./...`.
2. Merge to `main`, then tag a new version: `git tag v0.1.1 && git push origin v0.1.1`.
3. In each consumer run `go get github.com/WhitePredict/chain-configs@v0.1.1`.

A tag is immutable once the Go module proxy has seen it; publish a new tag
instead of moving an existing one.
