<!--
order: 2
-->

# State

The `fantoken` module keeps track of [**parameters**](#Params) and [**fan tokens**](#Token).

```
Params:			types.Params
FanTokens:		[]types.FanToken
```

## Params

In the state definition, we can find the **Params**. This section corresponds to a module-wide configuration structure that stores system parameters. In particular, it defines the overall fantoken module functioning and contains the **issueFee**, **mintFee** and **burnFee** for the _fan token_. Such an implementation allows governance to decide the issue fee, but also the mint and burn fees the users have to pay to perform these operations with the tokens, in an arbitrary way - since proposals can modify it.

```go
type Params struct {
	IssueFee	sdk.Coin
	MintFee		sdk.Coin
	BurnFee		sdk.Coin
}
```

## Fantoken

The state contains a list of **Fantokens**. They are [fan tokens](01_concepts.md#Fan-token) (fungible tokens deriving by the ERC-20 Standard), and their state information is:

- **Denom**, that corresponds to the identifier of the fan token. It is a `string`, automatically calculated on the first `Minter`, `Symbol`, `Name` and `Block Height` of the issuing transaction of the _fan token_ as explained in [concepts](01_concepts.md#Fan-token), and _cannot change_ for the whole life of the token;
- **MaxSupply**, that represents the upper limit for the total supply of the tokens. More specifically, it is an `integer number`, expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})) as explained in [concepts](01_concepts.md#Fan-token), that _cannot change_ for the whole life of the token and which corresponds to the maximum number the supply can reach in any moment;
- **Minter**, which corresponds to the address of the current `minter` for the token. It is an address and _can change_ during the token lifecycle thanks to the **minting ability transfer**. When the `minter` address is set to an empty value, the token can be minted no more;
- **MetaData**, which contains metadata for the _fan token_ and is made up of the `Name`, the `Symbol`, a `URI` and an `Authority` as described in [concepts](01_concepts.md#Fan-token).

More specifically, the `metadata` _can change_ during the life of the token according to:
- **URI** can be changed by the `authority`. It can be changed until when the authority is available;
- **Authority** which can be transferred by the current authority until when the `authority` itself is not set to an empty value.

```go
type FanToken struct {
	Denom		string
	MaxSupply	sdk.Int
	Minter		string
	MetaData	types.Metadata
}

type Metadata struct {
	Name		string
	Symbol      string
	URI         string
	Authority	string
}
```