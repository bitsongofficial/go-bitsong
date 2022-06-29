<!--
order: 2
-->

# State

The `fantoken` module keeps track of [**parameters**](#Params), [**fan tokens**](#Token), and [**burned coins**](#BurnedCoins).

```
Params:      types.Params
Tokens:      []types.FanToken
BurnedCoins: []sdk.Coin
```

## Params

In the state definition, we can find the **Params**. This section corresponds to a module-wide configuration structure that stores system parameters. In particular, it defines the overall fantoken module functioning and contains the **issueFee**, **mintFee**, **burnFee**, and the **transferFee** for the fan token. Such an implementation allows governance to decide the issue fee, but also the mint, burn and transfer fees, the users have to pay to perform these operations with the tokens, in an arbitrary way - since proposals can modify it.

```go
type Params struct {
	IssueFee    sdk.Coin
    MintFee     sdk.Coin
    BurnFee     sdk.Coin
    TransferFee sdk.Coin
}
```

## Fantoken

The state contains a list of **Fantokens**. They are [fan tokens](01_concepts.md#Fan-token) (fungible tokens deriving by the ERC-20 Standard), and their state information are:

- **Denom**, that corresponds to the Identifier of the fan token. It is a `string`, automatically calculated on `Owner`, `Symbol`, `Name` and `Block Height` of the issuing transaction of the fan token as explained in [concepts](01_concepts.md#Fan-token), and _cannot change_ for the whole life of the token;
- **MaxSupply**, that represents the maximum number of possible mintable tokens. It is an `integer number`, expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})) as explained in [concepts](01_concepts.md#Fan-token), and _cannot change_ for the whole life of the token;
- **Mintable**, indicating the ability of the token to be minted. It is a `boolean` value and \*can change **only once\*** during the token lifecycle. At the issuing it is set to true, and the token can be minted. When the owner change this value in the state, the token can be minted no more;
- **Owner**, which is the current owner of the token. It is an address and _can change_ during the token lifecycle thanks to the **ownership transfer**;
- **MetaData**, which contains metadata for the fan token and is made up of the `Name`, the `Symbol`, and a `URI`, as described in [concepts](01_concepts.md#Fan-token), and _cannot change_ for the whole life of the token.

```go
type FanToken struct {
	Denom		string
	MaxSupply	sdk.Int
	Mintable	bool
	Authority	string
	MetaData	types.Metadata
}

type Metadata struct {
	Name		string
    Symbol      string
	URI         string
}
```

## BurnedCoins

Another section in this module state is represented by **BurnedCoins**. It contains the total amount of all the burned tokens.

```go
BurnedCoins: []sdk.Coin
```

```go
type Coin struct {
	Denom  string
	Amount Int
}
```