<!-- 
order: 3
-->

# Messages

Messages (`msg`s) are objects that trigger state transitions. Messages are wrapped in transactions (`tx`s) that clients submit to the network. The BitSong SDK wraps and unwraps `fantoken` module messages from transactions.

## MsgIssue
The `MsgIssue` message is used to issue a new _fan token_. It takes as input `Symbol`, `Name`, `MaxSupply` (expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})) as explained in [concepts](01_concepts.md#Fan-token)), `Authority` (i.e., the address of the wallet which is able to modify the `metadata` of the _fan token_), `URI` (which is a link to the `fan token` metadata) and the `Minter` (i.e., the address of the wallet which is able to mint the _fan token_). Thanks to these values, the module can verify if the `Authority` and the `Minter` are valid addresses for the issue of a new token (they are not a blocked addresses or module accounts) and also verifies the values for the `name` (which can be any strings with max 128 characters, even the empty one), the `symbol` (that must match the regex `^[a-z0-9]{1,64}$`) and the `uri` (which can be any strings with less than 513 characters, even the empty one). At this point, it proceeds with token issuing and emitting of corresponding events. More specifically, the **module deduct the `issuing fee` from the `minter` wallet**, calculates the `denom`, generates the `metadata`, and finally creates the _fan token_. At this point, an `EventIssue` event is emitted.

```go
type MsgIssue struct {
	Symbol			string
	Name			string
	MaxSupply		sdk.Int
	Authority		string
	URI				string
	Minter			string
}
```

## MsgDisableMint
The `MsgDisableMint` message is used to **irreversibly** disable the minting ability for an existing _fan token_. It takes as input `Denom` and `Minter` (described in [fan token definition](01_concepts.md#Fan-token)). Thanks to these values, the module can verify whether the modifications are lawful (i.e., requested by the `Minter` and in accord with the state transition definition). The message permits to change the "*mintability*" of the _fan token_. In particular, at the issuing, the _fan token_ can be minted (in fact the `Minter` address is a value different from an empty one). Later on, during the lifecycle of the _fan token_, the `minter` can disable the possibility to mint new tokens (check the [relative docs](01_concepts.md#Lifecycle-of-a-fan-token) for more details). In such a scenario, it is possible to disable the mintability, by set an empty value as the address for the new `minter`) and, this operation, causes the `MaxSupply` of the token to be updated at the current value of the supply. At this point, an `EventDisableMint` event is emitted.

```go
type MsgDisableMint struct {
	Denom			string
	Minter			string
}
```

## MsgMint

The `MsgMint` message is used to mint an existing _fan token_. It takes as input `Recipient`, `Coin`, and `Minter` (all described in [fan token definition](01_concepts.md#Fan-token) except the `Coin`, which is an object made up of the `denom` of the _fan token_ to mint and its quantity, expressed in micro unit). In such a message, the `Recipient` is not required and its default value is the same of `Minter`. 
Thanks to these values, the module can verify whether the minting operation is lawful (i.e., requested: by the minter, on a mintable _fan token_, and for a quantity that allow to do not overcome the maximum supply), recalling that only the minter for of the _fan token_ can mint the token to any specified account.
At this point, the token is minted, the supply is increased, the coins are sent to the recipient, the **module deduct the `mint fee` from the `minter` wallet** and an `EventMint` event is emitted.

```go
type MsgMint struct {
	Recipient		string
	Coin			sdk.Coin
	Minter			string
}
```

## MsgBurn

The `MsgBurn` message is used to burn _fan token_. It takes as input `Coin`, and `Sender` (as above, the `Coin` is an object made up of the `denom` of the _fan token_ to burn and its quantity, expressed in micro unit, while `Sender` must be equal to the user who want to burn the tokens).
The module can verify whether the burning operation is lawful (i.e., the sender has a sufficient amount of token, in other words check if `sender balance` > `amount to burn`). At this point, the token is burned, the supply is lowered, the **module deduct the `burn fee` from the `owner` wallet** and an `EventBurn` event is emitted.
In such a way, that specific token ends its lifecycle, as shown in the [relative docs](01_concepts.md#Lifecycle-of-a-fan-token).

```go
type MsgBurn struct {
	Coin			sdk.Coin
	Sender			string
}
```

## MsgSetAuthority

The `MsgSetAuthority` message is used to transfer or disable the ability to change the metadata of a _fan token_. It takes as input `Denom`, `oldAuthority`, and `newAuthority` (`Denom` is described in [fan token definition](01_concepts.md#Fan-token), `old` and `new` `Authorities` are respectively the actual and the new addresses of the wallet who are able to change the metadata of the token). When the `newAuthority` is an empty address, the capability to change the metadata is **irreversibly** disabled.
The module can verify whether the operation is lawful (i.e., the requesting account is actually the authority for the _fan token_, the _fan token_ metadata can be changed and the destination account is neither blocked nor a module account). 
At this point, if the `newAuthority` is a _not empty_ address, it becomes the new token authority. On the other hand, the _fan token_ metadata cannot be changed anymore. Anyway, an `EventSetAuthority` event is emitted.
This operation enable the **authority transfer** transition described in the [lifecycle of a fan token](01_concepts.md#Lifecycle-of-a-fan-token).

```go
type MsgTransferAuthority struct {
	Denom			string
	oldAuthority	string
	newAuthority	string
}
```

## MsgSetMinter

The `MsgSetMinter` message is used to transfer the ability to mint a _fan token_. It takes as input `Denom`, `oldMinter`, and `newMinter` (`Denom` is described in [fan token definition](01_concepts.md#Fan-token), `old` and `new` `Minters` are respectively the actual and the new addresses of the wallet who are able to mint the token). When the `newMinter` is an empty address, it works as the `MsgDisableMint`.
The module can verify whether the operation is lawful (i.e., the requesting account is actually the minter for the _fan token_, the _fan token_ can be minted and the destination account is neither blocked nor a module account). 
At this point, if the `newMinter` is a _not empty_ address, it becomes the new token minter. On the other hand, the _fan token_ cannot be minted anymore. Anyway, an `EventSetMinter` event is emitted.
This operation enable the **minter transfer** transition described in the [lifecycle of a fan token](01_concepts.md#Lifecycle-of-a-fan-token).

```go
type MsgSetMinter struct {
	Denom		string
	oldMinter	string
	newMinter	string
}
```

## MsgSetUri

The `MsgSetMinter` message is used to modify the URI in the _fan token_ metadata. It takes as input `Denom`, new `URI`, and `Authority` (`Denom` and `URI` are described in [fan token definition](01_concepts.md#Fan-token), `Authority` is the actual address of the wallet who is able to modify the _fan token_ metadata).
The module can verify whether the operation is lawful (i.e., the requesting account is actually the authority for the _fan token_, the _fan token_ metadata can be changed and the new uri is a valid one, as described in [Fan Token parameters definition](01_concepts.md#Fan-token)). 
At this point, an `EventSetUri` event is emitted.

```go
type MsgSetUri struct {
	Denom			string
	URI				string
	Authority		string
}
```