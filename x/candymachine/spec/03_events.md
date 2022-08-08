# Events

## Messages

### MsgCreateCandyMachine

| Type                                                 | Attribute Key | Attribute Value                            |
| :--------------------------------------------------- | :------------ | :----------------------------------------- |
| bitsong.candymachine.v1beta1.EventCreateCandyMachine | creator       | {creator}                                  |
| bitsong.candymachine.v1beta1.EventCreateCandyMachine | collection_id | {collection_id}                            |
| message                                              | action        | /bitsong.marketplace.MsgCreateCandyMachine |

### MsgUpdateCandyMachine

| Type                                                 | Attribute Key | Attribute Value                            |
| :--------------------------------------------------- | :------------ | :----------------------------------------- |
| bitsong.candymachine.v1beta1.EventUpdateCandyMachine | creator       | {creator}                                  |
| bitsong.candymachine.v1beta1.EventUpdateCandyMachine | collection_id | {collection_id}                            |
| message                                              | action        | /bitsong.marketplace.MsgUpdateCandyMachine |

## MsgCloseCandyMachine

| Type                                                | Attribute Key | Attribute Value                           |
| :-------------------------------------------------- | :------------ | :---------------------------------------- |
| bitsong.candymachine.v1beta1.EventCloseCandyMachine | creator       | {creator}                                 |
| bitsong.candymachine.v1beta1.EventCloseCandyMachine | collection_id | {collection_id}                           |
| message                                             | action        | /bitsong.marketplace.MsgCloseCandyMachine |

## MsgMintNFTResponse

| Type                                               | Attribute Key | Attribute Value                 |
| :------------------------------------------------- | :------------ | :------------------------------ |
| bitsong.candymachine.v1beta1.EventMintNFT          | collection_id | {collection_id}                 |
| bitsong.candymachine.v1beta1.EventMintNFT          | nft_id        | {nft_id}                        |
| bitsong.candymachine.v1beta1.EventMetadataCreation | metadata_id   | {metadata_id}                   |
| bitsong.candymachine.v1beta1.EventMetadataCreation | creator       | {creator}                       |
| bitsong.candymachine.v1beta1.EventNFTCreation      | creator       | {creator}                       |
| bitsong.candymachine.v1beta1.EventNFTCreation      | nft_id        | {nft_id}                        |
| bitsong.candymachine.v1beta1.EventNFTTransfer      | nft_id        | {nft_id}                        |
| bitsong.candymachine.v1beta1.EventNFTTransfer      | sender        | {sender}                        |
| bitsong.candymachine.v1beta1.EventNFTTransfer      | receiver      | {receiver}                      |
| message                                            | action        | /bitsong.marketplace.MsgMintNFT |
