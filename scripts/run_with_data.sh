#!/usr/bin/env bash

#PASSWORD="Bitsong@123"

rm -rf "$HOME"/.bitsongd
rm -rf "$HOME"/.bitsongcli

make install # assumes currently in project directory

bitsongd init local --chain-id test-chain

bitsongcli config chain-id test-chain
bitsongcli config output json
bitsongcli config indent true
bitsongcli config trust-node true
bitsongcli config keyring-backend test

 bitsongcli keys delete raneet --force
 bitsongcli keys delete angelo --force

 bitsongcli keys add raneet
 bitsongcli keys add angelo


# Note: important to add 'raneet' as a genesis-account since this is the chain's validator
 bitsongd add-genesis-account "$(bitsongcli keys show raneet -a)" 1000btsg,100000000000ubtsg
 bitsongd add-genesis-account "$(bitsongcli keys show angelo -a)" 1000btsg,100000000000ubtsg


# Set staking token (both bond_denom and mint_denom)
STAKING_TOKEN="ubtsg"
FROM="\"bond_denom\": \"stake\""
TO="\"bond_denom\": \"$STAKING_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bitsongd/config/genesis.json
FROM="\"mint_denom\": \"stake\""
TO="\"mint_denom\": \"$STAKING_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bitsongd/config/genesis.json

# Set fee token (both for gov min deposit and crisis constant fee)
FEE_TOKEN="ubtsg"
FROM="\"stake\""
TO="\"$FEE_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bitsongd/config/genesis.json

# Set min-gas-prices (using fee token)
FROM="minimum-gas-prices = \"\""
TO="minimum-gas-prices = \"0.025$FEE_TOKEN\""
sed -i "s/$FROM/$TO/" "$HOME"/.bitsongd/config/app.toml

 bitsongd gentx --name raneet --amount 1000000ubtsg --keyring-backend test

echo "Collecting genesis txs..."
bitsongd collect-gentxs

echo "Validating genesis file..."
bitsongd validate-genesis

# Uncomment the below to broadcast node RPC endpoint
#FROM="laddr = \"tcp:\/\/127.0.0.1:26657\""
#TO="laddr = \"tcp:\/\/0.0.0.0:26657\""
#sed -i "s/$FROM/$TO/" "$HOME"/.bitsongd/config/config.toml

# Uncomment the below to broadcast REST endpoint
# Do not forget to comment the bottom lines !!
#bitsongd start --pruning "everything" &
#bitsongcli rest-server --chain-id test --laddr="tcp://0.0.0.0:1317" --trust-node && fg

bitsongd start --pruning "everything" &
bitsongcli rest-server --chain-id test-chain --trust-node && fg