# Client

## Query

### Params

Query launchpad module parameters

```sh
bitsongd query launchpad params
```

### Launchpad

Query a single launchpad by collection id

```
bitsongd query launchpad launchpad [collection_id]
bitsongd query launchpad launchpad 1
```

### Launchpads

Query all launchpads

```
bitsongd query launchpad pads
```

### Mintable metadata ids

Query all mintable metadata ids for a collection put on launchpad.

```
bitsongd query launchpad mintabe-metadata-ids [collection_id]
bitsongd query launchpad mintabe-metadata-ids 1
```

## Messages

### CreateLaunchPad

Create launchpad from parameters.

```sh
bitsongd tx launchpad create-launchpad \
 --collection-id=<collection_id>  \
 --price=<nft_mint_price> \
 --denom=<denom> \
 --metadata-baseurl=<metadata_base_url> \
 --end-timestamp=<end_timestamp> \
 --max-mint=<max_mint> \
 --treasury=<treasury_to_receive_payment> \
 --go-live-date=<launchpad_open_mint_timestamp> \
 --creators=<creators_of_collection_nfts> \
 --creator-shares=<shares_between_creators> \
 --mutable=<mutability_of_minted_nfts> \
 --seller-fee-basis-points=<fee_share_basis_points> \
 --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx launchpad create-launchpad --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### UpdateLaunchPad

Update launchpad by collection id and parameters.

```sh
bitsongd tx launchpad update-launchpad \
 --collection-id=<collection_id>  \
 --price=<nft_mint_price> \
 --denom=<denom> \
 --metadata-baseurl=<metadata_base_url> \
 --end-timestamp=<end_timestamp> \
 --max-mint=<max_mint> \
 --treasury=<treasury_to_receive_payment> \
 --go-live-date=<launchpad_open_mint_timestamp> \
 --creators=<creators_of_collection_nfts> \
 --creator-shares=<shares_between_creators> \
 --mutable=<mutability_of_minted_nfts> \
 --seller-fee-basis-points=<fee_share_basis_points> \
 --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx launchpad update-launchpad --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata2" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### CloseLaunchPad

Close launchpad by collection id.

```
bitsongd tx launchpad update-launchpad <collection_id> --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx launchpad close-launchpad 1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### MintNFT

Mint nft from collection id put on launchpad, and custom name.

```
bitsongd tx launchpad mint-nft <collection_id> <nft_name> --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx launchpad mint-nft 1 punk1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
```

### MintNFTs

Mint multiple nfts at a single time. The nft name is determined as `collection_name #{metadataId}`.

```
bitsongd tx launchpad mint-nfts <collection_id> <mint_count> --from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx launchpad mint-nfts 1 3 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block --gas=1000000
```
