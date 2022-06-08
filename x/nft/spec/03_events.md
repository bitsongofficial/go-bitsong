# Events

The nft module emits the following events:

## MsgCreateNFT

| Type                                      | Attribute Key | Attribute Value           |
| :---------------------------------------- | :------------ | :------------------------ |
| bitsong.nft.v1beta1.EventMetadataCreation | creator       | {creator}                 |
| bitsong.nft.v1beta1.EventMetadataCreation | metadata_id   | {metadataId}              |
| bitsong.nft.v1beta1.EventNFTCreation      | creator       | {creator}                 |
| bitsong.nft.v1beta1.EventNFTCreation      | nft_id        | {nftId}                   |
| burn                                      | burner        | {burner}                  |
| burn                                      | amount        | {amount}                  |
| coin_received                             | receiver      | {receiver}                |
| coin_received                             | amount        | {amount}                  |
| coin_spent[]                              | spender       | {spender}                 |
| coin_spent[]                              | amount        | {amount}                  |
| message                                   | action        | /bitsong.nft.MsgCreateNFT |
| message                                   | sender        | {sender}                  |
| transfer                                  | recipient     | {recipient}               |
| transfer                                  | sender        | {sender}                  |
| transfer                                  | amount        | {amount}                  |

## MsgPrintEdition

| Type                                  | Attribute Key | Attribute Value              |
| :------------------------------------ | :------------ | :--------------------------- |
| bitsong.nft.v1beta1.EventPrintEdition | printer       | {printer}                    |
| bitsong.nft.v1beta1.EventPrintEdition | metadata_id   | {metadataId}                 |
| burn                                  | burner        | {burner}                     |
| burn                                  | amount        | {amount}                     |
| coin_received                         | receiver      | {receiver}                   |
| coin_received                         | amount        | {amount}                     |
| coin_spent[]                          | spender       | {spender}                    |
| coin_spent[]                          | amount        | {amount}                     |
| message                               | action        | /bitsong.nft.MsgPrintEdition |
| message                               | sender        | {sender}                     |
| transfer                              | recipient     | {recipient}                  |
| transfer                              | sender        | {sender}                     |
| transfer                              | amount        | {amount}                     |

## MsgTransferNFT

| Type                                 | Attribute Key | Attribute Value                       |
| :----------------------------------- | :------------ | :------------------------------------ |
| bitsong.nft.v1beta1.EventNFTTransfer | nft_id        | {nft_id}                              |
| bitsong.nft.v1beta1.EventNFTTransfer | sender        | {sender}                              |
| bitsong.nft.v1beta1.EventNFTTransfer | receiver      | {receiver}                            |
| message                              | action        | /bitsong.nft.v1beta1.EventNFTTransfer |

## MsgSignMetadata

| Type                                  | Attribute Key | Attribute Value              |
| :------------------------------------ | :------------ | :--------------------------- |
| bitsong.nft.v1beta1.EventMetadataSign | signer        | {signer}                     |
| bitsong.nft.v1beta1.EventMetadataSign | metadata_id   | {metadata_id}                |
| message                               | action        | /bitsong.nft.MsgSignMetadata |

## MsgUpdateMetadata

| Type                                    | Attribute Key | Attribute Value                |
| :-------------------------------------- | :------------ | :----------------------------- |
| bitsong.nft.v1beta1.EventMetadataUpdate | updater       | {updater}                      |
| bitsong.nft.v1beta1.EventMetadataUpdate | metadata_id   | {metadata_id}                  |
| message                                 | action        | /bitsong.nft.MsgUpdateMetadata |

## MsgUpdateMetadataAuthority

| Type                                             | Attribute Key | Attribute Value                         |
| :----------------------------------------------- | :------------ | :-------------------------------------- |
| bitsong.nft.v1beta1.EventMetadataAuthorityUpdate | metadata_id   | {metadata_id}                           |
| bitsong.nft.v1beta1.EventMetadataAuthorityUpdate | new_authority | {new_authority}                         |
| message                                          | action        | /bitsong.nft.MsgUpdateMetadataAuthority |

## MsgCreateCollection

| Type                                        | Attribute Key | Attribute Value                  |
| :------------------------------------------ | :------------ | :------------------------------- |
| bitsong.nft.v1beta1.EventCollectionCreation | creator       | {creator}                        |
| bitsong.nft.v1beta1.EventCollectionCreation | collection_id | {collection_id}                  |
| message                                     | action        | /bitsong.nft.MsgCreateCollection |

## MsgVerifyCollection

| Type                                            | Attribute Key | Attribute Value                  |
| :---------------------------------------------- | :------------ | :------------------------------- |
| bitsong.nft.v1beta1.EventCollectionVerification | verifier      | {verifier}                       |
| bitsong.nft.v1beta1.EventCollectionVerification | collection_id | {collection_id}                  |
| bitsong.nft.v1beta1.EventCollectionVerification | nft_id        | {nft_id}                         |
| message                                         | action        | /bitsong.nft.MsgVerifyCollection |

## MsgUnverifyCollection

| Type                                              | Attribute Key | Attribute Value                    |
| :------------------------------------------------ | :------------ | :--------------------------------- |
| bitsong.nft.v1beta1.EventCollectionUnverification | verifier      | {verifier}                         |
| bitsong.nft.v1beta1.EventCollectionUnverification | collection_id | {collection_id}                    |
| bitsong.nft.v1beta1.EventCollectionUnverification | nft_id        | {nft_id}                           |
| message                                           | action        | /bitsong.nft.MsgUnverifyCollection |

## MsgUpdateCollectionAuthority

| Type                                               | Attribute Key | Attribute Value                           |
| :------------------------------------------------- | :------------ | :---------------------------------------- |
| bitsong.nft.v1beta1.EventUpdateCollectionAuthority | collection_id | {collection_id}                           |
| bitsong.nft.v1beta1.EventUpdateCollectionAuthority | new_authority | {new_authority}                           |
| message                                            | action        | /bitsong.nft.MsgUpdateCollectionAuthority |
