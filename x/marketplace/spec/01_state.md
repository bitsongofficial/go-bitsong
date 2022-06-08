# State

## Auction

```protobuf
message Auction {
    // unique identifier of the auction
    uint64 id = 1;
    // authority with permission to modify this auction.
    string authority = 2;
    // NFT being used to bid
    uint64 nft_id = 3;
    // Describes transfering nft ownership only or metadata ownership as well
    AuctionPrizeType prize_type = 4;
    // Duration of the auction
    google.protobuf.Duration duration = 5 [ (gogoproto.stdduration) = true ];
    // Denom to be used on bids
    string bid_denom = 6;
    // Minimum price for any bid to meet.
    uint64 price_floor = 7;
    // Instant sale price
    uint64 instant_sale_price = 8;
    // Tick size - how much higher the next bid must be to beat out the previous bid.
    uint64 tick_size = 9;
    // The state the auction is in, whether it has started or ended.
    AuctionState state = 10;
    // The amount of bid put last time
    uint64 last_bid_amount = 11;
    // The time the last bid was placed, used to keep track of auction timing.
    google.protobuf.Timestamp last_bid = 12 [ (gogoproto.stdtime) = true ];
    // Slot time the auction was officially ended by.
    google.protobuf.Timestamp ended_at = 13 [ (gogoproto.stdtime) = true ];
    // End time is the cut-off point that the auction is forced to end by.
    google.protobuf.Timestamp end_auction_at = 14 [ (gogoproto.stdtime) = true ];
    // Ticked to true when a prize is claimed by person who won it
    bool claimed = 15;
    // Only valid for LimitedEditionPrints auction
    uint64 printable_editions = 16;
}

/// Define valid auction state transitions.
enum AuctionState {
  EMPTY = 0 [ (gogoproto.enumvalue_customname) = "Empty" ];
  CREATED = 1 [ (gogoproto.enumvalue_customname) = "Created" ];
  STARTED = 2 [ (gogoproto.enumvalue_customname) = "Started" ];
  ENDED = 3 [ (gogoproto.enumvalue_customname) = "Ended" ];
}

```

### Auction type

```protobuf
enum AuctionPrizeType {
    // Transfer ownership of only nft without metadata
    NFT_ONLY_TRANSFER = 0 [ (gogoproto.enumvalue_customname) = "NftOnlyTransfer" ];
    // Transfer ownership of both nft and metadata
    FULL_RIGHTS_TRANSFER = 1 [ (gogoproto.enumvalue_customname) = "FullRightsTransfer" ];
    // Printing a new child edition from limited supply
    LIMITED_EDITION_PRINTS = 2 [ (gogoproto.enumvalue_customname) = "LimitedEditionPrints" ];
    // Printing a new child edition from unlimited supply
    OPEN_EDITION_PRINTS = 3 [ (gogoproto.enumvalue_customname) = "OpenEditionPrints" ];
}
```

1. `NftOnlyTransfer` is the auction for sending nft to winner bidder.
2. `FullRightsTransfer` is the auction to transfer both nft and metadata ownership.
3. `LimitedEditionPrints` is the auction to provide limited number of editions printed to auction winners. Editions can be printed after auction ends.
4. `OpenEditionPrints` is the auction to provide all auction participants to get printed versions. Editions can be instantly printed after bid even before auction ends.

Note: `Metadata` ownership is temporarily transfered to marketplace module during the auction phase for `FullRightsTransfer`, `LimitedEditionPrints` and `OpenEditionPrints`.
After auction ends, metadata ownership is returned back to owner for `LimitedEditionPrints` and `OpenEditionPrints` case.

### Storage

- Auction: `0x01 | format(id) -> Auction`
- Auction by Authority: `0x02 | owner | format(id) -> auction_id`
- Auction by EndTime: `0x03 | format(endTime) | format(id) -> auction_id`
- LastAuctionId `0x04 -> id`

## Bid

```protobuf
// Bids associate a bidding key with an amount bid.
message Bid {
  string bidder = 1;
  uint64 auction_id = 2;
  uint64 amount = 3;
}
```

- Bid: `0x05 | format(auction_id) | bidder -> Bid`
- Bid by bidder: `0x06 | bidder | format(auction_id) -> Bid`

## BidderMetadata

```protobuf
/// Models a set of metadata for a bidder, meant to be stored in a PDA. This allows looking up
/// information about a bidder regardless of if they have won, lost or cancelled.
message BidderMetadata {
    // Relationship with the bidder who's metadata this covers.
    string bidder = 1;
    // Relationship with the auction this bid was placed on.
    string last_auction_id = 2;
    // Amount that the user bid.
    uint64 last_bid = 3;
    // Tracks the last time this user bid.
    google.protobuf.Timestamp last_bid_timestamp = 4 [ (gogoproto.stdtime) = true ];
    // Whether the last bid the user made was cancelled. This should also be enough to know if the
    // user is a winner, as if cancelled it implies previous bids were also cancelled.
    bool last_bid_cancelled = 5;
}
```

- BidderMetadata: `0x07 | bidder -> BidderMetadata`
