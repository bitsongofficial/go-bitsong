# Messages

## MsgCreateAuction

Store auction record on-chain from provided Msg fields.

## MsgSetAuctionAuthority

1. Check Msg sender is auction authority
2. Ensure new authority is an accurate address
3. Update auction authority with new authority

## MsgStartAuction

1. Check sender is auction authority
2. Calculate auction end time from current time and auction duration
3. Set the state of auction as started

## MsgEndAuction

1.  Check executor is a correct authority of the auction
2.  Check auction is not already ended
3.  Set auction end time
4.  Set auction status as ended

## MsgPlaceBid

1. Verify bid is valid for the auction
2. Verify auction has not ended
3. Verify auction has started
4. Load bidder metadata or create one
5. Add new bid for the auction
6. Confirm payer does have enough token to pay the bid
7. Transfer amount of token to bid account
8. Serialize new auction state with new bid - update auction end time based on gap_tick_size if required
9. Update bidder metadata

## MsgCancelBid

1. Load the auction and verify this bid is valid.
2. Check instant_sale_price and update cancelled bids if auction still active
3. Check metadata exists for the bidder
4. Transfer tokens back to the bidder
5. Refuse to cancel if the auction ended and this person is a winning account.
6. Refuse to cancel if bidder set price above or equal instant_sale_price
7. Update bidder Metadata
8. Update auction with remaining bids

## MsgClaimBid

1. Load the auction and verify this bid is valid.
2. Ensure bidder is owner
3. Ensure user paid instant sale price or auction is ended
4. Send bid amount to auction creator (set primary_sale_happened as true, if it was previously true, process royalties)
5. Transfer ownership of NFT to bidder
