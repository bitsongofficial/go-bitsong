# Messages

## MsgCreateNFT

`MsgCreateNFT` is a message to be used to create an nft with specific metadata.
At the time of message execution, it creates new metadata object, and after that it creates an nft with the metadata.
After execution, it returns nft and metadata ids.
`sender` is set as nft creator and nft owner after the successful execution.

```protobuf
message MsgCreateNFT {
  string sender = 1;
  bitsong.nft.v1beta1.Metadata metadata = 2  [ (gogoproto.nullable) = false ];
}
message MsgCreateNFTResponse {
  string id = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
}
```

Steps:

1. Ensure collection with id value `msg.CollId` exists
2. Ensure `msg.Sender` is owner of the collection
3. Get unique metadata id by using last metadata id
4. Store newly generated metadata id as last metadata id
5. Create metadata object with newly generated metadataId
6. Set all the verified field as false at the time of creation
7. Store metadata object on storage
8. Emit event for metadata creation
9. Pay nft issue fee if it is set to positive value on params store
10. Create nft with owner `msg.Sender`, `msg.CollId`, `metadataId` and 0 as `Seq`.
11. Emit event for nft creation
12. Return nft id and `metadataId`

## MsgPrintEdition

`MsgPrintEdition` is a message to print new edition for a metadata.
Editions can be printed for metadata with `MasterEdition` field by metadata owner.

```protobuf
message MsgPrintEdition {
  string sender = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
  string owner = 4;
}
message MsgPrintEditionResponse {
  string id = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Ensure that master edition nft (0 as `Seq`) exists
3. Ensure that message is executed by `metadata.MintAuthority`
4. Ensure that meatadata is master edition metadata
5. Ensure total supply is of MasterEdition is lower than max supply
6. Generate new edition number for metadata
7. Pay nft issue fee if it is set to positive value on params store
8. Create a new NFT with `msg.CollId`, `msg.MetadataId` with new edition number and store
9. Store the updated edition number on the storage
10. Emit `EventPrintEdition` event for print edition
11. Return nft identifier

## MsgTransferNFT

`MsgTransferNFT` is a message to update the owner of nft with new one.

```protobuf
message MsgTransferNFT {
  string sender = 1;
  string id = 2;
  string new_owner = 3;
}
```

Steps:

1. Check nft exists with id `msg.Id`
2. Check nft owner is equal to `msg.Sender`
3. Update the owner of NFT to `msg.NewOwner`
4. Store the updated nft on the storage
5. Emit event for nft transfer

## MsgSignMetadata

`MsgSignMetadata` is a message to sign the creator field of metadata.
Once it is executed, it set the `Verified` field of creators on metadata as true.

```protobuf
message MsgSignMetadata {
  string sender = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Check if `msg.Sender` is one of the creators of metadata and return permission issue if not
3. Set `Verified` field to true for `metadata.Creators`
4. Store updated metadata on the storage
5. Emit event for metadata sign

## MsgUpdateMetadata

`MsgUpdateMetadata` is a message to update metadata by the metadata update authority.
`Name`, `URI`, `SellerFeeBasisPoints` and `Creators` fields can be changed when the metadata has `IsMutable` flag as true.

```protobuf
message MsgUpdateMetadata {
  string sender = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
  // The name of the asset
  string name = 4;
  // URI pointing to JSON representing the asset
  string uri = 5;
  // Royalty basis points that goes to creators in secondary sales (0-10000)
  uint32 seller_fee_basis_points = 6;
  // Array of creators, optional
  repeated bitsong.nft.v1beta1.Creator creators = 7
      [ (gogoproto.nullable) = false ];
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Check metadata is mutable and if not return immutable error
3. Check `msg.Sender` is authority that has permission to update metadata
4. Update metadata with passed `Name`, `Uri`, `SellerFeeBasisPoints` and `Creators`
5. Reset Verified field to `false` for Creators
6. Store updated metadata on the storage
7. Emit event for metadata update

## MsgUpdateMetadataAuthority

`MsgUpdateMetadataAuthority` is a message to update metadata authority to another address.

```protobuf
message MsgUpdateMetadataAuthority {
  string sender = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
  string new_authority = 4;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Ensure msg.Sender is the authority of metadata
3. Update metadata `MetadataAuthority` with `NewAuthority`
4. Store update metadata on the storage
5. Emit evnet for authority update for the metadata

## MsgUpdateMintAuthority

`MsgUpdateMetadataAuthority` is a message to update mint authority to another address.

```protobuf
message MsgUpdateMintAuthority {
  string sender = 1;
  uint64 coll_id = 2;
  uint64 metadata_id = 3;
  string new_authority = 4;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Ensure msg.Sender is the mint authority of metadata
3. Update metadata `MintAuthority` with `NewAuthority`
4. Store update metadata on the storage
5. Emit evnet for authority update for the metadata

## MsgCreateCollection

`MsgCreateCollection` is a message to create a new collection.

```protobuf
message MsgCreateCollection {
  string sender = 1;
  string symbol = 2;
  string name = 3;
  string uri = 4;
  bool is_mutable = 5;
  string update_authority = 6;
}
message MsgCreateCollectionResponse {
  uint64 id = 1;
}
```

Steps:

1. Get unique id by using last collection id
2. Store newly generated collection id as last collection id
3. Create collection with `collectionId`, `msg.Symbol`, `msg.Name`, `msg.Uri`, `msg.IsMutable` and `msg.UpdateAuthority`
4. Store collection on the storage
5. Emit event for new collection creation
6. Return `collectionId` as part of msg response.

## MsgUpdateCollectionAuthority

`MsgUpdateCollectionAuthority` is a message to update collection authority to a new one.
It should be executed by collection authority.

```protobuf
message MsgUpdateCollectionAuthority {
  string sender = 1;
  uint64 collection_id = 2;
  string new_authority = 3;
}
```

Steps:

1. Check collection exists with id `msg.CollectionId`
2. Ensure collection's `UpdateAuthority` is equal to `msg.Sender`
3. Update collection authority with `msg.NewAuthority`
4. Store updated collection object into storage
5. Emit event for collection authority update
