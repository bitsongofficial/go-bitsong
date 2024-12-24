#!/bin/bash
# Run this script to quickly install, setup, and run the current version of bitsong without docker.
#
# Example:
# CHAIN_ID="local-1" HOME_DIR="~/.btsg1" TIMEOUT_COMMIT="500ms" CLEAN=true sh scripts/test_node.sh
# CHAIN_ID="local-2" HOME_DIR="~/.btsg2" CLEAN=true RPC=36657 REST=2317 PROFF=6061 P2P=36656 GRPC=8090 GRPC_WEB=8091 ROSETTA=8081 TIMEOUT_COMMIT="500ms" sh scripts/test_node.sh
#
# To use unoptomized wasm files up to ~5mb, add: MAX_WASM_SIZE=5000000

export KEY="btsg1"
export KEY2="btsg2"

export CHAIN_ID=${CHAIN_ID:-"local-1"}
export MONIKER="localbitsong"
export KEYALGO="secp256k1"
export KEYRING=${KEYRING:-"test"}
export HOME_DIR=$(eval echo "${HOME_DIR:-"~/.bitsongd"}")
export BINARY=${BINARY:-bitsongd}

export CLEAN=${CLEAN:-"true"}
export RPC=${RPC:-"26657"}
export REST=${REST:-"1317"}
export PROFF=${PROFF:-"6060"}
export P2P=${P2P:-"26656"}
export GRPC=${GRPC:-"9090"}
export GRPC_WEB=${GRPC_WEB:-"9091"}
export ROSETTA=${ROSETTA:-"8080"}
export TIMEOUT_COMMIT=${TIMEOUT_COMMIT:-"5s"}

alias BINARY="$BINARY --home=$HOME_DIR"

command -v $BINARY > /dev/null 2>&1 || { echo >&2 "$BINARY command not found. Ensure this is setup / properly installed in your GOPATH (make install)."; exit 1; }
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

$BINARY config keyring-backend $KEYRING
$BINARY config chain-id $CHAIN_ID

from_scratch () {
  # Fresh install on current branch
  make install

  # remove existing daemon.
  rm -rf $HOME_DIR && echo "Removed $HOME_DIR"

  # bitsong1e7f6j0006x5csqra593uvm7q850ek72jke37z9
  echo "satoshi pioneer scrub lunch survey connect cable brick meadow wife way possible become visit puzzle track north bleak clap dentist pattern awake orange grit" | BINARY keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
  # bitsong1ww54e9y0ysfq3ljglj9rx3qvcekksksju93l5f
  echo "course render gym design oppose focus van impulse surround chase lock setup unhappy priority mammal video chapter under rocket child sadness always coyote sample" | BINARY keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO --recover

  BINARY init $MONIKER --chain-id $CHAIN_ID --default-denom ubtsg

  # Function updates the config based on a jq argument as a string
  update_test_genesis () {
    cat $HOME_DIR/config/genesis.json | jq "$1" > $HOME_DIR/config/tmp_genesis.json && mv $HOME_DIR/config/tmp_genesis.json $HOME_DIR/config/genesis.json
  }

  # Block
  update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
  # Gov
  update_test_genesis '.app_state["gov"]["params"]["min_deposit"]=[{"denom": "ubtsg","amount": "1000000"}]'
  update_test_genesis '.app_state["gov"]["params"]["voting_period"]="15s"'
  update_test_genesis '.app_state["gov"]["params"]["expedited_voting_period"]="6s"'
  # staking
  update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="ubtsg"'
  update_test_genesis '.app_state["staking"]["params"]["min_commission_rate"]="0.050000000000000000"'
  # mint
  update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="ubtsg"'
  # crisis
  update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "ubtsg","amount": "1000"}'

  # Custom Modules

  # Allocate genesis accounts
  BINARY genesis add-genesis-account $KEY 10000000ubtsg,1000utest --keyring-backend $KEYRING
  BINARY genesis add-genesis-account $KEY2 1000000ubtsg,1000utest --keyring-backend $KEYRING
  BINARY genesis add-genesis-account bitsong1e7f6j0006x5csqra593uvm7q850ek72jke37z9 100000000000ubtsg --keyring-backend $KEYRING

  BINARY genesis gentx $KEY 1000000ubtsg --keyring-backend $KEYRING --chain-id $CHAIN_ID

  # Collect genesis tx
  BINARY genesis collect-gentxs

  # Run this to ensure the genesis file is setup correctly
  BINARY genesis validate-genesis
}

# check if CLEAN is not set to false
if [ "$CLEAN" != "false" ]; then
  echo "Starting from a clean state"
  from_scratch
fi

echo "Starting node..."

# Opens the RPC endpoint to outside connections
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:'$RPC'"/g' $HOME_DIR/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' $HOME_DIR/config/config.toml

# REST endpoint
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:'$REST'"/g' $HOME_DIR/config/app.toml
sed -i 's/enable = false/enable = true/g' $HOME_DIR/config/app.toml

# replace pprof_laddr = "localhost:6060" binding
sed -i 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:'$PROFF_LADDER'"/g' $HOME_DIR/config/config.toml

# change p2p addr laddr = "tcp://0.0.0.0:26656"
sed -i 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:'$P2P'"/g' $HOME_DIR/config/config.toml

# GRPC
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:'$GRPC'"/g' $HOME_DIR/config/app.toml
sed -i 's/address = "localhost:9091"/address = "0.0.0.0:'$GRPC_WEB'"/g' $HOME_DIR/config/app.toml

# Rosetta Api
# sed -i 's/address = ":8080"/address = "0.0.0.0:'$ROSETTA'"/g' $HOME_DIR/config/app.toml

# faster blocks
sed -i 's/timeout_commit = "5s"/timeout_commit = "'$TIMEOUT_COMMIT'"/g' $HOME_DIR/config/config.toml

# Start the node with 0 gas fees
BINARY start --pruning=nothing  --minimum-gas-prices=0.0025ubtsg --rpc.laddr="tcp://0.0.0.0:$RPC"