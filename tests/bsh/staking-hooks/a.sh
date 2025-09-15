#!/bin/bash
####################################################################
# A. START
####################################################################
# TODO: refactor to build both binaries and 2 networks in parallel
# bitsongd sub-1 ./data 26657 26656 6060 9090 ubtsg
BIND=bitsongd
CHAINID=test-1
CHAINDIR=./data
defaultCoins="100000000000ubtsg"  # 100K
nonSlashedDelegation="100000000ubtsg" # 100
delegate="1000000ubtsg" # 1btsg

VAL1HOME=$CHAINDIR/$CHAINID/val1
VAL2HOME=$CHAINDIR/$CHAINID/val2
VAL3HOME=$CHAINDIR/$CHAINID/val3
# Define the new ports for val1
VAL1_API_PORT=1317
VAL1_GRPC_PORT=9090
VAL1_GRPC_WEB_PORT=9091
VAL1_PROXY_APP_PORT=26658
VAL1_RPC_PORT=26657
VAL1_PPROF_PORT=6060
VAL1_P2P_PORT=26656

# Define the new ports for val2
VAL2_API_PORT=1318
VAL2_GRPC_PORT=9393
VAL2_GRPC_WEB_PORT=9394
VAL2_PROXY_APP_PORT=9395
VAL2_RPC_PORT=26357
VAL2_PPROF_PORT=6361
VAL2_P2P_PORT=26356

echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "Creating $BINARY instance for VAL1: home=$VAL1HOME | chain-id=$CHAINID | p2p=:$VAL1_P2P_PORT | rpc=:$VAL1_RPC_PORT | profiling=:$VAL1_PPROF_PORT | grpc=:$VAL1_GRPC_PORT"
echo "Creating $BINARY instance for VAL2: home=$VAL2HOME | chain-id=$CHAINID | p2p=:$VAL2_P2P_PORT | rpc=:$VAL2_RPC_PORT | profiling=:$VAL2_PPROF_PORT | grpc=:$VAL2_GRPC_PORT"
trap 'pkill -f '"$BIND" EXIT
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo "««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««««"
echo "performing sanity tests"
rm -rf $VAL1HOME $VAL2HOME 
git clone https://github.com/permissionlessweb/go-bitsong
cd go-bitsong
git checkout main
make install 
cd ../ &&

rm -rf $VAL1HOME/test-keys
rm -rf $VAL2HOME/test-key

$BIND init $CHAINID --overwrite --home $VAL1HOME --chain-id $CHAINID &&
$BIND init $CHAINID --overwrite --home $VAL2HOME --chain-id $CHAINID

mkdir $VAL1HOME/test-keys
mkdir $VAL2HOME/test-keys

$BIND --home $VAL1HOME config keyring-backend test &&
$BIND --home $VAL2HOME config keyring-backend test &&
$BIND --home $VAL2HOME config node tcp://localhost:$VAL2_RPC_PORT &&

# remove val2 genesis
rm -rf $VAL2HOME/config/genesis.json &&
# modify val1 genesis 
jq ".app_state.crisis.constant_fee.denom = \"ubtsg\" |
      .app_state.staking.params.bond_denom = \"ubtsg\" |
      .app_state.mint.params.blocks_per_year = \"20000000\" |
      .app_state.mint.params.mint_denom = \"ubtsg\" |
      .app_state.merkledrop.params.creation_fee.denom = \"ubtsg\" |
      .app_state.gov.voting_params.voting_period = \"15s\" |
      .app_state.gov.params.expedited_voting_period = \"5s\" |
      .app_state.gov.params.voting_period = \"15s\" |
      .app_state.gov.params.min_deposit[0].denom = \"ubtsg\" |
      .app_state.fantoken.params.burn_fee.denom = \"ubtsg\" |
      .app_state.fantoken.params.issue_fee.denom = \"ubtsg\" |
      .app_state.slashing.params.signed_blocks_window = \"10\" |
      .app_state.slashing.params.min_signed_per_window = \"1.000000000000000000\" |
      .app_state.fantoken.params.mint_fee.denom = \"ubtsg\"" $VAL1HOME/config/genesis.json > $VAL1HOME/config/tmp.json
# give val2 a genesis
mv $VAL1HOME/config/tmp.json $VAL1HOME/config/genesis.json

# setup test keys.
yes | $BIND  --home $VAL1HOME keys add validator1  --output json > $VAL1HOME/test-keys/val.json 2>&1 &&
yes | $BIND --home $VAL2HOME keys add validator2  --output json > $VAL2HOME/test-keys/val.json 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add user    --output json > $VAL1HOME/test-keys/key_seed.json 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add relayer --output json > $VAL2HOME/test-keys/relayer_seed.json 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add delegator1 --output json > $VAL1HOME/test-keys/del.json 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add delegator2  --output json > $VAL2HOME/test-keys/del.json 2>&1 &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL1HOME keys show user -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL2HOME keys show relayer -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL1HOME keys show validator1 -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL2HOME keys show validator2 -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL1HOME keys show delegator1 -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $($BIND --home $VAL2HOME keys show delegator2 -a) $defaultCoins &&
$BIND --home $VAL1HOME genesis gentx validator1 $delegate --chain-id $CHAINID  &&
$BIND genesis collect-gentxs --home $VAL1HOME
sleep 1

cp $VAL1HOME/config/genesis.json $VAL2HOME/config/genesis.json
VAL1_P2P_ADDR=$($BIND tendermint show-node-id --home $VAL1HOME)@localhost:$VAL1_P2P_PORT


# keys 
DEL1=$(jq -r '.name' $CHAINDIR/$CHAINID/val1/test-keys/del.json)
DEL1ADDR=$(jq -r '.address' $CHAINDIR/$CHAINID/val1/test-keys/del.json)
DEL2=$(jq -r '.name'  $CHAINDIR/$CHAINID/val2/test-keys/del.json)
DEL2ADDR=$(jq -r '.address'  $CHAINDIR/$CHAINID/val2/test-keys/del.json)
VAL1=$(jq -r '.name' $CHAINDIR/$CHAINID/val1/test-keys/val.json)
VAL1ADDR=$(jq -r '.address'  $CHAINDIR/$CHAINID/val1/test-keys/val.json)
VAL2=$(jq -r '.name'  $CHAINDIR/$CHAINID/val2/test-keys/val.json)
VAL2ADDR=$(jq -r '.address'  $CHAINDIR/$CHAINID/val2/test-keys/val.json)
# Get validator's public key (ensure jq is installed)
PUBKEY_KEY=$($BIND tendermint show-validator --home $VAL2HOME | jq -r '.key')

# app & config modiifications
# config.toml
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL1_PROXY_APP_PORT\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL1_P2P_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"2s\"/" "$VAL1HOME/config/config.toml"
 
# val2
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL2_PROXY_APP_PORT\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL2_P2P_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$VAL1_P2P_ADDR\"/" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"2s\"/" "$VAL1HOME/config/config.toml"

# app.toml
sed -i.bak "/^\[api\]/,/^\[/ s/minimum-gas-prices.*/minimum-gas-prices = \"0.0ubtsg\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[api\]/,/^\[/ s/address.*/address = \"tcp:\/\/0.0.0.0:$VAL1_API_PORT\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[grpc\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_PORT\"/" $VAL1HOME/config/app.toml &&
sed -i.bak "/^\[grpc-web\]/,/^\[/ s/address.*/address = \"localhost:$VAL1_GRPC_WEB_PORT\"/" $VAL1HOME/config/app.toml &&
# val2
sed -i.bak "/^\[api\]/,/^\[/ s/minimum-gas-prices.*/minimum-gas-prices = \"0.0ubtsg\"/" $VAL2HOME/config/app.toml &&
sed -i.bak "/^\[api\]/,/^\[/ s/address.*/address = \"tcp:\/\/0.0.0.0:$VAL2_API_PORT\"/" $VAL2HOME/config/app.toml &&
sed -i.bak "/^\[grpc\]/,/^\[/ s/address.*/address = \"localhost:$VAL2_GRPC_PORT\"/" $VAL2HOME/config/app.toml &&
sed -i.bak "/^\[grpc-web\]/,/^\[/ s/address.*/address = \"localhost:$VAL2_GRPC_WEB_PORT\"/" $VAL2HOME/config/app.toml &&

# Start bitsong
echo "Starting Genesis validator..."
$BIND start --home $VAL1HOME & 
VAL1_PID=$!
echo "VAL1_PID: $VAL1_PID"
sleep 4
####################################################################
# B. SLASH
####################################################################
bitsongd start --home $VAL2HOME &
VAL2_PID=$!
echo "VAL2_PID: $VAL2_PID"
sleep 4
#!/bin/bash
# Create JSON file in the validator's home directory
cat <<EOF > "$VAL2HOME/validator.json"
{
  "pubkey": {
    "@type": "/cosmos.crypto.ed25519.PubKey",
    "key": "$PUBKEY_KEY"
  },
  "amount": "9000000000ubtsg",
  "moniker": "VAL2",
  "identity": "",
  "website": "",
  "security": "",
  "details": "",
  "commission-rate": "0.10",
  "commission-max-rate": "0.20",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF
echo "Validator JSON created at $VAL2HOME/validator.json"
bitsongd tx staking create-validator $VAL2HOME/validator.json --gas auto --gas-adjustment 1.5 --fees="600ubtsg"  --chain-id=$CHAINID --home $VAL2HOME --from=$VAL2 -y
sleep 3
VAL1_OP_ADDR=$(jq -r '.body.messages[0].validator_address' $VAL1HOME/config/gentx/gentx-*.json)
VAL2_OP_ADDR=$($BIND q staking validators --home $VAL2HOME -o json | jq -r ".validators[] | select(.operator_address!= \"$VAL1_OP_ADDR\") |.operator_address" | head -1)
$BIND tx staking delegate $VAL1_OP_ADDR 99000000ubtsg --from $DEL1 --gas auto  --fees 200ubtsg --gas-adjustment 1.2 --chain-id $CHAINID --home $VAL1HOME -y 
$BIND tx staking delegate $VAL2_OP_ADDR 400000000ubtsg --from $DEL2 --gas auto --fees 800ubtsg --gas-adjustment 1.4 --chain-id $CHAINID --home $VAL2HOME -y
sleep 3
$BIND tx staking delegate $VAL2_OP_ADDR 99000000ubtsg --from $DEL1 --gas auto  --fees 800ubtsg --gas-adjustment 1.2 --chain-id $CHAINID --home $VAL1HOME -y 
kill $VAL1_PID
sleep 12
$BIND start --home $VAL1HOME &
sleep 10
DEL1_QUERY=$($BIND q staking delegation $DEL1ADDR $VAL1_OP_ADDR --home $VAL1HOME -o json)
DEL2_QUERY=$($BIND q staking delegation $DEL2ADDR $VAL2_OP_ADDR --home $VAL2HOME -o json)
DEL1_REWARD=$($BIND q distribution rewards $DEL1ADDR --home $VAL1HOME --output json)
DEL2_REWARD=$($BIND q distribution rewards $DEL2ADDR --home $VAL1HOME --output json)
VAL1_DEL1_SHARES=$(echo "$DEL1_QUERY" | jq -r '.delegation.shares')
VAL1_DEL1_BTSG=$(echo "$DEL1_QUERY" | jq -r '.balance.amount')
VAL2_DEL2_SHARES=$(echo "$DEL2_QUERY" | jq -r '.delegation.shares')
VAL2_DEL2_BTSG=$(echo "$DEL2_QUERY" | jq -r '.balance.amount')
VAL1_OUTSTANDING_REWARDS=$($BIND q distribution validator-outstanding-rewards $VAL1_OP_ADDR --home $VAL1HOME -o json | jq -r '.rewards[] | select(.denom == "ubtsg") | .amount')
VAL1_TOTAL_SHARES=$($BIND q staking validator $VAL1_OP_ADDR --home $VAL1HOME -o json | jq -r '.delegator_shares')
VAL1_TOTAL_TOKENS=$($BIND q staking validator $VAL1_OP_ADDR --home $VAL1HOME -o json | jq -r '.tokens')
VAL2_OUTSTANDING_REWARDS=$($BIND q distribution validator-outstanding-rewards $VAL2_OP_ADDR --home $VAL1HOME -o json | jq -r '.rewards[] | select(.denom == "ubtsg") | .amount')
VAL2_TOTAL_SHARES=$($BIND q staking validator $VAL2_OP_ADDR --home $VAL1HOME -o json | jq -r '.delegator_shares')
VAL2_TOTAL_TOKENS=$($BIND q staking validator $VAL2_OP_ADDR --home $VAL1HOME -o json | jq -r '.tokens')
DEL1_PRE_UPGR_BALANCE=$($BIND q bank balances $DEL1ADDR --home $VAL2HOME --output json | jq -r '.balances[] | select(.denom == "ubtsg") | .amount')
DEL2_PRE_UPGR_BALANCE=$($BIND q bank balances $DEL2ADDR --home $VAL2HOME --output json | jq -r '.balances[] | select(.denom == "ubtsg") | .amount')
VAL1_PRE_UPGR_BALANCE=$($BIND q bank balances $VAL1ADDR --home $VAL1HOME --output json | jq -r '.balances[] | select(.denom == "ubtsg") | .amount')
VAL2_PRE_UPGR_BALANCE=$($BIND q bank balances $VAL2ADDR --home $VAL1HOME --output json | jq -r '.balances[] | select(.denom == "ubtsg") | .amount')
echo "DEL1_REWARD: $DEL1_REWARD"
echo "DEL2_REWARD: $DEL2_REWARD"
echo "VAL1_OP_ADDR: $VAL1_OP_ADDR"
echo "VAL2_OP_ADDR: $VAL2_OP_ADDR"
echo "DEL1ADDR: $DEL1ADDR"
echo "DEL2ADDR: $DEL2ADDR"
echo "DEL1_QUERY: $DEL1_QUERY"
echo "DEL2_QUERY: $DEL2_QUERY"
echo "VAL1_DEL1_SHARES: $VAL1_DEL1_SHARES"
echo "VAL1_DEL1_BTSG: $VAL1_DEL1_BTSG"
echo "VAL2_DEL2_SHARES: $VAL2_DEL2_SHARES"
echo "VAL2_DEL2_BTSG: $VAL2_DEL2_BTSG"
echo "VAL1_OUTSTANDING_REWARDS:$VAL1_OUTSTANDING_REWARDS"
echo "VAL1_TOTAL_SHARES:$VAL1_TOTAL_SHARES"
echo "VAL1_TOTAL_TOKENS:$VAL1_TOTAL_TOKENS"
echo "VAL2_OUTSTANDING_REWARDS:$VAL2_OUTSTANDING_REWARDS"
echo "VAL2_TOTAL_SHARES:$VAL2_TOTAL_SHARES"
echo "VAL2_TOTAL_TOKENS:$VAL2_TOTAL_TOKENS"
echo "DEL1_PRE_UPGR_BALANCE:$DEL1_PRE_UPGR_BALANCE"
echo "DEL2_PRE_UPGR_BALANCE:$DEL2_PRE_UPGR_BALANCE"
echo "VAL1_PRE_UPGR_BALANCE:$VAL1_PRE_UPGR_BALANCE"
echo "VAL2_PRE_UPGR_BALANCE:$VAL2_PRE_UPGR_BALANCE"
echo "VAL2_PRE_UPGR_BALANCE:$VAL2_PRE_UPGR_BALANCE"
sleep 1
if [ -z "$VAL1_DEL1_SHARES" ] || [ -z "$VAL1_DEL1_BTSG" ] || [ -z "$VAL2_DEL2_SHARES" ] || [ -z "$VAL2_DEL2_BTSG" ]; then
  echo "Error: unable to extract delegation information."
  exit 1
fi
sleep 1
####################################################################
# C. CONFIRM SLASH BREAKS WITH BAD APP VERSION:
# we make use of a bitsong app version we know is broken.
####################################################################
echo "performing insanity tests"
pkill -f bitsongd
cd go-bitsong || exit
git checkout v0.24.BROKEN
make install 
cd ../
bitsongd start --home $VAL2HOME &
VAL2_PID=$!
echo "VAL2_PID: $VAL2_PID"
bitsongd start --home $VAL1HOME &
VAL1_PID=$!
echo "VAL1_PID: $VAL1_PID"
sleep 6
kill $VAL1_PID
sleep 12
$BIND start --home $VAL1HOME &
sleep 6
### CONFIRM HOOKS ARE BROKEN
DEL1_REWARD=$($BIND q distribution rewards $DEL1ADDR --home $VAL1HOME --output json)
DEL2_REWARD=$($BIND q distribution rewards $DEL2ADDR --home $VAL1HOME --output json)
echo "DEL1_REWARD: $DEL1_REWARD"
echo "DEL2_REWARD: $DEL2_REWARD"
pkill -f bitsongd
# bitsongd export --for-zero-height --home $VAL1HOME 