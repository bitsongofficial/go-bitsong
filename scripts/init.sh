#!/bin/bash

PASSWORD="12345678"

rm -rf ~/.bitsong*

# Initialize the genesis.json file that will help you to bootstrap the network
bitsongd init MyValidator --chain-id=bitsong-dev-network-1

bitsongcli config chain-id bitsong-dev-network-1
bitsongcli config output json
bitsongcli config indent true
bitsongcli config trust-node true

# Change default bond token genesis.json
sed -i 's/stake/ubtsg/g' ~/.bitsongd/config/genesis.json

# Create a key to hold your validator account
echo ${PASSWORD} | bitsongcli keys add validator

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some coins
# with the genesis.app_state.staking.params.bond_denom denom, the default is staking
bitsongd add-genesis-account validator 1000000000000000000ubtsg

# Generate the transaction that creates your validator
echo ${PASSWORD} | bitsongd gentx --name validator --amount=1000000000000ubtsg

# Add the generated bonding transaction to the genesis file
bitsongd collect-gentxs
bitsongd validate-genesis

# Now its safe to start `bitsongd`
bitsongd start
