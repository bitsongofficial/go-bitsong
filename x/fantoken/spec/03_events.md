# Events

The fantoken module emits the following events:
## MsgIssueFanToken

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| issue_fan_token | denom         | {Denom}          |
| issue_fan_token | symbol        | {Symbol}         |
| issue_fan_token | creator       | {creatorAddress} |
| message         | module        | fantoken         |
| message         | sender        | {ownerAddress}   |

## MsgEditFanToken

| Type            | Attribute Key | Attribute Value |
| :-------------- | :------------ | :-------------- |
| edit_fan_token  | denom         | {Denom}         |
| edit_fan_token  | owner         | {ownerAddress}  |
| message         | module        | fantoken        |
| message         | sender        | {ownerAddress}  |

## MsgTransferFanTokenOwner

| Type                     | Attribute Key | Attribute Value   |
| :----------------------- | :------------ | :---------------- |
| transfer_fan_token_owner | denom         | {Denom}           |
| transfer_fan_token_owner | owner         | {ownerAddress}    |
| transfer_fan_token_owner | dst_owner     | {dstOwnerAddress} |
| message                  | module        | fantoken          |
| message                  | sender        | {ownerAddress}    |

## MsgMintFanToken

| Type           | Attribute Key | Attribute Value    |
| :------------- | :------------ | :----------------- |
| mint_fan_token | denom         | {Denom}            |
| mint_fan_token | amount        | {amount}           |
| mint_fan_token | recipient     | {recipientAddress} |
| message        | module        | fantoken           |
| message        | sender        | {ownerAddress}     |

## MsgBurnToken

| Type           | Attribute Key | Attribute Value |
| :------------- | :------------ | :-------------- |
| burn_fan_token | denom         | {Denom}         |
| burn_fan_token | amount        | {amount}        |
| message        | module        | fantoken        |
| message        | sender        | {ownerAddress}  |
