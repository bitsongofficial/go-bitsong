# State

## LaunchPad

```protobuf
message LaunchPad {
  // id of the collection
  uint64 coll_id = 1;
  // mint price
  uint64 price = 2;
  // wallet to receive payments
  string treasury = 3;
  // denom for the auction
  string denom = 4;
  // Timestamp when minting is allowed
  uint64 go_live_date = 5;
  // mint end timestamp - not considered when set as 0
  uint64 end_timestamp = 6;
  // max mintable amount
  uint64 max_mint = 7;
  // minted amount
  uint64 minted = 8;
  // owner of launchpad
  string authority = 9;
  // all metadata url is generated from metadata_base_url
  string metadata_base_url = 10;
  // mutability of the minted nfts
  bool mutable = 11;
  // Royalty basis points that goes to creators in secondary sales (0-10000)
  uint32 seller_fee_basis_points = 12;
  // Creators of metadata
  repeated bitsong.nft.v1beta1.Creator creators = 13 [ (gogoproto.nullable) = false ];
}
```

| Parameter            | Description                                                    |
| :------------------- | :------------------------------------------------------------- |
| collId               | Id of collection to mint                                       |
| price                | price of nft to pay for minting an nft                         |
| treasury             | wallet to receive payments                                     |
| denom                | denom to be used for payment                                   |
| goLiveDate           | Timestamp when minting is allowed                              |
| endTimestamp         | Describes automatically ending timestamp                       |
| maxMint              | Describes maximum number of nfts that can be minted            |
| metadataBaseUrl      | Base url of metadata for the collection                        |
| mutable              | Mutability of the nft items minted via launchpad               |
| authority            | Authority of the launchpad with update and close permission    |
| minted               | Number of minted nfts via the launchpad.                       |
| sellerFeeBasisPoints | Seller fee basis points of the nft items minted via launchpad. |
| creators             | Creators of the nft. Collection nfts share the same creators.  |

When a collection owner creates launchpad, ownership of collection is sent to launchpad, and it can be returned only after launchpad ends.

The tokens spent on minting nfts, is sent to the treasury account.

Note: How do we handle rare nfts on the collection with different price?

- `endTimestamp` describes launchpad automatically ends when that time comes.
- `maxMint` describes launchpad automatically ends when that number of nfts are minted.

Store:

- Launchpad: `0x01 | format(collection_id) | bidder -> Bid`
- Launchpad by EndTime: `0x02 | format(EndTime) | format(collection_id) -> Bid`

Notes: Launchpad by EndTime queue is only set when `endTimestamp` is not `0`.
