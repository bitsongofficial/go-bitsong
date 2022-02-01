# Messages

## MsgIssueFanToken
A new fantoken is created using the `MsgIssueFanToken` message.

```go
type MsgIssueFanToken struct {
	Symbol		string
	Name		string
	MaxSupply	sdk.Int
	Mintable	bool
	Description string
	Owner		string
	IssueFee	sdk.Coin
}
```

# MsgEditFanToken
The `Mintable` of a fantoken can be updated using the `MsgEditFanToken`.

```go
type MsgEditFanToken struct {
	Denom		string
	Mintable	bool
	Owner		string
}
```

## MsgMintFanToken
Only the owner of the fantoken can mint new fantoken to a specified account. It fails if the total supply > max supply

```go
type MsgMintFanToken struct {
	Recipient	string
	Denom		string
	Amount		sdk.Int
	Owner		string
}
```

## MsgBurnFanToken
The action will be completed if the sender balance > amount to burn

```go
type MsgBurnFanToken struct {
	Denom		string
	Amount		sdk.Int
	Sender		string
}
```

## MsgTransferFanTokenOwner

Transfer the ownership of the fantoken to another account owner

```go
type MsgTransferFanTokenOwner struct {
	Denom		string
	SrcOwner	string
	DstOwner	string
}
```