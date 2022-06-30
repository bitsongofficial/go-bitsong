<!-- 
order: 4
-->

# Events

The fantoken module emits the following events:
## EventIssue

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgIssue` |
| bitsong.fantoken.v1beta1.EventIssue | denom        | {denom}         |

## EventDisableMint

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgDisableMint` |
| bitsong.fantoken.v1beta1.EventDisableMint | denom        | {denom}         |

## EventMint

| Type                     | Attribute Key | Attribute Value   |
| :----------------------- | :------------ | :---------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgMint` |
| bitsong.fantoken.v1beta1.EventMint | recipient        | {recipient}         |
| bitsong.fantoken.v1beta1.EventMint | coin        | {coin}         |

## EventBurn

| Type           | Attribute Key | Attribute Value    |
| :------------- | :------------ | :----------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgBurn` |
| bitsong.fantoken.v1beta1.EventBurn | sender        | {sender}         |
| bitsong.fantoken.v1beta1.EventBurn | coin        | {coin}         |

## EventSetAuthority

| Type           | Attribute Key | Attribute Value |
| :------------- | :------------ | :-------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgSetAuthority` |
| bitsong.fantoken.v1beta1.EventTransferAuthority | denom        | {denom}         |
| bitsong.fantoken.v1beta1.EventTransferAuthority | old_authority        | {old_authority}         |
| bitsong.fantoken.v1beta1.EventTransferAuthority | new_authority        | {new_authority}         |

## EventSetMinter

| Type           | Attribute Key | Attribute Value |
| :------------- | :------------ | :-------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgSetMinter` |
| bitsong.fantoken.v1beta1.EventTransferMinter | denom        | {denom}         |
| bitsong.fantoken.v1beta1.EventTransferMinter | old_minter        | {old_minter}         |
| bitsong.fantoken.v1beta1.EventTransferMinter | new_authority        | {new_minter}         |

## EventSetUri

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgSetUri` |
| bitsong.fantoken.v1beta1.EventSetUri | denom        | {denom}         |
