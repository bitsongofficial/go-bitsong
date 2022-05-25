#!/usr/bin/env bash

bitsongd query nft metadata 1
bitsongd query nft nft-info 1
bitsongd query nft collection 1
bitsongd query nft nfts-by-owner $(bitsongd keys show -a validator --keyring-backend=test)
          
bitsongd tx nft create-nft --name="Punk10" --symbol="PUNK" --uri="https://punk.com/10" --seller-fee-basis-points=100 --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=false --update-authority="$(bitsongd keys show -a validator --keyring-backend=test)" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft sign-metadata --metadata-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block       
bitsongd tx nft transfer-nft --nft-id=1 --new-owner="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft update-metadata --metadata-id=1 --name="Punk11" --symbol="PUNK" --uri="https://punk.com/10" --seller-fee-basis-points=100 --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft update-metadata-authority --metadata-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft verify-collection --collection-id=1 --nft-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft unverify-collection --collection-id=1 --nft-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft update-collection-authority --collection-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
