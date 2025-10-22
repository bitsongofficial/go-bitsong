#!/bin/bash
####################################################################
# A. START
####################################################################

# bitsongd sub-1 ./data 26657 26656 6060 9090 ubtsg
BIND=bitsongd
CV=cosmovisor
CHAINID=test-1
CHAINDIR=./data
DAEMON_HOME=$CHAINDIR/$CHAINID/val1

# upgrade 
OLD_RELEASE_GIT=https://github.com/permissionlessweb/go-bitsong
OLD_TAG=feat/v0.23.0

# cosmovisor 
REALPATH=$(realpath)
export DAEMON_HOME="$REALPATH/data/test-1/val1"
export DAEMON_NAME=$BIND
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export COSMOVISOR_CUSTOM_PREUPGRADE="preupgrade.sh"

# Define the new ports for val1
VAL1_API_PORT=1317
VAL1_GRPC_PORT=9090
VAL1_GRPC_WEB_PORT=9091
VAL1_PROXY_APP_PORT=26658
VAL1_RPC_PORT=26657
VAL1_PPROF_PORT=6060
VAL1_P2P_PORT=26656

 
# upgrade details
UPGRADE_VERSION_TITLE="v0.24.0"
UPGRADE_VERSION_TAG="v024"
UPGRADE_INFO="https://raw.githubusercontent.com/permissionlessweb/networks/refs/heads/master/bitsong-2b/upgrades/v0.24.0/cosmovisor.json"

echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "Creating $BINARY instance for VAL1: home=$DAEMON_HOME | chain-id=$CHAINID | p2p=:$VAL1_P2P_PORT | rpc=:$VAL1_RPC_PORT | profiling=:$VAL1_PPROF_PORT | grpc=:$VAL1_GRPC_PORT"
trap 'pkill -f '"$BIND" EXIT
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"

defaultCoins="100000000000000ubtsg"  # 1M
fundCommunityPool="1000000000ubtsg" # 1K
delegate="1000000ubtsg" # 1btsg

rm -rf $DAEMON_HOME  
# install bitsongd in background
(
    echo "Starting repository clone and build..."
    git clone $OLD_RELEASE_GIT
    cd go-bitsong &&
    git checkout $OLD_TAG
    make install 
    cd ../
    echo " Repository clone and build completed"
) &


rm -rf $DAEMON_HOME/test-keys
$BIND init $CHAINID --overwrite --home $DAEMON_HOME --chain-id $CHAINID
mkdir $DAEMON_HOME/test-keys
$BIND --home $DAEMON_HOME config keyring-backend test
sleep 1


#       
# modify val1 genesis 
jq ".app_state.crisis.constant_fee.denom = \"ubtsg\" |
      .app_state.staking.params.bond_denom = \"ubtsg\" |
      .app_state.mint.params.blocks_per_year = \"10000000\" |
      .app_state.mint.params.mint_denom = \"ubtsg\" |
      .app_state.gov.voting_params.voting_period = \"30s\" |
      .app_state.gov.params.expedited_voting_period = \"10s\" | 
      .app_state.gov.params.voting_period = \"15s\" |
      .app_state.gov.params.min_deposit[0].denom = \"ubtsg\" |
      .app_state.gov.params.expedited_min_deposit[0].denom = \"ubtsg\" |
      .app_state.fantoken.params.burn_fee.denom = \"ubtsg\" |
      .app_state.fantoken.params.issue_fee.denom = \"ubtsg\" |
      .app_state.slashing.params.signed_blocks_window = \"10\" |
      .app_state.slashing.params.min_signed_per_window = \"1.000000000000000000\" |
      .app_state.fantoken.params.mint_fee.denom = \"ubtsg\"" $DAEMON_HOME/config/genesis.json > $DAEMON_HOME/config/tmp.json
# give val2 a genesis
mv $DAEMON_HOME/config/tmp.json $DAEMON_HOME/config/genesis.json

# setup test keys.
yes | $BIND  --home $DAEMON_HOME keys add validator1 --output json > $DAEMON_HOME/test-keys/val.json 2>&1 &&
yes | $BIND  --home $DAEMON_HOME keys add user --output json > $DAEMON_HOME/test-keys/user.json 2>&1 &&
yes | $BIND  --home $DAEMON_HOME keys add delegator1 --output json > $DAEMON_HOME/test-keys/del.json 2>&1 &&
$BIND --home $DAEMON_HOME genesis add-genesis-account "$($BIND --home $DAEMON_HOME keys show user -a)" $defaultCoins &&
$BIND --home $DAEMON_HOME genesis add-genesis-account "$($BIND --home $DAEMON_HOME keys show validator1 -a)" $defaultCoins &&
$BIND --home $DAEMON_HOME genesis add-genesis-account "$($BIND --home $DAEMON_HOME keys show delegator1 -a)" $defaultCoins &&
$BIND --home $DAEMON_HOME genesis gentx validator1 $delegate --chain-id $CHAINID &&
$BIND genesis collect-gentxs --home $DAEMON_HOME &&


# keys 
DEL1=$(jq -r '.name' $CHAINDIR/$CHAINID/val1/test-keys/del.json)
DEL1ADDR=$(jq -r '.address' $CHAINDIR/$CHAINID/val1/test-keys/del.json)
VAL1=$(jq -r '.name' $CHAINDIR/$CHAINID/val1/test-keys/val.json)
USERADDR=$(jq -r '.address'  $CHAINDIR/$CHAINID/val1/test-keys/user.json)


# app & config modiifications
# config.toml
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL1_PROXY_APP_PORT\"/g" $DAEMON_HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $DAEMON_HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $DAEMON_HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL1_P2P_PORT\"/" $DAEMON_HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $DAEMON_HOME/config/config.toml &&
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"8s\"/" "$DAEMON_HOME/config/config.toml"

# app.toml
sed -i.bak "/^\[api\]/,/^\[/ s/minimum-gas-prices.*/minimum-gas-prices = \"0.0ubtsg\"/" $DAEMON_HOME/config/app.toml &&
sed -i.bak "/^\[api\]/,/^\[/ s/address.*/address = \"tcp:\/\/0.0.0.0:$VAL1_API_PORT\"/" $DAEMON_HOME/config/app.toml &&
sed -i.bak "/^\[grpc\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_PORT\"/" $DAEMON_HOME/config/app.toml &&
sed -i.bak "/^\[grpc-web\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_WEB_PORT\"/" $DAEMON_HOME/config/app.toml &&
 
####################################################################
# C. COSMOVISOR
####################################################################
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
$CV init "$HOME"/go/bin/$BIND
curl -o $DAEMON_HOME/cosmovisor/preupgrade.sh https://raw.githubusercontent.com/permissionlessweb/networks/refs/heads/master/bitsong-2b/upgrades/v0.24.0/preUpgradeScript.sh
 
$CV run genesis validate-genesis  --home $DAEMON_HOME
# Start bitsong
echo "Starting Genesis validator..."
$CV run start --home $DAEMON_HOME & 
VAL1_PID=$!
echo "VAL1_PID: $VAL1_PID"
sleep 7

####################################################################
# C. UPGRADE
####################################################################
echo "lets upgrade "
sleep 6

LATEST_HEIGHT=$( $BIND status --home $DAEMON_HOME | jq -r '.sync_info.latest_block_height' )
UPGRADE_HEIGHT=$(( $LATEST_HEIGHT + 15 ))
echo "UPGRADE HEIGHT: $UPGRADE_HEIGHT"
sleep 6


cat <<EOF > "$DAEMON_HOME/upgrade.json" 
{
 "messages": [
  {
   "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
   "authority": "bitsong10d07y265gmmuvt4z0w9aw880jnsr700jktpd5u",
   "plan": {
    "name": "$UPGRADE_VERSION_TAG",
    "time": "0001-01-01T00:00:00Z",
    "height": "$UPGRADE_HEIGHT",
    "info": "$UPGRADE_INFO",
    "upgraded_client_state": null
   }
  }
 ],
 "metadata": "ipfs://CID",
 "deposit": "5000000000ubtsg",
 "title": "$UPGRADE_VERSION_TITLE",
 "summary": "mememe",
 "expedited": true 
}
EOF

echo "propose upgrade using expedited proposal..."
$CV run tx gov submit-proposal $DAEMON_HOME/upgrade.json --gas auto --gas-adjustment 1.5 --fees="2000ubtsg" --chain-id=$CHAINID --home=$DAEMON_HOME --from="$VAL1" -y
sleep 6

# echo "vote upgrade"
$CV run tx gov vote 1 yes --from "$DEL1" --gas auto --gas-adjustment 1.2 --fees 1000ubtsg --chain-id $CHAINID --home $DAEMON_HOME -y
$CV run tx gov vote 1 yes --from "$VAL1" --gas auto --gas-adjustment 1.2 --fees 1000ubtsg --chain-id $CHAINID --home $DAEMON_HOME -y
sleep 10
 
####################################################################
# C. CONFIRM
####################################################################
echo "performing v024 upgrade"
sleep 120
pkill -f $BIND
echo "UPGRADE APPLIED SUCCESSFULLY"


