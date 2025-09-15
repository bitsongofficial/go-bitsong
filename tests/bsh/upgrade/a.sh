#!/bin/bash

BIND=bitsongd
CHAINID=test-1
UPGRADE_VERSION=v024

OLD_RELEASE_GIT=https://github.com/bitsongofficial/go-bitsong
NEW_RELEASE_GIT=https://github.com/permissionlessweb/go-bitsong
OLD_TAG=main
NEW_TAG=main

SNAPSHOT_PATH=./data/bitsong-snapshot.tar.lz4
SNAPSHOT_URL=$(curl -s "https://www.polkachu.com/tendermint_snapshots/bitsong" | grep -o 'https://snapshots\.polkachu\.com/snapshots/bitsong/bitsong_[0-9]*\.tar\.lz4' | head -1)
echo "Current snapshot url: $SNAPSHOT_URL"

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


echo "Creating $BIND instance for VAL1: home=$VAL1HOME | chain-id=$CHAINID | p2p=:$VAL1_P2P_PORT | rpc=:$VAL1_RPC_PORT | profiling=:$VAL1_PPROF_PORT | grpc=:$VAL1_GRPC_PORT"
trap 'pkill -f '"$BIND" EXIT

# install bitsongd in background
(
    echo "Starting repository clone and build..."
    git clone $OLD_RELEASE_GIT
    cd go-bitsong &&
    git checkout $OLD_TAG
    make install 
    cd ../
    echo "✅ Repository clone and build completed"
) &
BUILD_PID=$!

# download live snapshot in background
(
    echo "Starting snapshot download..."
    mkdir $CHAINDIR
    if [ -f "$SNAPSHOT_PATH" ]; then
        echo "Snapshot already exists at $SNAPSHOT_PATH, skipping download"
    else
        echo "Downloading snapshot from $SNAPSHOT_URL"
        curl "$SNAPSHOT_URL" -o "$SNAPSHOT_PATH"
    fi
    echo "✅ Snapshot download completed"
) &
DOWNLOAD_PID=$!

# wait for both to complete
echo "Waiting for clone/build and snapshot download to complete..."
wait $BUILD_PID
BUILD_EXIT=$?
wait $DOWNLOAD_PID  
DOWNLOAD_EXIT=$?

# Check if both succeeded
if [ $BUILD_EXIT -eq 0 ] && [ $DOWNLOAD_EXIT -eq 0 ]; then
    echo "✅ Both operations completed successfully"
else
    echo "❌ One or both operations failed (Build: $BUILD_EXIT, Download: $DOWNLOAD_EXIT)"
    exit 1
fi

####################################################################
# A. CHAINS CONFIG
####################################################################

rm -rf $VAL1HOME 
rm -rf $VAL1HOME/test-keys

# initialize chains
$BIND init $CHAINID --overwrite --home $VAL1HOME --chain-id $CHAINID &&
mkdir $VAL1HOME/test-keys

# cli config
$BIND --home $VAL1HOME config keyring-backend test &&
$BIND --home $VAL1HOME config chain-id $CHAINID &&
$BIND --home $VAL1HOME config node tcp://localhost:$VAL1_RPC_PORT &&
  
# setup test keys.
yes | $BIND  --home $VAL1HOME keys add validator1 --output json > $VAL1HOME/test-keys/val.json 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add user --output json > $VAL1HOME/test-keys/user.json 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add delegator1 --output json > $VAL1HOME/test-keys/del.json 2>&1 &&

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
sleep 55

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

# TODO: run basic test suite ensuring functionality post upgrade

if kill -0 $VAL1_PID 2>/dev/null; then
    echo "SUCCESS: Node started successfully without panic"
    kill $VAL1_PID 2>/dev/null
    wait $VAL1_PID 2>/dev/null
    exit 0
else
    echo "FAILED: Node process died (likely panicked)"
    exit 1
fi