<!-- 
order: 3
-->

# Messages

Messages (`msg`s) are objects that trigger state transitions. Messages are wrapped in transactions (`tx`s) that clients submit to the network. The BitSong SDK wraps and unwraps `merkledrop` module messages from transactions.

## MsgCreate

The `MsgCreate` message is used to create a new _merkledrop_. It takes as input `Owner`, `MerkleRoot`, `StartHeight`, `EndHeight`, and `Coin`. The value of the block height at which the drop become available (the **starting** block) must be greater or equal to the block height where the transaction is included. For this reason, if the users select **0** as `StartHeight` it will be automatically set to the current block height (the one where the transaction is included). Moreover, there exists an upper bound for this value, that corresponds to the value of the `actual block height + 100000`. This choice derives from a design pattern that avoid the generation of _spam_ _merkledrop_. At the same time, the `EndHeight` value, which corresponds to the block height where the _merkledrop_ is closed and the withdrawal is executed if part of the tokens were not claimed. This value must be greater than the `StartHeight` and lower than a maximum value of `selected start block height + 5000000`. The `Coin` is made up of the `denom` of the token to distribute and the `amount`, which corresponds to the sum of all the tokens to drop. Once the module has verified that the `owner` address is valid and that the `merkletree root` is a hexadecimal character string, it **deduct the `creation fee` from the owner wallet** and send the `coin` (the amount of token to drop), from the owner address to the module. At this point, the `LastMerkleDropId` is increased and the _merkledrop_ is created, by assigning **zero to the claimed value** (since at the creation time, no one claimed any token). They are added three indexes:
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
The `MsgClaim` message is used to claim tokens from an active _merkledrop_. It takes as input `Sender`, `MerkledropId`, `Index`, the `Amount` to claim, and a list of `Proofs`. In such a scenario, verified the validity of the `sender` address and the existence of the _merkledrop_ by the ID, if the airdrop is currently active (i.e., its `start block height` is lower than the current block height and its `end block height` is greater than the current one), the module verifies if the `sender` already claimed his tokens (by querying at an index). In case he didn't, the module proceeds retriving the merkletree root for the _merkledrop_ from the chain and verifies the proofs (as described in the [verification process](01_concepts.md#Verification-process)). 
After tese verifications, the module only checks if the coin the `sender` wants to claim are available, and send those tokens from the module to the `sender` wallet. At this point, the claim is stored through its index, the claimed tokens are added to the actually claimed amount and, if all the drops are claimed with this operation, the merkledrop is cleaned by the state. 
An event of type `EventClaim` is emitted at the end of the claim process.

```go
type MsgClaim struct {
	Sender			string
	MerkledropId	uint64
	Index			uint64
	Amount			sdk.Int 
	Proofs			[]string
}
```