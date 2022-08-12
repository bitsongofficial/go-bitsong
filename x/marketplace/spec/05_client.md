# Client

## Query

### auction

Query all auctions

```
bitsongd query marketplace auctions
```

Query auction by id

```
bitsongd query marketplace auction <auction_id>

bitsongd query marketplace auction 1
```

### bids

Query bids by auction id

```
bitsongd query marketplace bids-by-auction <auction_id>

bitsongd query marketplace bids-by-auction 1
```

Query bids by bidder

```
bitsongd query marketplace bids-by-bidder <bidder>

bitsongd query marketplace bids-by-bidder $(bitsongd keys show -a validator --keyring-backend=test)
```

### bidder metadata

Query bidder metadata by bidder address

```
bitsongd query marketplace bidder-metadata <bidder>

bitsongd query marketplace bidder-metadata $(bitsongd keys show -a validator --keyring-backend=test)
```

## Messages

### CreateAuction

Create auction from nft id, auction type, bid denom, auction duration, price floor, instant sale price and tick size

```
bitsongd tx marketplace create-auction \
--nft-id=<nft_id> \
--prize-type=<auction_type> \
--bid-denom=<auction_bid_denom> \
--duration=<auction_duration> \
--price-floor=<price_floor> \
--instant-sale-price=<instant_sale_price> \
--tick-size=<tick_size> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace create-auction --nft-id=1:1:0 --prize-type="NFT_ONLY_TRANSFER" --bid-denom="ubtsg" --duration="864000s" --price-floor=1000000 --instant-sale-price=100000000 --tick-size=100000 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### SetAuctionAuthority

Update auction authority by auction owner

```
bitsongd tx marketplace set-auction-authority \
--auction-id=<auction_id> \
--new-authority=<new_authority> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace set-auction-authority --auction-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### StartAuction

Start auction by auction authority

```
bitsongd tx marketplace start-auction \
--auction-id=<auction_id> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace start-auction --auction-id=1 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### EndAuction

End auction by auction authority

```
bitsongd tx marketplace end-auction \
--auction-id=<auction_id> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace end-auction --auction-id=1 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### PlaceBid

Place bid on an auction by bidder with amount

```
bitsongd tx marketplace place-bid \
--auction-id=<auction_id> \
--amount=<bid_amount> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace place-bid --auction-id=1 --amount="1000000ubtsg" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### CancelBid

Cancel bid from an auction

```
bitsongd tx marketplace cancel-bid \
--auction-id=<auction_id> \
--from=<bidder> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace cancel-bid --auction-id=1 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### ClaimBid

Claim winner bid from auction

```
bitsongd tx marketplace claim-bid \
--auction-id=<auction_id> \
--from=<bidder> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx marketplace claim-bid --auction-id=1 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```
