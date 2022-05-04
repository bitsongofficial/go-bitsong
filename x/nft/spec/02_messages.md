# Messages

## MsgCreateNFT

`MsgCreateNFT` is a message to be used to create an nft with specific metadata.
At the time of message execution, it creates new metadata object, and after that it creates an nft with the metadata.
After execution, it returns nft and metadata ids.
`sender` is set as nft creator and nft owner after the successful execution.

```protobuf
message MsgCreateNFT {
  string sender = 1;
  bitsong.nft.v1beta1.Metadata metadata = 2 [ (gogoproto.nullable) = false ];
}
message MsgCreateNFTResponse {
  uint64 id = 1;
  uint64 metadata_id = 2;
}
```

Steps:

1. Get unique id by using last metadata id
2. Store newly generated metadata id as last metadata id
3. Create metadata object with newly generated metadataId
4. Set all the verified field as false at the time of creation
5. Store metadata object on storage
6. Pay nft issue fee if it is set to positive value on params store
7. Get unique id by using last nft id
8. Store newly generated nft id as last `nftId`
9. Create nft with owner `msg.Sender`, nftId and metadataId
10. Emit event for nft creation
11. Return `nftId` and `metadataId`

## MsgTransferNFT

`MsgCreateNFT` is a message to update the owner of nft with new one.

```protobuf
message MsgTransferNFT {
  string sender = 1;
  uint64 id = 2;
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
    uint64 metadata_id = 2;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Check if `msg.Sender` is one of the creators of metadata and return permission issue if not
3. Set `Verified` field to true for `metadataa.Data.Creators`
4. Store updated metadata on the storage
5. Emit event for metadata sign

## MsgUpdateMetadata

`MsgUpdateMetadata` is a message to update metadata by the metadata update authority.
`Data` and `PrimarySaleHappened` field can be changed when the metadata has `IsMutable` flag as true.

```protobuf
message MsgUpdateMetadata {
  string sender = 1;
  uint64 metadata_id = 2;
  bitsong.nft.v1beta1.Data data = 3;
  // Immutable, once flipped, all sales of this metadata are considered
  // secondary.
  bool primary_sale_happened = 4;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Check metadata is mutable and if not return immutable error
3. Check `msg.Sender` is authority that has permission to update metadata
4. Update metadata with passed `PrimarySaleHappened` and `Data`
5. Reset Verified field to `false` for Creators
6. Store updated metadata on the storage
7. Emit event for metadata update

## MsgUpdateMetadataAuthority

`MsgUpdateMetadataAuthority` is a message to update metadata authority to another address.

```protobuf
message MsgUpdateMetadataAuthority {
  string sender = 1;
  uint64 metadata_id = 2;
  string new_authority = 3;
}
```

Steps:

1. Check metadata exists with id `msg.MetadataId`
2. Ensure msg.Sender is the authority of metadata
3. Update metadata `UpdateAuthority` with `NewAuthority`
4. Store update metadata on the storage
5. Emit evnet for authority update for the metadata

## MsgCreateCollection

`MsgCreateCollection` is a message to create a new collection.

```protobuf
message MsgCreateCollection {
  string sender = 1;
  string name = 2;
  string uri = 3;
  string update_authority = 4;
}
message MsgCreateCollectionResponse {
    uint64 id = 1;
}
```

Steps:

1. Get unique id by using last collection id
2. Store newly generated collection id as last collection id
3. Create collection with `collectionId`, `msg.Name`, `msg.Uri`, `msg.UpdateAuthority`
4. Store collection on the storage
5. Emit event for new collection creation
6. Return `collectionId` as part of msg response.

## MsgVerifyCollection

`MsgCreateCollection` is a message to verify that an nft is part of collection.
It should be executed by collection authority.

```protobuf
message MsgVerifyCollection {
  string sender = 1;
  uint64 collection_id = 2;
  uint64 nft_id = 3;
}
```

Steps:

1. Check collection exists with id `msg.CollectionId`
2. Ensure collection's `UpdateAuthority` is equal to `msg.Sender`
3. Ensure `msg.CollectionId` and `msg.NftId` are valid ids
4. Set connection between collection and nft
5. Emit event for collection verification

## MsgUnverifyCollection

`MsgUnverifyCollection` is a message to unverify that an nft is part of collection.
It should be executed by collection authority.

```protobuf
message MsgUnverifyCollection {
  string sender = 1;
  uint64 collection_id = 2;
  uint64 nft_id = 3;
}
```

Steps:

1. Check collection exists with id `msg.CollectionId`
2. Ensure collection's `UpdateAuthority` is equal to `msg.Sender`
3. Ensure `msg.CollectionId` and `msg.NftId` are valid ids
4. Remove connection between collection and nft
5. Emit event for collection unverification

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
