<!-- 
order: 3
-->

# Messages

Messages (`msg`s) are objects that trigger state transitions. Messages are wrapped in transactions (`tx`s) that clients submit to the network. The BitSong SDK wraps and unwraps `merkledrop` module messages from transactions.

## MsgCreate
The `MsgCreate` message is used to create a new _merkledrop_. It takes as input `Owner`, `MerkleRoot`, `StartHeight`, `EndHeight`, and `Coin`.

```go
type MsgCreate struct {
	Owner			string
	MerkleRoot		string
	StartHeight		int64
	EndHeight		int64
	Coin			sdk.Coin
}
```

## MsgClaim


```go
type MsgClaim struct {
	Sender			string
	MerkledropId	uint64
	Index			uint64
	Amount			sdk.Int 
	Proofs			[]string
}
```