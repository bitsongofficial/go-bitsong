<!-- 
order: 3
-->

# Messages

Messages (`msg`s) are objects that trigger state transitions. Messages are wrapped in transactions (`tx`s) that clients submit to the network. The BitSong SDK wraps and unwraps `fantoken` module messages from transactions.

## MsgIssue
The `MsgIssue` message is used to issue a new fan token. It takes as input `Symbol`, `Name`, `MaxSupply` (expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})) as explained in [concepts](01_concepts.md#Fan-token)), `Authority`, and `URI` which is a link to the fan token metadata. Thanks to these values, the module can verify if the `Authority` can issue a new token (it is not a blocked address or module account). At this point, it proceeds with token issuing and emitting of corresponding events. More specifically, the module deduct the `issuing fee` from the owner wallet, calculates the `denom`, generates the `metadata`, and finally creates the fan token. At this point, an `EventTypeIssue` event is emitted.

```go
type MsgIssue struct {
	Symbol			string
	Name			string
	MaxSupply		sdk.Int
	Authority		string
	URI				string
}
```

## MsgDisableMint
The `MsgDisableMint` message is used to **irreversibly** disable the minting ability for an existing fan token. It takes as input `Denom` and `Authority` (described in [fan token definition](01_concepts.md#Fan-token)). Thanks to these values, the module can verify whether the modifications are lawful (i.e., requested by the `Authority` and in accord with the state transition definition). The message permits to change the "*mintability*" of the fan token. In particular, at the issuing, the fan token can be minted (in fact the `Mintable` bool is set to true). Later on, during the lifecycle of the fan token, the authority can disable the possibility to mint it (check the [relative docs](01_concepts.md#Lifecycle-of-a-fan-token) for more details). After disabling the mintability, the `MaxSupply` of the token is updated. At this point, an `EventTypeDisableMint` event is emitted.

```go
type MsgEdit struct {
	Denom			string
	Authority		string
}
```

## MsgMint

The `MsgMint` message is used to mint an existing fan token. It takes as input `Recipient`, `Denom`, `Amount`, and `Authority` (all described in [fan token definition](01_concepts.md#Fan-token) except the `Amount`, which is the quantity of token to mint, expressed in micro unit). In such a message, the `Recipient` is not required and its default value is the same of `Authority`. 
Thanks to these values, the module can verify whether the minting operation is lawful (i.e., requested: by the authority, on a mintable fan token, and for a quantity that allow to do not overcome the maximum supply), recalling that only the authority for of the fan token can mint the token to a specified account. 
At this point, the token is minted, the result is sent to the recipient, and an `EventTypeMint` event is emitted.

```go
type MsgMintFanToken struct {
	Recipient		string
	Denom			string
	Amount			sdk.Int
	Authority		string
}
```

## MsgBurn

The `MsgBurn` message is used to burn fan token. It takes as input `Denom`, `Amount`, and `Sender` (`Denom` is described in [fan token definition](01_concepts.md#Fan-token), `Amount` is the quantity of token to burn, and `Sender` must be equal to the user who want to burn the tokens).
The module can verify whether the burning operation is lawful (i.e., the sender has a sufficient amount of token, in other words check if `sender balance` > `amount to burn`). At this point, the token is burned and an `EventTypeBurnFanToken` event is emitted.
In such a way, that specific token ends its lifecycle, as shown in the [relative docs](01_concepts.md#Lifecycle-of-a-fan-token).

```go
type MsgBurn struct {
	Denom			string
	Amount			sdk.Int
	Sender			string
}
```

## MsgTransferAuthority

The `MsgTransferAuthority` message is used to transfer the ownership of a fan token. It takes as input `Denom`, `SrcAuthority`, and `DstAuthority` (`Denom` is described in [fan token definition](01_concepts.md#Fan-token), `Src` and `Dst` `Authorities` are respectively the "*old*" and "*new*" authorities of the token).

The module can verify whether the operation is lawful (i.e., the requesting account is actually the authority for the fan token and the destination account is neither blocked nor a module account). 
At this point, the `DstAuthority` becomes the new token authority and an `EventTypeTransferAuthority` event is emitted.
This operation enable the **ownership transfer** transition described in the [lifecycle of a fan token documentation](01_concepts.md#Lifecycle-of-a-fan-token).

```go
type MsgTransferAuthority struct {
	Denom			string
	SrcAuthority	string
	DstAuthority	string
}
```