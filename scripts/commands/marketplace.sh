#!/usr/bin/env bash

bitsongd query marketplace auctions
bitsongd query marketplace auction 1
bitsongd query marketplace bids-by-auction 1
bitsongd query marketplace bids-by-bidder $(bitsongd keys show -a validator --keyring-backend=test)
bitsongd query marketplace bidder-metadata $(bitsongd keys show -a validator --keyring-backend=test)

bitsongd tx marketplace create-auction --nft-id=1:1:0 --prize-type="NFT_ONLY_TRANSFER" --bid-denom="ubtsg" --duration="864000s" --price-floor=1000000 --instant-sale-price=100000000 --tick-size=100000 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace set-auction-authority --auction-id=1 --new-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace start-auction --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace end-auction --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace place-bid --auction-id=1 --amount="1000000ubtsg" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace cancel-bid --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace claim-bid --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd keys add bidder1 --keyring-backend=test
bitsongd keys add bidder2 --keyring-backend=test
bitsongd keys add bidder3 --keyring-backend=test
bitsongd tx bank send validator $(bitsongd keys show -a bidder1 --keyring-backend=test) 100000000ubtsg --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx bank send validator $(bitsongd keys show -a bidder2 --keyring-backend=test) 100000000ubtsg --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx bank send validator $(bitsongd keys show -a bidder3 --keyring-backend=test) 100000000ubtsg --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace place-bid --auction-id=1 --amount="2000000ubtsg" --from=bidder1 --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace place-bid --auction-id=1 --amount="2000000ubtsg" --from=bidder2 --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace place-bid --auction-id=1 --amount="3000000ubtsg" --from=bidder3 --chain-id=test --keyring-backend=test -y --broadcast-mode=block

bitsongd tx marketplace create-auction --nft-id=1:1:0 --prize-type="NFT_ONLY_TRANSFER" --bid-denom="ubtsg" --duration="1s" --price-floor=1000000 --instant-sale-price=100000000 --tick-size=100000 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace start-auction --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block

# edition test
bitsongd tx nft create-collection --name="punk-collection" --uri="https://punk.com" --update-authority=$(bitsongd keys show -a validator --keyring-backend=test) --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx nft create-nft --collection-id=1 --name="Punk10" --symbol="PUNK" --uri="https://punk.com/10" --seller-fee-basis-points=100 --creators=$(bitsongd keys show -a validator --keyring-backend=test) --creator-shares="10" --mutable=false --update-authority="$(bitsongd keys show -a validator --keyring-backend=test)" --master-edition-max-supply=4 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace create-auction --nft-id=1:1:0 --prize-type="LIMITED_EDITION_PRINTS" --bid-denom="ubtsg" --duration="864000s" --price-floor=1000000 --instant-sale-price=100000000 --tick-size=100000 --edition-limit=2 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace start-auction --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace end-auction --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace place-bid --auction-id=1 --amount="1000000ubtsg" --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace cancel-bid --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
bitsongd tx marketplace claim-bid --auction-id=1 --from=validator --chain-id=test --keyring-backend=test -y --broadcast-mode=block
