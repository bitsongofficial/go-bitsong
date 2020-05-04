#!/bin/bash

# Generate a new StdTx
bitsongcli tx player register test $(bitsongcli keys show faucet -a --keyring-backend=test) $(bitsongcli keys show validator -a --bech val --keyring-backend=test) --generate-only > player_tx.json

# Sign the tx with player address
bitsongcli tx sign player_tx.json --from faucet --keyring-backend=test > player_signed.json

# Sign the tx with validator address
bitsongcli tx sign player_signed.json --from validator --keyring-backend=test > player_signed2.json

# Broadcast the tx
bitsongcli tx broadcast player_signed2.json --from faucet --keyring-backend=test -b block

