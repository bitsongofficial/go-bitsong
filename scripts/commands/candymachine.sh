#!/usr/bin/env bash

bitsongd query candymachine params
bitsongd query candymachine candymachines
bitsongd query candymachine candymachine 1

bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine create-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine update-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata2" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine close-candymachine 1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine mint-nft 1 punk1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk2 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk3 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk4 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk5 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk6 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk7 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk8 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk9 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk10 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx candymachine mint-nft 1 punk11 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd query nft nfts-by-owner  $(bitsongd keys show -a validator --keyring-backend=test)
bitsongd query nft metadata 1
bitsongd query nft metadata 2
bitsongd query nft metadata 3
bitsongd query nft metadata 4
bitsongd query nft metadata 5
bitsongd query nft metadata 6
bitsongd query nft metadata 7
bitsongd query nft metadata 8
bitsongd query nft metadata 9
bitsongd query nft metadata 10
bitsongd query nft metadata 11
