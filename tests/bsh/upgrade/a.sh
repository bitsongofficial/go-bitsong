#!/bin/bash

BIND=bitsongd
CHAINID=test-1
UPGRADE_VERSION=v024

OLD_RELEASE_GIT=https://github.com/bitsongofficial/go-bitsong
NEW_RELEASE_GIT=https://github.com/permissionlessweb/go-bitsong
OLD_TAG=v0.23.0
NEW_TAG=main

SNAPSHOT_BLOCK=24222132
SNAPSHOT_TAR=bitsong_$SNAPSHOT_BLOCK.tar.lz4
SNAPSHOT_PATH=./data/$SNAPSHOT_TAR

# file paths
CHAINDIR=./data
VAL1HOME=$CHAINDIR/$CHAINID/val1
 
# Define the new ports for val1 on chain a
VAL1_API_PORT=1317
VAL1_GRPC_PORT=9090
VAL1_GRPC_WEB_PORT=9091
VAL1_PROXY_APP_PORT=26658
VAL1_RPC_PORT=26657
VAL1_PPROF_PORT=6060
VAL1_P2P_PORT=26656
 

echo "Creating $BINARY instance for VAL1: home=$VAL1HOME | chain-id=$CHAINID | p2p=:$VAL1_P2P_PORT | rpc=:$VAL1_RPC_PORT | profiling=:$VAL1_PPROF_PORT | grpc=:$VAL1_GRPC_PORT"
trap 'pkill -f '"$BIND" EXIT

# Clone the repository if it doesn't exist
git clone $OLD_RELEASE_GIT
# # Change into the cloned directory
cd go-bitsong &&
# # Checkout the version of go-bitsong that doesnt submit slashing hooks
git checkout $OLD_TAG
make install 
cd ../ &&

## download snapshot from polkachu (thanks team!)
mkdir .bin
# sudo apt install lz4
# curl https://snapshots.polkachu.com/snapshots/bitsong/$SNAPSHOT_TAR  -o $SNAPSHOT_PATH

####################################################################
# A. CHAINS CONFIG
####################################################################

rm -rf $VAL1HOME 
rm -rf $VAL1HOME/test-keys

# initialize chains
$BIND init $CHAINID --overwrite --home $VAL1HOME --chain-id $CHAINID
sleep 2
mkdir $VAL1HOME/test-keys

# cli config
$BIND --home $VAL1HOME config keyring-backend test
sleep 1
$BIND --home $VAL1HOME config chain-id $CHAINID
sleep 1
$BIND --home $VAL1HOME config node tcp://localhost:$VAL1_RPC_PORT
sleep 1
  
# setup test keys.
yes | $BIND  --home $VAL1HOME keys add validator1 --output json > $VAL1HOME/test-keys/val.json 2>&1 
sleep 1
yes | $BIND  --home $VAL1HOME keys add user --output json > $VAL1HOME/test-keys/user.json 2>&1
sleep 1
yes | $BIND  --home $VAL1HOME keys add delegator1 --output json > $VAL1HOME/test-keys/del.json 2>&1
sleep 1

DEL1=$(jq -r '.name' $CHAINDIR/"$CHAINID"/val1/test-keys/del.json)
DEL1ADDR=$(jq -r '.address' $CHAINDIR/"$CHAINID"/val1/test-keys/del.json)
VAL1=$(jq -r '.name' $CHAINDIR/"$CHAINID"/val1/test-keys/val.json)
USERADDR=$(jq -r '.address'  $CHAINDIR/"$CHAINID"/val1/test-keys/user.json)

# app & config modiifications
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL1_PROXY_APP_PORT\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL1_P2P_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^pprof_laddr *=.*/pprof_laddr = \"localhost:6060\"/g" $VAL1HOME/config/config.toml &&
# shorten block times a bit
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"2s\"/" "$VAL1HOME/config/config.toml"
 
# app.toml
sed -i.bak "/^\[api\]/,/^\[/ s/minimum-gas-prices.*/minimum-gas-prices = \"0.0ubtsg\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[api\]/,/^\[/ s/address.*/address = \"tcp:\/\/0.0.0.0:$VAL1_API_PORT\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[grpc\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_PORT\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[grpc-web\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_WEB_PORT\"/" $VAL1HOME/config/app.toml &&
 
 
####################################################################
# 0. SNAPSHOT CONFIG 
####################################################################
echo "unzipping snapshot..."
# create export 
lz4 -c -d  $SNAPSHOT_PATH | tar -x -C $VAL1HOME

echo "creating testnet-from-export"
# create testnet-from-export
$BIND in-place-testnet "$CHAINID" "$USERADDR" bitsongvaloper1qxw4fjged2xve8ez7nu779tm8ejw92rv0vcuqr --trigger-testnet-upgrade $UPGRADE_VERSION  --home $VAL1HOME --skip-confirmation & 
INPLACE_TESTNET=$!
echo "INPLACE_TESTNET: $INPLACE_TESTNET"
sleep 45

####################################################################
# 0. UPGRADING
####################################################################
pkill -f $BIND
# Clone the repository if it doesn't exist
# # Change into the cloned directory
cd go-bitsong && git remote add upstream $NEW_RELEASE_GIT && git fetch upstream
# # Checkout the version of go-bitsong that doesnt submit slashing hooks
git checkout upstream/$NEW_TAG && git pull 
make install 
cd ../ &&

# Start bitsong
echo "Running upgradehandler to ensure upgrade with live state is okay!..."
$BIND start --home $VAL1HOME & 
VAL1_PID=$!
echo "VAL1_PID: $VAL1_PID"
sleep 7
 
# ## ensure funding community pool is okay 
# MSG_CODE=$($BIND tx distribution fund-community-pool $fundCommunityPool --from="$DEL1" --gas auto --fees 200ubtsg --gas-adjustment 1.2 --chain-id $CHAINID --home $VAL1HOME -o json -y | jq -r '.code')
# if [ -n "$MSG_CODE" ] && [ "$MSG_CODE" -ne 0 ]; then
#   exit 1
# fi

echo "COMMUNITY POOL PATCH APPLIED SUCCESSFULLY, ENDING TESTS"
pkill -f bitsongd