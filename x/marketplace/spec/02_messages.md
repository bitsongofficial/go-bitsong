# Messages

## MsgCreateAuction

`MsgCreateAuction` is a message to be used to create an auction with provided params.
The NFT should not be already being used by another auction and it should be owned by the sender.
After the execution, nft is owned by marketplace module - not by an individual account.

```protobuf
message MsgCreateAuction {
  string sender = 1;
  /// NFT being used to bid
  uint64 nft_id = 2;
  // Describes transfering nft ownership only or metadata ownership as well
  bitsong.marketplace.v1beta1.AuctionPrizeType prize_type = 3;
  // Denom to be used on bids
  string bid_denom = 4;
  // Duration of the auction
  google.protobuf.Duration duration = 5 [ (gogoproto.stdduration) = true ];
  // Minimum price for any bid to meet.
  uint64 price_floor = 6;
  // Instant sale price
  uint64 instant_sale_price = 7;
  // Tick size - how much higher the next bid must be to beat out the previous bid.
  uint64 tick_size = 8;
}
message MsgCreateAuctionResponse {
  uint64 id = 1;
}
```

Steps:

1. Ensure nft is owned by the sender
2. Ensure nft metadata is owned by the sender if auction prize type is `FullRightsTransfer`
3. Send nft ownership to marketplace module
4. If auction is for transferring metadata ownership as well, metadata authority is transferred to marketplace module
5. Create auction object from provided params
6. Emit event for auction creation
7. Return auction id

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
2.  Check auction is not already ended
3.  Set auction end time
4.  Set auction status as ended
5.  Check auction has winning bid
6.  If winning bid does not exists, send nft and metadata ownership to auction authority
7.  Emit event for auction end

Notes: Metadata ownership should be returned back to owner regardless winner bidder exists or not when it's `LimitedEditionPrints` or `OpenEditionPrints`

## MsgPlaceBid

`MsgPlaceBid` is a message to place bid on an auction.

```protobuf
message MsgPlaceBid {
  string sender = 1;
  uint64 auction_id = 2;
  cosmos.base.v1beta1.Coin amount = 3;
}
```

Steps:

1. Verify auction is `Started` status
2. Verify bid is valid for the auction (check `bid_denom`, `tick_size` and `last_bid_amount`)
3. Check if previous bid exists for this auction by the bidder and if exists reject
4. Add new bid for the auction on the storage
5. Confirm payer does have enough token to pay the bid
6. Transfer amount of token to bid account
7. Serialize new auction state with new bid
8. Update or create bidder metadata
9. If the amount exceeds `instant_sale_price`, end the auction
10. Emit event for placing bid

TODO: process should be different per auction type

- `NftOnlyTransfer`: Single winner
- `FullRightsTransfer`: Single winner
- `LimitedEditionPrints`: Multiple winners
- `OpenEditionPrints`: Multiple winners

Notes: `LimitedEditionPrints` will be checking if the bid will be on top `n` - number of editions to print - not top bid.
`OpenEditionPrints` will be checking only floor price.

## MsgCancelBid

`MsgPlaceBid` is a message to cancel bid on an auction.

```protobuf
message MsgCancelBid {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1. Load the auction and verify this bid is valid.
2. Refuse to cancel if the auction ended and this person is a winning account.
3. Remove bid from the storage
4. Transfer tokens back to the bidder
5. Update bidder Metadata
6. Emit event for cancelling bid

Note: Possibility to add cancel fee for the bid.

TODO: process should be different per auction type

- `NftOnlyTransfer`: Single winner
- `FullRightsTransfer`: Single winner
- `LimitedEditionPrints`: Multiple winners
- `OpenEditionPrints`: Multiple winners

Notes: `LimitedEditionPrints` will be checking if the bid is not on top `n` - to be able to cancel.
`OpenEditionPrints` won't be able to cancel.

## MsgClaimBid

`MsgClaimBid` is a message to claim auction prize.

```protobuf
message MsgClaimBid {
  string sender = 1;
  uint64 auction_id = 2;
}
```

Steps:

1. Load the auction and verify this bid is valid.
2. Ensure the sender is winner bidder
3. Send bid amount to auction authority
4. If `primary_sale_happened` is true, process royalties from NFT's `seller_fee_basis_points` field to creators
5. Set `primary_sale_happened` as true if it was not set already
6. Transfer ownership of NFT to bidder
7. If auction type is for transferring metadata ownership as well, transfer metadata ownership as well
8. Update auction with claimed status
9. Emit event for claiming bid

TODO: process should be different per auction type

- `NftOnlyTransfer`: Single winner
- `FullRightsTransfer`: Single winner
- `LimitedEditionPrints`: Multiple winners
- `OpenEditionPrints`: Multiple winners

Notes: `OpenEditionPrints` will be able to claim even auction not ends.
`LimitedEditionPrints` will be able to claim only after auction ends.
