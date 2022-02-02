# State

## FanToken
Definition of data structure of Fungible Token

```go
type FanToken struct {
	Name		string
	MaxSupply	sdk.Int
	Mintable	bool
	Owner		string
	MetaData	bank.Metadata
}
```

## Params
Params is a module-wide configuration structure that stores system parameters and defines overall functioning of the fan token module.

```go
type Params struct {
	IssuePrice	sdk.Coin
}
```