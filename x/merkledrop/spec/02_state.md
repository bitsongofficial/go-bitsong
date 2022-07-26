<!--
order: 2
-->

# State

The `merkledrop` module keeps track of [**LastMerkledropId**](#LastMerkledropId), [**Merkledrops**](#Merkledrops), [**Indexes**](#Indexes) and [**Parameters**](#Params).

```
LastMerkledropId: Uint64,
Merkledrops:      []types.Merkledrop,
Indexes:          []*types.Indexes,
Params:           types.Params
```

## LastMerkledropId

This value is an integer that corresponds to the number of _merkledrops_ already created. It is used at the creation of a new merkledrop as its id.

## Merkledrops

The state contains a list of **Merkledrops**. They are [airdrop configuration](01_concepts.md#Merkledrop), and their state information is:

- **Id** that corresponds to the identifier of the _merkledrop_. It is an `uint64`, automatically incremented everytime a new merkledrop is created;
- **MerkleRoot**, that represent the root hash (in hex format) of the _merkle tree_ containing the data of the airdrop;
- **StartHeight**, that is the block height value at which the drop allows the user to claim the tokens;
- **EndHeight**, which corresponds to the block height value where the _merkledrop_ is considered expired and an automatic withdrawal is executed if part of the tokens were not claimed;
- **Denom**, which corresponds to the `denom` of the token to drop;
- **Amount**, that is the total `amount` of token to drop;
- **Claimed** which corresponds to the value of claimed tokens from the users. At the beginning it is 0 and is increased at each claim;
- **Owner** which is to the address of the wallet which is creating the _merkledrop_.

```go
type Merkledrop struct {
	Id 			uint64
	MerkleRoot 	string
	StartHeight int64
	EndHeight 	int64
	Denom 		string
	Amount 		sdk.Int
	Claimed 	sdk.Int
	Owner 		string
}
```

## Indexes
To perform the check operations, a list of index is also stored in the state for each merkledrop.

```go
type Indexes struct {
	MerkledropId uint64
	Index        []uint64
}
```


## Params

In the state definition, we can find the **Params**. This section corresponds to a module-wide configuration structure that stores system parameters. In particular, it defines the overall merkledrop module functioning and contains the **creationFee** for the _merkledrop_. Such an implementation allows governance to decide the creation fee, in an arbitrary way - since proposals can modify it.

```go
type Params struct {
	CreationFee	sdk.Coin
}
```