<!--
order: 1
-->

# Concepts

## Conventions

By looking at numbers, we separate the decimals by point and the thousands by comma. For instance, the number _one thousand two hundred thirty-four and fifty-six hundredths_, is written as:

![formula](https://render.githubusercontent.com/render/math?math=\color{gray}1,234.56)

## Fan token

Fan tokens, conceptually based on the [ERC-20 Standard](https://ethereum.org/it/developers/docs/standards/tokens/erc-20), are **fungible tokens** issued for fan communities. They borns to create new connections between fans and any content creator, like star performers, actors, designers, musicians, photographers, writers, models, influencers, etc.
They enable the growth of a private and (most importantly) custom economy creating new channels for fans' engagement.
_Fan tokens_ have enormous potential. By using them, you can build myriad applications allowing fans a deeper interaction in the artistic life of their top performers.

To provide you with some examples, you can think that it is possible to use them for creating loyalty programs to provide privileged access to exclusive content. To allow your fan to crowdfund a tour or studio album and share part of the revenue with your fans. To enable your fans with the opportunity to vote on the cities for an upcoming tour. Or even to accept _fan tokens_ as payment for NFTs.

In the design of the _fan token_ functionalities, big part of the reasonings were based on the [OpenZeppelin standard](https://docs.openzeppelin.com/contracts/4.x/api/token/erc20). For example, the concept of *burning* the tokens lowering the `totalSupply` directly derives from the standard [documentation](https://docs.openzeppelin.com/contracts/4.x/api/token/erc20#ERC20-_burn-address-uint256-).

A **fan token** is characterized by:
| Attribute | Type | Description |
| --------------------- | ---------------------------- | ---------------------------------------------- |
| denom | `string` | It is an hash calculated on the first `Minter`, the `Symbol`, the `Name` and the `Block Height` of the issuing transaction of the _fan token_. It is the hash identifying the _fan token_ and is used to [prevent the creation of identical tokens](#Uniqueness-of-the-denom). Moreover, to fastly identify a _fan token_ from its `denom`, it starts with the prefix `ft`.|
| max_supply | `sdk.Int` | It is chosen once by the user. It is the maximum supply value of mintable tokens from its definition. It is expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})). For this reason, to indicate a maximum supply of ![formula](https://render.githubusercontent.com/render/math?math=\color{gray}456) tokens, this value must be equal to ![formula](https://render.githubusercontent.com/render/math?math=\color{gray}456\cdot10^{6}=456,000,000).|
| Minter | `sdk.AccAddress` | It is the address of the minter for the _fan token_. It can be changed to trasfer the minting ability of the token during the time. |
| metadata | `Metadata` | It is generated once and it is made up of `Name`, `Symbol`, `URI` and `Authority` (i.e., is the address of the wallet which is able to perform edits on the `URI`). More specifically, the URI contains a link to a resource with a set of information linked to the _fan token_.|

**Metadata** are characterized by:
| Attribute | Type | Description |
| --------------------- | ---------------------------- | ---------------------------------------------- |
| name | `string` | It is chosen once by the user. It should correspond to the long name the user want to associate to the symbol (e.g., Dollar, Euro, BitSong). It can also be empty and its max length is of 128 characters. |
| symbol | `string` | It is chosen once by the user and can be any string matching the pattern `^[a-z0-9]{1,64}$`, i.e., any lowercase string containing letters and digits with a length between 1 and 64 characters. It should follow the ISO standard for the [alphabetic code](https://www.iso.org/iso-4217-currency-codes.html) (e.g., USD, EUR, BTSG, etc.).|
| uri | `string` | It is a link to a resource which contains a set of information linked to the _fan token_. |
| authority | `sdk.AccAddress` | It is the address of the authority for the _fan token_ `metadata` managment. It can be changed to trasfer the ability of changing the metadata the token during the time. |


## Lifecycle of a fan token

It is possible to entirely represent the lifecycle of a fan token through Finite State Machine (FSM) diagrams. We will present two representations:

- the first refers to the fan token **object**. We can compare such a definition with that of currency (e.g., Euro, Dollar, BitSong);
- the second, instead, is referred to the lifecycle of the fan token **instance**. Such definition is comparable with that of coin/money (e.g., the specific 1 Euro coin you could have in your pocket at a particular moment in time).

We can describe the lifecycle of a fan token **object** through two states.

![Fantoken object lifecycle](img/fantoken_object_lifecycle.png "Fantoken object lifecycle")

Referring to the figure above, as detailed in the documentation, to "create" the fan token, we need to **issue it**. This operation leads to the birth of the object and thus to its first state, state _1_. Here, the token is related to a `minter`, who is able to mint the token to different wallets, and an `authority`, that is responsible for managing the `metadata`. From this state, it is possible:

- to **transfer the ability to mint to a new address**, which produces the changing of the minter, without modifying the landing state. This operation can be done on every state until when the minter address is set to empty;
- to **transfer the ability to edit the metadata to a new address**, which produces the changing of the authority, without modifying the landing state. This operation can be done on every state until when the authority address is set to empty;
- to **disable the minting ability**, which is achived by setting the `minter` address to an empty one. This produces a state change to the state _2_. Here, no one can mint the _fan token_ anymore.
Once the _fan token_ lands in state _2_, the only possible action is to transfer its authority (the ability to change its `metadata`) to another address or to disable this feature landing the _fan token_ to the state _4_. Once the `minter` address is set to empty, the minting **ability can be enabled** no more;
- to **disable the ability to edit the metadata**, which is achived by setting the `authority` address to an empty one. This produces a state change to the state _3_. Here, no one can manage the _fan token_ `metadata` anymore.
Once the _fan token_ lands in state _3_, the only possible action is to transfer its `minter` address (that corresponds to the address of who is able to mint the `fan token`) to another address or to disable this feature landing the _fan token_ to the state _4_. Once the `authority` address is set to empty, the **ability to change the fan token metadata can be enabled** no more.

Also referring to the lifecycle of a fan token **instance**, it is possible to identify two states.

![Fantoken instance lifecycle](img/fantoken_instance_lifecycle.png "Fantoken instance lifecycle")

Concerning to the figure above, when the fan token object is issued, we can **mint** it. Minting leads to the birth of a new instance, moving the fan token instance to state _1_. In this state, the token can be:

- **traded**, which produces the changing of the owner of the instance, without modifying the landing state. To make it clearer, it can be considered as the simple exchange of money between two users. This does not modify the landing state;
- **burned**, which produces a state change to the state _2_, where the authority cannot operate on the fan token instance anymore.

## Uniqueness of the denom

The _denom_ is calculated on first `Minter`, `Symbol`, `Name` and `Block Height` of the issuing transaction of the fan token. 

```go
func GetFantokenDenom(height int64, minter sdk.AccAddress, symbol, name string) string {
	bz := []byte(fmt.Sprintf("%d%s%s%s", height, minter.String(), symbol, name))
	return "ft" + tmcrypto.AddressHash(bz).String()
}
```

The _denom_ of every fan token starts with the prefix `ft`. Follows a **hash** of `Block Height`, first `Minter`, `Symbol` and `Name` of the _fan token_. This _denom_ is used as base denom for the fan token, and, for this reason, it should be **unique**. In this sense, since the hash depends both on the first `Minter` and the `Block Height`, multiple fan tokens with the same name and symbol can co-exist even created by the same address but they must be created from transactions in different blocks.
