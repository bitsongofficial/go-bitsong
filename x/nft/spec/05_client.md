# Client

## Queries

### collection

Query collection by id

```
bitsongd query nft collection <collection_id>
bitsongd query nft collection 1
```

### metadata

Query metadata by Id

```
bitsongd query nft metadata <metadata_id>
bitsongd query nft metadata 1
```

### nft

Query nft info by Id

```
# Note: nft id is expressed as <collection_id>:<metadata_id>:<sequence>
bitsongd query nft nft-info <nft_id>
bitsongd query nft nft-info 1:1:0
```

Query nfts by owner

```
bitsongd query nft nfts-by-owner <nft_owner>

bitsongd query nft nfts-by-owner $(bitsongd keys show -a validator --keyring-backend=test)
```

## Transactions

### CreateCollection

Create collection from name, uri and authority params

```
bitsongd tx nft create-collection \
--name=<collection_name> \
--uri=<collection_uri> \
--update-authority=<update_authority> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block


bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### UpdateCollectionAuthority

Update collection authority by previous authority

```
bitsongd tx nft update-collection-authority \
--collection-id=<collection_id> \
--new-authority=<new_authority> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft update-collection-authority --collection-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### CreateNft

Create Nft by collection owner with name, symbol, uri, seller fee basis points, creators, creator shares, mutability and update authority.

```
bitsongd tx nft create-nft \
--collection-id=<collection_id> \
--name=<nft_name> \
--symbol=<symbol> \
--uri=<uri> \
--seller-fee-basis-points=<seller_fee_basis_points> \
--creators=<creators_list_separated_by_comma> \
--creator-shares=<shares_list_separated_by_comma> \
--mutable=<mutability> \
--update-authority=<update_authority> \
--from=<creator> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft create-nft --collection-id=1 --name="Punk10" --symbol="PUNK" --uri="https://punk.com/10" --seller-fee-basis-points=100 --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=false --update-authority="$(bitsongd keys show -a validator --keyring-backend=test)" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### PrintEdition

Print edition from master edition nft.

```
bitsongd tx nft print-edition \
--collection-id=<collection_id> \
--metadata-id=<metadata_id> \
--owner=<owner_of_new_edition> \
--from=<printer> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block


bitsongd tx nft print-edition --collection-id=1 --metadata-id=1 --owner=$(bitsongd keys show -a validator --keyring-backend=test --home=$HOME/.bitsongd) --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### SignMetadata

Sign the metadata creators field by creator account configured on metadata on metadata creation.

```
bitsongd tx nft sign-metadata \
--metadata-id=<metadata_id> \
--from=<signer> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft sign-metadata --metadata-id=1 --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### TransferNft

Transfer nft from original owner to new owner

```
bitsongd tx nft transfer-nft \
--nft-id=<nft_id> \
--new-owner=<new_owner> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft transfer-nft --nft-id=1:1:0 --new-owner="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### UpdateMetadata

Update metadata with new name, symbol, uri, seller fee basis points, creators and creator shares.

```
bitsongd tx nft update-metadata \
--metadata-id=<metadata_id> \
--name=<name> \
--symbol=<symbol> \
--uri=<uri> \
--seller-fee-basis-points=<seller_fee_basis_points> \
--creators=<creators_list_separated_by_comma> \
--creator-shares=<shares_list_separated_by_comma> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft update-metadata --metadata-id=1 --name="Punk11" --symbol="PUNK" --uri="https://punk.com/10" --seller-fee-basis-points=100 --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### UpdateMetadataAuthority

Update metadata authority by previous metadata authority

```
bitsongd tx nft update-metadata-authority \
--metadata-id=<metadata_id> \
--new-authority=<new_authority> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft update-metadata-authority --metadata-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```

### UpdateMintAuthority

Update mint authority by previous mint authority.

```
bitsongd tx nft update-mint-authority \
--metadata-id=<metadata_id> \
--new-authority=<new_authority> \
--from=<sender> --chain-id=<chain_id> --keyring-backend=<keyring> -y --broadcast-mode=block

bitsongd tx nft update-mint-authority --metadata-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=nft-devnet-1 --keyring-backend=test -y --broadcast-mode=block
```
