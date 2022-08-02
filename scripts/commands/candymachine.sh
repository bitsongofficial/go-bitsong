#!/usr/bin/env bash

bitsongd query candymachine candymachines

bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine create-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata" --end-settings-type="BY_MINT" --end-settings-value=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine update-candymachine --collection-id=1 --price=1000 --denom=ubtsg --metadata-baseurl="https://punk.com/metadata2" --end-settings-type="BY_MINT" --end-settings-value=10 --treasury=$(bitsongd keys show -a validator --keyring-backend=test) --go-live-date="1659404536" --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=true --seller-fee-basis-points=100 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine close-candymachine 1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx candymachine mint-nft 1 punk1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
