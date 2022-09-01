# Events

## Messages

### MsgCreateLaunchPad

| Type                                           | Attribute Key | Attribute Value                         |
| :--------------------------------------------- | :------------ | :-------------------------------------- |
| bitsong.launchpad.v1beta1.EventCreateLaunchPad | creator       | {creator}                               |
| bitsong.launchpad.v1beta1.EventCreateLaunchPad | collection_id | {collection_id}                         |
| message                                        | action        | /bitsong.marketplace.MsgCreateLaunchPad |

### MsgUpdateLaunchPad

| Type                                           | Attribute Key | Attribute Value                         |
| :--------------------------------------------- | :------------ | :-------------------------------------- |
| bitsong.launchpad.v1beta1.EventUpdateLaunchPad | creator       | {creator}                               |
| bitsong.launchpad.v1beta1.EventUpdateLaunchPad | collection_id | {collection_id}                         |
| message                                        | action        | /bitsong.marketplace.MsgUpdateLaunchPad |

## MsgCloseLaunchPad

| Type                                          | Attribute Key | Attribute Value                        |
| :-------------------------------------------- | :------------ | :------------------------------------- |
| bitsong.launchpad.v1beta1.EventCloseLaunchPad | creator       | {creator}                              |
| bitsong.launchpad.v1beta1.EventCloseLaunchPad | collection_id | {collection_id}                        |
| message                                       | action        | /bitsong.marketplace.MsgCloseLaunchPad |

## MsgMintNFT

| Type                                            | Attribute Key | Attribute Value                 |
| :---------------------------------------------- | :------------ | :------------------------------ |
| bitsong.launchpad.v1beta1.EventMintNFT          | collection_id | {collection_id}                 |
| bitsong.launchpad.v1beta1.EventMintNFT          | nft_id        | {nft_id}                        |
| bitsong.launchpad.v1beta1.EventMetadataCreation | metadata_id   | {metadata_id}                   |
| bitsong.launchpad.v1beta1.EventMetadataCreation | creator       | {creator}                       |
| bitsong.launchpad.v1beta1.EventNFTCreation      | creator       | {creator}                       |
| bitsong.launchpad.v1beta1.EventNFTCreation      | nft_id        | {nft_id}                        |
| message                                         | action        | /bitsong.marketplace.MsgMintNFT |

## MsgMintNFTs

| Type                                            | Attribute Key | Attribute Value                  |
| :---------------------------------------------- | :------------ | :------------------------------- |
| bitsong.launchpad.v1beta1.EventMintNFT          | collection_id | {collection_id}                  |
| bitsong.launchpad.v1beta1.EventMintNFT          | nft_id        | {nft_id}                         |
| bitsong.launchpad.v1beta1.EventMetadataCreation | metadata_id   | {metadata_id}                    |
| bitsong.launchpad.v1beta1.EventMetadataCreation | creator       | {creator}                        |
| bitsong.launchpad.v1beta1.EventNFTCreation      | creator       | {creator}                        |
| bitsong.launchpad.v1beta1.EventNFTCreation      | nft_id        | {nft_id}                         |
| message                                         | action        | /bitsong.marketplace.MsgMintNFTs |
