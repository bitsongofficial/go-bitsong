<!-- 
order: 3
-->

# Messages

Messages (`msg`s) are objects that trigger state transitions. Messages are wrapped in transactions (`tx`s) that clients submit to the network. The BitSong SDK wraps and unwraps `merkledrop` module messages from transactions.

## MsgCreate

The `MsgCreate` message is used to create a new _merkledrop_. It takes as input `Owner`, `MerkleRoot`, `StartHeight`, `EndHeight`, and `Coin`. The value of the block height at which the drop become available (the **starting** block) must be greater or equal to the block height where the transaction is included. For this reason, if the users select **0** as `StartHeight` it will be automatically set to the current block height (the one where the transaction is included). Moreover, there exists an upper bound for this value, that corresponds to the value of the `actual block height + 100000`. This choice derives from a design pattern that avoid the generation of _spam_ merkledrop. At the same time, the `EndHeight` value, which corresponds to the block height where the merkledrop is closed and the withdrawal is executed if part of the tokens were not claimed. This value must be greater than the `StartHeight` and lower than a maximum value of `selected start block height + 5000000`. The `Coin` is made up of the `denom` of the token to distribute and the `amount`, which corresponds to the sum of all the tokens to drop. Once the module has verified that the `owner` address is valid and that the `merkletree root` is a hexadecimal character string, it **deduct the `creation fee` from the owner wallet** and send the `coin` (the amount of token to drop), from the owner address to the module. At this point, the `LastMerkleDropId` is increased and the merkledrop is created, by assigning **zero to the claimed value** (since at the creation time, no one claimed any token). They are added three indexes:
- on the `merkledrop_id`;
- on the `owner`;
- on the `end_height`.

These indexes improve the query operations and some process described in the [end_block operations](04_end_block.md).
An event of type `EventCreate` is emitted at the end of the creation process.

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