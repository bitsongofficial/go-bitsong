# Client

## Query

### Params

Query candy machine module parameters

```sh
bitsongd query candymachine params
```

### Candymachine

Query a single candy machine by collection id

```
bitsongd query candymachine candymachine [collection_id]
bitsongd query candymachine candymachine 1
```

### Candymachines

Query all candy machines

```
bitsongd query candymachine machines
```

## Messages

### CreateCandyMachine

Create candy machine from parameters.

```sh
bitsongd tx candymachine create-candymachine \
 --collection-id=<collection_id>  \
 --price=<nft_mint_price> \
 --denom=<denom> \
 --metadata-baseurl=<metadata_base_url> \
 --end-timestamp=<end_timestamp> \
 --max-mint=<max_mint> \
 --treasury=<treasury_to_receive_payment> \
 --go-live-date=<candymachine_open_mint_timestamp> \
 --creators=<creators_of_collection_nfts> \
 --creator-shares=<shares_between_creators> \
 --mutable=<mutability_of_minted_nfts> \
 --seller-fee-basis-points=<fee_share_basis_points> \
 --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx candymachine create-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### UpdateCandyMachine

Update candy machine by collection id and parameters.

```sh
bitsongd tx candymachine update-candymachine \
 --collection-id=<collection_id>  \
 --price=<nft_mint_price> \
 --denom=<denom> \
 --metadata-baseurl=<metadata_base_url> \
 --end-timestamp=<end_timestamp> \
 --max-mint=<max_mint> \
 --treasury=<treasury_to_receive_payment> \
 --go-live-date=<candymachine_open_mint_timestamp> \
 --creators=<creators_of_collection_nfts> \
 --creator-shares=<shares_between_creators> \
 --mutable=<mutability_of_minted_nfts> \
 --seller-fee-basis-points=<fee_share_basis_points> \
 --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx candymachine update-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata2" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### CloseCandyMachine

Close candy machine by collection id.

```
bitsongd tx candymachine update-candymachine <collection_id> --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx candymachine close-candymachine 1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### MintNFT

Mint nft from collection id put on candy machine, and custom name.

```
bitsongd tx candymachine mint-nft <collection_id> <nft_name> --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx candymachine mint-nft 1 punk1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```
