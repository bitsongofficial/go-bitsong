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
Fan tokens have enormous potential. By using them, you can build myriad applications allowing fans a deeper interaction in the artistic life of their top performers.

To provide you with some examples, you can think that it is possible to use them for creating loyalty programs to provide privileged access to exclusive content. To allow your fan to crowdfund a tour or studio album and share part of the revenue with your fans. To enable your fans with the opportunity to vote on the cities for an upcoming tour. Or even to accept fan tokens as payment for NFTs.

A **fan token** is characterized by:
| Attribute | Type | Description |
| --------------------- | ---------------------------- | ---------------------------------------------- |
| denom | `string` | It is an hash calculated on `Authority`, `Symbol`, `Name` and `Block Height` of the issuing transaction of the fan token. It is the hash identifying the fan token and is used to [prevent the creation of identical tokens](#Uniqueness-of-the-denom). |
| max_supply | `sdk.Int` | It is chosen once by the user. It is the maximum supply value of mintable tokens from its definition. It is expressed in micro unit (![formula](https://render.githubusercontent.com/render/math?math=\color{gray}\mu=10^{-6})). For this reason, to indicate a maximum supply of ![formula](https://render.githubusercontent.com/render/math?math=\color{gray}456) tokens, this value must be equal to ![formula](https://render.githubusercontent.com/render/math?math=\color{gray}456\cdot10^{6}=456,000,000).|
| mintable | `bool` | It is `true` at issuing. It can be later changed by the authority. If it is `true`, the fan token authority can mint the token. |
| authority | `sdk.AccAddress` | It is the address of the authority for the fan token. It can be changed to trasfer the authorityship of the token during the time. It is mainly used to verify the ability to perform operations.|
| metadata | `Metadata` | It is generated once and it is made up of `Name`, `Symbol`, and `URI`. More specifically, the URI contains a link to a resource with a set of information linked to the fan token.|

**Metadata** are characterized by:
| Attribute | Type | Description |
| --------------------- | ---------------------------- | ---------------------------------------------- |
| name | `string` | It is chosen once by the user. It should correspond to the long name the user want to associate to the symbol (e.g., Dollar, Euro, BitSong). |
| symbol | `string` | It is chosen once by the user and can be any string matching the pattern `^[a-z0-9]{1,64}$`, i.e., any lowercase string containing letters and digits with a length between 1 and 64 characters.It should follow the ISO standard for the [alphabetic code](https://www.iso.org/iso-4217-currency-codes.html) (e.g., USD, EUR, BTSG, etc.).|
| uri | `string` |  It is a link to a resource which contains a set of information linked to the fan token. |


## Lifecycle of a fan token

It is possible to entirely represent the lifecycle of a fan token through Finite State Machine (FSM) diagrams. We will present two representations:

- the first refers to the fan token **object**. We can compare such a definition with that of currency (e.g., Euro, Dollar, BitSong);
- the second, instead, is referred to the lifecycle of the fan token **instance**. Such definition is comparable with that of coin/money (e.g., the specific 1 Euro coin you could have in your pocket at a particular moment in time).

We can describe the lifecycle of a fan token **object** through two states.

![Fantoken object lifecycle](img/fantoken_object_lifecycle.png "Fantoken object lifecycle")

Referring to the figure above, as detailed in the documentation, to "create" the fan token, we need to **issue it**. This operation leads to the birth of the object and thus to its first state, state _1_. Here, the token is related to an authority, which can mint it. From this state, the authority can perform two actions on the object:

- to **transfer the authorityship**, which produces the changing of the authority, without modifying the landing state;
- to **disable the minting ability**, which produces a state change to the state _2_. Here, the authority cannot mint the fan token anymore.

Once the fan token lands in state _2_, the only possible action is to transfer its authorityship. Here, the authority **can enable the minting ability** no more.

Also referring to the lifecycle of a fan token **instance**, it is possible to identify two states.

![Fantoken instance lifecycle](img/fantoken_instance_lifecycle.png "Fantoken instance lifecycle")

Concerning to the figure above, when the fan token object is issued, we can **mint** it. Minting leads to the birth of a new instance, moving the fan token instance to state _1_. In this state, the token can be:

- **traded**, which produces the changing of the authority, without modifying the landing state;
- **burned**, which produces a state change to the state _2_, where the authority cannot operate on the fan token instance anymore.

## Uniqueness of the denom

The _denom_ is calculated on `Authority`, `Symbol`, `Name` and `Block Height` of the issuing transaction of the fan token. 

```go
func GetFantokenDenom(height int64, authority sdk.AccAddress, symbol, name string) string {
	bz := []byte(fmt.Sprintf("%d%s%s%s", height, authority.String(), symbol, name))
	return "ft" + tmcrypto.AddressHash(bz).String()
}
```

The _denom_ of every fan token starts with the prefix `ft`. Follows a **hash** of `Block Height`, `Authority`, `Symbol` and `Name` of the fan token. This _denom_ is used as base denom for the fan token, and, for this reason, it should be **unique**. In this sense, since the hash depends both on the `Authority` and the `Block Height`, multiple fan tokens with the same name and symbol can co-exist even created by the same address.
