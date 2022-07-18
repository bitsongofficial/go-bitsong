# Messages

## MsgCreateAuction

`MsgCreateAuction` is a message to be used to create an auction with provided params.
The NFT should not be already being used by another auction and it should be owned by the sender.
After the execution, nft is owned by marketplace module - not by an individual account.

```protobuf
message MsgCreateAuction {
  string sender = 1;
  // NFT being used to bid
  string nft_id = 2;
  // Describes transfering nft ownership only or metadata ownership as well
  bitsong.marketplace.v1beta1.AuctionPrizeType prize_type = 3;
  // Denom to be used on bids
  string bid_denom = 4;
  // Duration of the auction
  google.protobuf.Duration duration = 5
      [ (gogoproto.stdduration) = true, (gogoproto.nullable) = false ];
  // Minimum price for any bid to meet.
  uint64 price_floor = 6;
  // Instant sale price
  uint64 instant_sale_price = 7;
  // Tick size - how much higher the next bid must be to beat out the previous bid.
  uint64 tick_size = 8;
  // Edition limitation for limited edition auction
  uint64 edition_limit = 9;
}
message MsgCreateAuctionResponse {
  uint64 id = 1;
}
```

Steps:

1. Pay auction creation fee if fee exists
2. Ensure nft is owned by the sender
3. Ensure nft metadata is owned by the sender if auction prize type is `FullRightsTransfer`
4. Send ownerships (nft, metadata) to marketplace module based on auction type
   - If full rights transfer auction, transfer metadata, mint, nft ownerships to marketplace module
   - If mint authority transfer auction, transfer mint authority to marketplace module
   - If metadata authority transfer auction, transfer meatadata authority to marketplace module
   - If nft only transfer auction, transfer nft ownership to marketplace module
   - If print editions auction, transfer mint authority to marketplace module
5. Generate new auction id and store last auction id
6. Create auction object from provided params and generated auction id
7. Emit event for auction creation
8. Return auction id

Notes: Send only metadata ownership to auction when it's `LimitedEditionPrints` and `OpenEditionPrints`

## MsgSetAuctionAuthority

`MsgSetAuctionAuthority` is a message to send authority of the auction to a new address.
This should be executed by auction authority.

```protobuf
message MsgSetAuctionAuthority {
  string sender = 1;
  uint64 auction_id = 2;
  string new_authority = 3;
}
```

Steps:

1. Check Msg sender is auction authority
2. Ensure new authority is an accurate address
3. Update auction authority with new authority
4. Emit event for authority update

## MsgStartAuction

`MsgStartAuction` is a message to start the auction that is created via `MsgCreateAuction`.

```protobuf
message MsgStartAuction {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1. Check sender is auction authority
2. Ensure auction status is `Created`
3. Calculate auction end time from current time and auction duration
4. Set the state of auction to `Started`
5. Store updated auction into store
6. Emit event for auction start

## MsgEndAuction

`MsgEndAuction` is a message to end the auction by auction authority before auction ends.

```protobuf
message MsgEndAuction {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1.  Check executor is a correct authority of the auction
2.  Ensure auction is not already ended
3.  Set auction end time
4.  Set auction status as ended
5.  Check auction has winning bid
6.  Process nft and metadata ownerships based on auction type
7.  Autoclaim of bids for print auctions
8.  Emit event for auction end

Notes: Metadata mint authority is returned back to owner regardless winner bidder exists or not when it's `LimitedEditionPrints` or `OpenEditionPrints`

## MsgPlaceBid

`MsgPlaceBid` is a message to place bid on an auction.

```protobuf
message MsgPlaceBid {
  string sender = 1;
  uint64 auction_id = 2;
  cosmos.base.v1beta1.Coin amount = 3 [ (gogoproto.nullable) = false ];
}
```

Steps:

1. Ensure auction is `Started` status
2. Verify bid is valid for the auction (check `bid_denom`, `tick_size`, previous bids)
3. Check if previous bid exists for this auction by the bidder and if exists reject
4. Add new bid for the auction on the storage
5. Transfer amount of token to bid account
6. Serialize new auction state with new bid as last bid
7. Update or create bidder metadata
8. If the amount exceeds `instant_sale_price`, end the auction for non-print auctions (`NftOnlyTransfer`, `FullRightsTransfer`, `MintAuthorityTransfer`, `MetadataAuthorityTransfer`)
9. Emit event for placing bid

## MsgCancelBid

`MsgCancelBid` is a message to cancel bid on an auction.

```protobuf
message MsgCancelBid {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1. Load the auction and verify this bid is valid.
2. Refuse to cancel the bid is winner bid.
3. Remove bid from the storage
4. Transfer tokens back to the bidder
5. Update bidder Metadata
6. Emit event for cancelling bid

Note: Possibility to add cancel fee for the bid.

## MsgClaimBid

`MsgClaimBid` is a message to claim auction prize.

```protobuf
message MsgClaimBid {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1. Ensure that auction is already ended
2. Load the auction and verify this bid is valid.
3. Ensure the sender is winner bidder
4. Send bid amount to auction authority
5. If `primary_sale_happened` is true for the metadata, process royalties from NFT's `seller_fee_basis_points` field to creators
6. Set `primary_sale_happened` as true if it's not set already
7. Transfer ownerships or print an edition to the bidder based on auction type
   - If full rights transfer auction, transfer metadata, mint, nft ownerships to winner bidder
   - If mint authority transfer auction, transfer mint authority to winner bidder
   - If metadata authority transfer auction, transfer meatadata authority to winner bidder
   - If nft only transfer auction, transfer nft ownership to winner bidder
   - If print editions auction, print edition with mint authority
8. Increase auction claim count by 1
9. Delete claimed bid object
10. Emit event for claiming bid
