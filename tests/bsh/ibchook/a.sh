#!/bin/bash
BIND=bitsongd
CHAINID_A=test-1
CHAINID_B=test-2

# setup test keys.
VAL=val
RELAYER=relayer
DEL=del
USER=user
DELFILE="test-keys/$DEL.json"
VALFILE="test-keys/$VAL.json"
RELAYERFILE="test-keys/$RELAYER.json"
USERFILE="test-keys/$USER.json"

# file paths
CHAINDIR=./data
VAL1HOME=$CHAINDIR/$CHAINID_A/val1
VAL2HOME=$CHAINDIR/$CHAINID_B/val1
HERMES=~/.hermes
HERMES_CFG_TEMPLATE_PATH="helpers/relayer/hermes.toml"
 

# Define the new ports for val1 on chain a
VAL1_API_PORT=1317
VAL1_GRPC_PORT=9090
VAL1_GRPC_WEB_PORT=9091
VAL1_PROXY_APP_PORT=26658
VAL1_RPC_PORT=26657
VAL1_PPROF_PORT=6060
VAL1_P2P_PORT=26656

# Define the new ports for val1 on chain b
VAL2_API_PORT=1318
VAL2_GRPC_PORT=10090
VAL2_GRPC_WEB_PORT=10091
VAL2_PROXY_APP_PORT=9395
VAL2_RPC_PORT=27657
VAL2_PPROF_PORT=6361
VAL2_P2P_PORT=26356

echo "Creating $BIND instance for VAL1_A: home=$VAL1HOME | chain-id=$CHAINID_A | p2p=:$VAL1_P2P_PORT | rpc=:$VAL1_RPC_PORT | profiling=:$VAL1_PPROF_PORT | grpc=:$VAL1_GRPC_PORT"
echo "Creating $BIND instance for VAL: home=$VAL2HOME | chain-id=$CHAINID_B | p2p=:$VAL2_P2P_PORT | rpc=:$VAL2_RPC_PORT | profiling=:$VAL2_PPROF_PORT | grpc=:$VAL2_GRPC_PORT"
trap 'pkill -f '"$BIND" EXIT

defaultCoins="100000000000ubtsg"  # 100K
delegate="1000000ubtsg" # 1btsg

####################################################################
# A. CHAINS CONFIG
####################################################################

rm -rf $VAL1HOME $VAL2HOME 
rm -rf $VAL1HOME/test-keys
rm -rf $VAL2HOME/test-keys

# initialize chains
$BIND init $CHAINID_A --overwrite --home $VAL1HOME --chain-id $CHAINID_A &&
$BIND init $CHAINID_B --overwrite --home $VAL2HOME --chain-id $CHAINID_B

mkdir $VAL1HOME/test-keys
mkdir $VAL2HOME/test-keys

# cli config
$BIND --home $VAL1HOME config keyring-backend test
$BIND --home $VAL2HOME config keyring-backend test &&
$BIND --home $VAL1HOME config chain-id $CHAINID_A
$BIND --home $VAL2HOME config chain-id $CHAINID_B &&
$BIND --home $VAL1HOME config node tcp://localhost:$VAL1_RPC_PORT
$BIND --home $VAL2HOME config node tcp:\/\/127.0.0.1:$VAL2_RPC_PORT &&

# optimize val1 genesis for testing
jq ".app_state.crisis.constant_fee.denom = \"ubtsg\" |
      .app_state.staking.params.bond_denom = \"ubtsg\" |
      .app_state.mint.params.blocks_per_year = \"20000000\" |
      .app_state.mint.params.mint_denom = \"ubtsg\" |
      .app_state.merkledrop.params.creation_fee.denom = \"ubtsg\" |
      .app_state.gov.voting_params.voting_period = \"15s\" |
      .app_state.gov.voting_params.voting_period = \"15s\" |
      .app_state.gov.params.voting_period = \"15s\" |
      .app_state.gov.params.expedited_voting_period = \"12s\" |
      .app_state.gov.params.min_deposit[0].denom = \"ubtsg\" |
      .app_state.fantoken.params.burn_fee.denom = \"ubtsg\" |
      .app_state.fantoken.params.issue_fee.denom = \"ubtsg\" |
      .app_state.slashing.params.signed_blocks_window = \"15\" |
      .app_state.slashing.params.min_signed_per_window = \"0.500000000000000000\" |
      .app_state.fantoken.params.mint_fee.denom = \"ubtsg\"" $VAL1HOME/config/genesis.json > $VAL1HOME/config/tmp.json
# give val2 genesis optimized genesis
mv $VAL1HOME/config/tmp.json $VAL1HOME/config/genesis.json
cp $VAL1HOME/config/genesis.json $VAL2HOME/config/genesis.json
jq ".chain_id = \"$CHAINID_B\"" $VAL2HOME/config/genesis.json > $VAL2HOME/config/tmp.json
mv $VAL2HOME/config/tmp.json $VAL2HOME/config/genesis.json

yes | $BIND  --home $VAL1HOME keys add $VAL --output json > $VAL1HOME/$VALFILE 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add $VAL --output json > $VAL2HOME/$VALFILE 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add $USER --output json > $VAL1HOME/$USERFILE 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add $USER --output json > $VAL2HOME/$USERFILE 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add $DEL --output json > $VAL1HOME/$DELFILE 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add $DEL  --output json > $VAL2HOME/$DELFILE 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add $RELAYER  --output json > $VAL1HOME/$RELAYERFILE 2>&1 &&

RELAYERADDR=$(jq -r '.address' $VAL1HOME/$RELAYERFILE)
DEL1ADDR=$(jq -r '.address' $VAL1HOME/$DELFILE)
DEL2ADDR=$(jq -r '.address'  $VAL2HOME/$DELFILE)
VAL1A_ADDR=$(jq -r '.address'  $VAL1HOME/$VALFILE)
VAL1B_ADDR=$(jq -r '.address'  $VAL2HOME/$VALFILE)
USERAADDR=$(jq -r '.address' $VAL1HOME/$USERFILE)
USERBADDR=$(jq -r '.address' $VAL2HOME/$USERFILE)


$BIND --home $VAL1HOME genesis add-genesis-account $USERAADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $RELAYERADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $VAL1A_ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $DEL1ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $DEL2ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis gentx $VAL $delegate --chain-id $CHAINID_A &&
$BIND genesis collect-gentxs --home $VAL1HOME

# setup second chain 
$BIND genesis add-genesis-account $USERBADDR $defaultCoins --home $VAL2HOME && 
$BIND genesis add-genesis-account $VAL1B_ADDR $defaultCoins --home $VAL2HOME &&
$BIND genesis add-genesis-account $RELAYERADDR $defaultCoins --home $VAL2HOME &&
$BIND genesis add-genesis-account $DEL2ADDR $defaultCoins --home $VAL2HOME &&
$BIND genesis gentx $VAL $delegate --home $VAL2HOME --chain-id $CHAINID_B &&
$BIND genesis collect-gentxs --home $VAL2HOME 

# app & config modiifications
# config.toml
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL1_PROXY_APP_PORT\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL1_P2P_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^pprof_laddr *=.*/pprof_laddr = \"localhost:6060\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"1s\"/" "$VAL1HOME/config/config.toml"
# val2
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL2_PROXY_APP_PORT\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL2_P2P_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^pprof_laddr *=.*/pprof_laddr = \"localhost:6070\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"1s\"/" "$VAL1HOME/config/config.toml"
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


# Start chains
echo "Starting chain 1..."
$BIND start --home $VAL1HOME & 
VAL1A_PID=$!
echo "VAL1A_PID: $VAL1A_PID"
echo "Starting chain 2..."
$BIND start --home $VAL2HOME & 
VAL1B_PID=$!
echo "VAL1B_PID: $VAL1B_PID"
sleep 10

echo "RELAYERADDR: $RELAYERADDR"
echo "DEL1ADDR: $DEL1ADDR"
echo "DEL2ADDR: $DEL2ADDR"
echo "VAL1A_ADDR: $VAL1A_ADDR"
echo "VAL1B_ADDR: $VAL1B_ADDR"
echo "USERAADDR: $USERAADDR"
echo "USERBADDR: $USERBADDR"

## if polytone wasm files dont exist in  ./bin, download 
echo "Uploading ibc_hooks_counter WASM file on the destination chain…"
$BIND tx wasm upload --home $VAL2HOME ../../ict/contracts/ibchooks_counter.wasm --from $USER --chain-id $CHAINID_B --gas auto --gas-adjustment 1.4 --gas auto --fees 400000ubtsg -y 
sleep 2
echo "Uploaded ibc_hooks_counter WASM file on the destination chain successfully…"

####################################################################
# B. RELAYER CONFIG
####################################################################

## create mnemonic file, grab menmonic from relayer key file, print to new txt file
REL_MNEMONIC=$(jq -r '.mnemonic' $VAL1HOME/$RELAYERFILE)
echo "$REL_MNEMONIC" >  $VAL1HOME/mnemonic.txt
## if hermes command does not exist, install hermes
if ! command -v hermes &> /dev/null
then
    cargo install ibc-relayer-cli --bin hermes --locked
fi

## configure hermes with chain & and b
rm -rf $HERMES && mkdir -p $HERMES
cp ../$HERMES_CFG_TEMPLATE_PATH $HERMES/config.toml

## modify $HERMES_CFG toml with correct values 
sed -i.bak "/^\[chains\]/,/^\[/ { 
    /id = \"$CHAINID_A\"/ { 
        s/rpc_addr.*/rpc_addr = \"http:\/\/127.0.0.1:$VAL1_RPC_PORT\"/; 
        s/grpc_addr.*/grpc_addr = \"http:\/\/127.0.0.1:$VAL1_GRPC_PORT\"/; 
        s/event_source.url.*/event_source.url = \"ws:\/\/127.0.0.1:$VAL1_RPC_PORT\/websocket\"/; 
        s/key_name.*/key_name = \"$VAL\"/; 
    } 
}" "$HERMES/config.toml"

sed -i.bak "/^\[chains\]/,/^\[/ { 
    /id = \"$CHAINID_B\"/ { 
        s/rpc_addr.*/rpc_addr = \"http:\/\/127.0.0.1:$VAL2_RPC_PORT\"/; 
        s/grpc_addr.*/grpc_addr = \"http:\/\/127.0.0.1:$VAL2_GRPC_PORT\"/; 
        s/event_source.url.*/event_source.url = \"ws:\/\/127.0.0.1:$VAL2_RPC_PORT\/websocket\"/; 
        s/key_name.*/key_name = \"$VAL\"/; 
    } 
}" "$HERMES/config.toml"


echo "Clean up hermes"
hermes keys delete --chain "$CHAINID_A" --all
hermes keys delete --chain "$CHAINID_B" --all

echo "import keys"
hermes keys add --key-name $RELAYER --chain $CHAINID_A --hd-path "m/44'/639'/0'/0/0" --mnemonic-file $VAL1HOME/mnemonic.txt
hermes keys add --key-name $RELAYER --chain $CHAINID_B --hd-path "m/44'/639'/0'/0/0" --mnemonic-file $VAL1HOME/mnemonic.txt


## create channel 
echo "starting relayer" 
echo "Creating IBC transfer channel"
hermes create channel \
  --a-chain "$CHAINID_A" \
  --b-chain "$CHAINID_B" \
  --a-port "transfer" \
  --b-port "transfer" \
  --order unordered \
  --new-client-connection \
  --yes


# start relayer
hermes start & 
HERMES_PID=$! 
trap 'pkill '"$HERMES_PID" EXIT

####################################################################
# C. IBC_HOOKS COUNTER CONFIG
####################################################################
$BIND tx wasm i 1 '{"count":0}' --from $DEL --home $VAL2HOME --no-admin --label="listener contract chain1" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y -o json
sleep 4
# get ibc-hook contract from destination chain
IBC_HOOKS_CC=$($BIND q wasm lca 1 --home $VAL2HOME -o json | jq -r .contracts[0])
MEMO_JSON='{"wasm":{"contract":"'$IBC_HOOKS_CC'","msg":{"increment":{}}}}'
echo "IBC_HOOKS_CC:$IBC_HOOKS_CC"
# init ibc-hooks contract
 
# ibc transfer with ibc-hook in memo
$BIND tx ibc-transfer transfer transfer channel-0 "$DEL2ADDR" 100ubtsg \
    --from $DEL --home $VAL1HOME \
    --fees 400000ubtsg --gas auto --gas-adjustment 1.3 \
    --memo "$MEMO_JSON" \
    -y 

# wait for packet to relay.
sleep 30



$BIND tx ibc-transfer transfer transfer channel-0 "$IBC_HOOKS_CC" 100ubtsg \
    --from $DEL --home $VAL1HOME \
    --fees 400000ubtsg --gas auto --gas-adjustment 1.3 \
    --memo "$MEMO_JSON" \
    -y &&

# GetIBCHooksUserAddress
IBC_HOOKS_Q=$($BIND q ibchooks wasm-sender channel-0 "$DEL1ADDR" --home $VAL1HOME --output json)
echo "IBC_HOOKS_Q: $IBC_HOOKS_Q"
sleep 30

# GetIBCHookTotalFunds
TOTAL_FUNDS_Q=$($BIND q wasm contract-state smart "$IBC_HOOKS_CC" "{\"get_total_funds\":{\"addr\":\"$IBC_HOOKS_Q\"}}" --home $VAL2HOME -o json)
HOOK_COUNT=$($BIND q wasm contract-state smart "$IBC_HOOKS_CC" "{\"get_count\":{\"addr\":\"$IBC_HOOKS_Q\"}}" --home $VAL2HOME -o json)
echo "TOTAL_FUNDS_Q: $TOTAL_FUNDS_Q"
echo "HOOK_COUNT: $HOOK_COUNT"
echo "IBC_HOOKS_Q: $IBC_HOOKS_Q"


pkill -f hermes
pkill -f bitsongd
## history should exist, and the callback initiator should equal the test addr
