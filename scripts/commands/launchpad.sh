#!/usr/bin/env bash

bitsongd query launchpad params
bitsongd query launchpad launchpads
bitsongd query launchpad launchpad 1
bitsongd query launchpad mintabe-metadata-ids 1

bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

# shuffle launchpad
bitsongd tx launchpad create-launchpad --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --shuffle=true --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
# sequence launchpad
bitsongd tx launchpad create-launchpad --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --shuffle=false --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx launchpad update-launchpad --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata2" --end-timestamp=0 --max-mint=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --shuffle=true --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx launchpad close-launchpad 1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx launchpad mint-nft 1 punk1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk2 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk3 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk4 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk5 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk6 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk7 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk8 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk9 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk10 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx launchpad mint-nft 1 punk11 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

# mint multiple nfts at a single time
bitsongd tx launchpad mint-nfts 1 10 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block --gas=1000000
bitsongd tx launchpad mint-nfts 1 3 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block --gas=1000000

bitsongd query nft nfts-by-owner  $(bitsongd keys show -a validator --keyring-backend=test)
bitsongd query nft metadata 1 1
bitsongd query nft metadata 1 2
bitsongd query nft metadata 1 3
bitsongd query nft metadata 1 4
bitsongd query nft metadata 1 5
bitsongd query nft metadata 1 6
bitsongd query nft metadata 1 7
bitsongd query nft metadata 1 8
bitsongd query nft metadata 1 9
bitsongd query nft metadata 1 10
bitsongd query nft metadata 1 11
