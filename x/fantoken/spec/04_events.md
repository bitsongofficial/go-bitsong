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

### EventDisableMint

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgDisableMint` |
| bitsong.fantoken.v1beta1.EventDisableMint | denom        | {denom}         |

### EventMint

| Type                     | Attribute Key | Attribute Value   |
| :----------------------- | :------------ | :---------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgMint` |
| bitsong.fantoken.v1beta1.EventMint | denom        | {denom}         |
| bitsong.fantoken.v1beta1.EventMint | amount        | {amount}         |
| bitsong.fantoken.v1beta1.EventMint | recipient        | {recipient}         |

### EventBurn

| Type           | Attribute Key | Attribute Value    |
| :------------- | :------------ | :----------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgBurn` |
| bitsong.fantoken.v1beta1.EventBurn | denom        | {denom}         |
| bitsong.fantoken.v1beta1.EventBurn | amount        | {amount}         |

### EventTransferAuthority

| Type           | Attribute Key | Attribute Value |
| :------------- | :------------ | :-------------- |
| message         | action        | `/bitsong.fantoken.v1beta1.MsgTransferAuthority` |
| bitsong.fantoken.v1beta1.EventTransferAuthority | denom        | {denom}         |
| bitsong.fantoken.v1beta1.EventTransferAuthority | src_authority        | {src_authority}         |
| bitsong.fantoken.v1beta1.EventTransferAuthority | dst_authority        | {dst_authority}         |