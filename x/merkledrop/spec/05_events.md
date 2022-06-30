<!-- 
order: 5
-->

# Events

The merkledrop module emits the following events:
## EventCreate

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.merkledrop.v1beta1.MsgCreate` |
| bitsong.merkledrop.v1beta1.EventCreate | owner        | {owner}         |
| bitsong.merkledrop.v1beta1.EventCreate | merkledrop_id        | {merkledrop_id}         |

## EventClaim

| Type            | Attribute Key | Attribute Value  |
| :-------------- | :------------ | :--------------- |
| message         | action        | `/bitsong.merkledrop.v1beta1.MsgClaim` |
| bitsong.merkledrop.v1beta1.EventClaim | merkledrop_id        | {merkledrop_id}         |
| bitsong.merkledrop.v1beta1.EventClaim | index        | {index}         |
| bitsong.merkledrop.v1beta1.EventClaim | coin        | {coin}         |

## EventWithdraw

| Type                     | Attribute Key | Attribute Value   |
| :----------------------- | :------------ | :---------------- |
| bitsong.fantoken.v1beta1.EventWithdraw | merkledrop_id        | {merkledrop_id}         |
| bitsong.fantoken.v1beta1.EventWithdraw | coin        | {coin}         |