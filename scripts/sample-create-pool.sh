#!/bin/bash

# Set localnet configuration
# Reference localnet script to see which tokens are given to the user accounts in genesis state
BINARY=bitsongd
CHAIN_ID=localnet
CHAIN_DIR=./data
USER_1_ADDRESS=bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw
USER_2_ADDRESS=bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw

# Ensure jq is installed
if [[ ! -x "$(which jq)" ]]; then
  echo "jq (a tool for parsing json in the command line) is required..."
  echo "https://stedolan.github.io/jq/download/"
  exit 1
fi

# Ensure bitsongd is installed
if ! [ -x "$(which $BINARY)" ]; then
  echo "Error: bitsongd is not installed. Try building $BINARY by 'make install'" >&2
  exit 1
fi

# Ensure localnet is running
if [[ "$(pgrep $BINARY)" == "" ]];then
    echo "Error: localnet is not running. Try running localnet by 'make localnet" 
    exit 1
fi

# bitsongd q bank balances bitsong1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu --home ./data/localnet --output json | jq
echo "-> Checking user1 account balances..."
$BINARY q bank balances $USER_1_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# bitsongd q bank balances bitsong185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny --home ./data/localnet --output json | jq
echo "-> Checking user2 account balances..."
$BINARY q bank balances $USER_2_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# bitsongd tx liquidity create-pool 1 100000000ubtsg,100000000token --home ./data/localnet --chain-id localnet --from user1 --keyring-backend test --yes
echo "-> Creating liquidity pool 1..."
$BINARY tx liquidity create-pool 1 100000000ubtsg,100000000token \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user1 \
--keyring-backend test \
--yes

sleep 2

# bitsongd tx liquidity create-pool 1 100000000ubtsg,100000000atom --home ./data/localnet --chain-id localnet --from user2 --keyring-backend test --yes
echo "-> Creating liquidity pool 2..."
$BINARY tx liquidity create-pool 1 100000000ubtsg,100000000atom \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user2 \
--keyring-backend test \
--yes

sleep 2

# bitsongd q bank balances bitsong1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu --home ./data/localnet --output json | jq
echo "-> Checking user1 account balances after..."
$BINARY q bank balances $USER_1_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# bitsongd q bank balances bitsong185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny --home ./data/localnet --output json | jq
echo "-> Checking user2 account balances after..."
$BINARY q bank balances $USER_2_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# bitsongd q liquidity pools --home ./data/localnet --output json | jq
echo "-> Querying liquidity pools..."
$BINARY q liquidity pools \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq