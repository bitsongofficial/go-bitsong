# State

## Collection

A `Collection` is a collection of nfts on certain criteria. It stores `id`, `symbol`, `name`, `uri`, `is_mutable`, `update_authority` fields.

```protobuf
message Collection {
  uint64 id = 1;
  /// The symbol for the asset
  string symbol = 2;
  string name = 3;
  string uri = 4;
  // Whether or not the data struct is mutable, default is not
  bool is_mutable = 5;
  // who can update metadata (if is_mutable is true)
  string update_authority = 6;
}
```

- Collection: `0x04 | format(id) -> Collection`
- LastCollectionId `0x06 -> id`

## Metadata

A `Metadata` is a metadata that is attached to an nft.

```protobuf
message MasterEdition {
  uint64 supply = 1;
  uint64 max_supply = 2;
}

message Metadata {
  uint64 id = 1;
  uint64 coll_id = 2;
  // The name of the asset
  string name = 3;
  // URI pointing to JSON representing the asset
  string uri = 4;
  // Royalty basis points that goes to creators in secondary sales (0-10000)
  uint32 seller_fee_basis_points = 5;
  // Immutable, once flipped, all sales of this metadata are considered
  // secondary.
  bool primary_sale_happened = 6;
  // Whether or not the data struct is mutable, default is not
  bool is_mutable = 7;
  // Array of creators, optional
  repeated Creator creators = 8 [ (gogoproto.nullable) = false ];
  // who can update metadata (if is_mutable is true)
  string metadata_authority = 9;
  // who can mint the editions
  string mint_authority = 10;
  MasterEdition master_edition = 11;
}

message Creator {
  string address = 1;
  bool verified = 2;
  // In percentages, NOT basis points ;) Watch out!
  uint32 share = 3;
}
```

- Metadata: `0x03 | format(coll_id) | format(id) -> Metadata`

### Edition

Metadata has `MasterEdition` object integrated for print ability.
It involves `supply` and `max_supply` fields.
When new print is created, supply is increased and new `NFT` object with unique `edition` is created.
Print cannot exceed `max_supply`.

### LastMetadataId

`LastMetadataIdInfo` is the record to track last metadata id per collection to avoid duplication in metadata ids.

```protobuf
message LastMetadataIdInfo {
  uint64 coll_id = 1;
  uint64 last_metadata_id = 2;
}
```

- LastMetadataId `0x05 | format(coll_id) -> LastMetadataIdInfo`

## NFT

A `NFT` is a single unit of a non-fungible token. It has `coll_id`, `metadata_id`, `seq` and `owner`.

The string identifier of nft is expressed as `{coll_id}:{metadata_id}:{seq}`

```protobuf
message NFT {
  uint64 coll_id = 1;
  uint64 metadata_id = 2;
  uint64 seq = 3; // edition nr (0 mean normal nft)
  string owner = 4;
}
```

- NFT: `0x01 | format(coll_id) | format(metadata_id) | format(seq) -> NFT`
- NFT by Owner: `0x02 | owner | format(coll_id) | format(metadata_id) | format(seq) -> nft_identifier`

## Params

Params is a module-wide configuration structure that stores nft module's system parameters.

- Params: `Paramsspace("nft") -> Params`

```protobuf
// Params defines nft module's parameters
message Params {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  cosmos.base.v1beta1.Coin issue_price = 1 [
    (gogoproto.moretags) = "yaml:\"issue_price\"",
    (gogoproto.nullable) = false
  ];
}
```
