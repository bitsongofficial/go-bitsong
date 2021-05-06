# Messages

## MsgIssueFanToken
A new token is created using the `MsgIssueFanToken` message.

```go
type MsgIssueFanToken struct {
	Denom		string
	Name		string
	MaxSupply	uint64
	Mintable	bool
	Owner		string
}
```

## MsgMintFanToken
Only the owner of the fan token can mint new token to a specified account. It fail if the total supply > max supply

```go
type MsgMintFanToken struct {
	Recipient	string
	Denom		string
	Amount		uint64
	Owner		string
}
```

## MsgBurnFanToken
The action will be completed if the sender balance > balance to burn

```go
type MsgBurnFanToken struct {
	Denom		string
	Amount		uint64
	Sender		string
}
```

## MsgTransferFanTokenOwner

Transfer the ownership of the Fan token to another account owner

```go
type MsgTransferFanTokenOwner struct {
	Denom		string
	SrcOwner	string
	DstOwner	string
}
```