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
HERMES_CFG_TEMPLATE_PATH="helpers/relayer/hermes.toml"

# file paths
CHAINDIR=../data/polytone
VAL1HOME=$CHAINDIR/$CHAINID_A/val1
VAL2HOME=$CHAINDIR/$CHAINID_B/val1
HERMES=~/.hermes

POLYTONE_CONTRACTS=(
  "polytone_listener.wasm"
  "polytone_note.wasm"
  "polytone_proxy.wasm"
  "polytone_voice.wasm"
  "polytone_tester.wasm"
  )

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
sleep 1

mkdir $VAL1HOME/test-keys
mkdir $VAL2HOME/test-keys

# cli config
$BIND --home $VAL1HOME config keyring-backend test
$BIND --home $VAL2HOME config keyring-backend test &&
$BIND --home $VAL1HOME config chain-id $CHAINID_A &&
$BIND --home $VAL2HOME config chain-id $CHAINID_B &&
$BIND --home $VAL1HOME config node tcp://localhost:$VAL1_RPC_PORT &&
$BIND --home $VAL2HOME config node tcp:\/\/127.0.0.1:$VAL2_RPC_PORT &&
sleep 1

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
yes | $BIND  --home $VAL2HOME keys add user --output json > $VAL2HOME/$USERFILE 2>&1 &&
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


$BIND --home $VAL1HOME genesis add-genesis-account $USERAADDR $defaultCoins
$BIND --home $VAL1HOME genesis add-genesis-account $RELAYERADDR $defaultCoins
$BIND --home $VAL1HOME genesis add-genesis-account $VAL1A_ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $DEL1ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account $DEL2ADDR $defaultCoins &&
$BIND --home $VAL1HOME genesis gentx $VAL $delegate --chain-id $CHAINID_A &&
$BIND genesis collect-gentxs --home $VAL1HOME &&

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
sed -i.bak "/^\[consensus\]/,/^\[/ s/^[[:space:]]*timeout_commit[[:space:]]*=.*/timeout_commit = \"1s\"/" "$VAL2HOME/config/config.toml"
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
if [ -z "$(ls -A ./bin)" ]; then
  sh download.sh
  while [ -z "$(ls -A ./bin)" ]; do
    sleep 1
  done
fi

## upload polytone 
  for contract in "${POLYTONE_CONTRACTS[@]}"; do
    echo "Uploading $contract WASM file..."
    # get tx hash 
    $BIND tx wasm upload --home $VAL2HOME ../../ict/contracts/$contract --from $USER --chain-id $CHAINID_B --gas auto --gas-adjustment 1.4 --gas auto --fees 400000ubtsg -y 
    $BIND tx wasm upload --home $VAL1HOME ../../ict/contracts/$contract --from $DEL  --chain-id $CHAINID_A --gas auto --gas-adjustment 1.4 --gas auto --fees 400000ubtsg -y 
    sleep 2
    echo "Uploaded $contract WASM file successfully."
    sleep 4
done

sleep 1
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

sleep 15

####################################################################
# C. POLYTONE CONFIG
####################################################################
POLYONE_LISTENER_ID=1
POLYONE_NOTE_ID=2
POLYONE_PROXY_ID=3
POLYONE_VOICE_ID=4
POLYONE_TESTER_ID=5

# init note
 $BIND tx wasm i $POLYONE_NOTE_ID '{"block_max_gas": "100000000" }' --from $DEL --home $VAL1HOME --chain-id $CHAINID_A --no-admin --label="note contract chain1" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y 
 $BIND tx wasm i $POLYONE_NOTE_ID '{"block_max_gas": "100000000" }' --from $USER --home $VAL2HOME --chain-id $CHAINID_B --no-admin --label="note contract chain2" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y
sleep 3
# init voice
 $BIND tx wasm i $POLYONE_VOICE_ID '{"proxy_code_id":"3","block_max_gas":"100000000" }' --from $DEL --home $VAL1HOME --chain-id $CHAINID_A --no-admin --label="voice contract chain1" -y --fees 400000ubtsg --gas auto --gas-adjustment 1.3
 $BIND tx wasm i $POLYONE_VOICE_ID '{"proxy_code_id":"3","block_max_gas":"100000000" }' --from $DEL --home $VAL2HOME --chain-id $CHAINID_B --no-admin --label="voice contract chain1" -y --fees 400000ubtsg --gas auto --gas-adjustment 1.3
sleep 5
# init tester
 $BIND tx wasm i $POLYONE_TESTER_ID '{}' --from $DEL --home $VAL1HOME --no-admin --label="tester contract chain1" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y 
 $BIND tx wasm i $POLYONE_TESTER_ID '{}' --from $DEL --home $VAL2HOME --no-admin --label="tester contract chain2" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y 
sleep 3

# get polytone contracts
# POLYONE_PROXY_ADDR_A=$($BIND q wasm lca $POLYONE_PROXY_ID --home $VAL1HOME -o json | jq -r .contracts[0])
# POLYONE_PROXY_ADDR_B=$($BIND q wasm lca $POLYONE_PROXY_ID --home $VAL2HOME -o json | jq -r .contracts[0])
# echo "POLYONE_PROXY_ADDR_A: $POLYONE_PROXY_ADDR_A"
# echo "POLYONE_PROXY_ADDR_B: $POLYONE_PROXY_ADDR_B"
POLYONE_NOTE_ADDR_A=$($BIND q wasm lca $POLYONE_NOTE_ID  --home $VAL1HOME -o json | jq -r .contracts[0])
POLYONE_NOTE_ADDR_B=$($BIND q wasm lca $POLYONE_NOTE_ID  --home $VAL2HOME -o json | jq -r .contracts[0])
POLYONE_TESTER_ADDR_A=$($BIND q wasm lca $POLYONE_TESTER_ID --home $VAL1HOME -o json | jq -r .contracts[0])
POLYONE_TESTER_ADDR_B=$($BIND q wasm lca $POLYONE_TESTER_ID --home $VAL2HOME -o json | jq -r .contracts[0])
POLYONE_VOICE_ADDR_A=$($BIND q wasm lca $POLYONE_VOICE_ID  --home $VAL1HOME -o json | jq -r .contracts[0])
POLYONE_VOICE_ADDR_B=$($BIND q wasm lca $POLYONE_VOICE_ID  --home $VAL2HOME -o json | jq -r .contracts[0])
echo "POLYONE_NOTE_ADDR_A: $POLYONE_NOTE_ADDR_A"
echo "POLYONE_NOTE_ADDR_B: $POLYONE_NOTE_ADDR_B"
echo "POLYONE_TESTER_ADDR_A: $POLYONE_TESTER_ADDR_A"
echo "POLYONE_TESTER_ADDR_B: $POLYONE_TESTER_ADDR_B"
echo "POLYONE_VOICE_ADDR_A: $POLYONE_VOICE_ADDR_A"
echo "POLYONE_VOICE_ADDR_B: $POLYONE_VOICE_ADDR_B"

# init listener
 $BIND tx wasm i $POLYONE_LISTENER_ID "{\"note\":\"$POLYONE_NOTE_ADDR_A\"}" --from $DEL --home $VAL1HOME --no-admin --label="listener contract chain1" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y 
 $BIND tx wasm i $POLYONE_LISTENER_ID "{\"note\":\"$POLYONE_NOTE_ADDR_B\"}" --from $DEL --home $VAL2HOME --no-admin --label="listener contract chain2" --fees 400000ubtsg --gas auto --gas-adjustment 1.3 -y 

## create channel 
echo "starting relayer" 
echo "Creating IBC transfer channel"
hermes create channel --a-chain $CHAINID_A --b-chain $CHAINID_B\
    --a-port "wasm.$POLYONE_NOTE_ADDR_A"\
    --b-port "wasm.$POLYONE_VOICE_ADDR_A"\
    --order unordered\
    --chan-version polytone-1\
    --new-client-connection\
    --yes


# start relayer
hermes start & 
HERMES_PID=$! 
trap 'pkill '"$HERMES_PID" EXIT


####################################################################
# C. POLYTONE INTEGRATION
####################################################################

# send msg to note
$BIND tx wasm e $POLYONE_NOTE_ADDR_A "{\"execute\":{\"msgs\":[],\"timeout_seconds\":\"100\",\"callback\": {\"receiver\": \"$POLYONE_TESTER_ADDR_A\", \"msg\":\"aGVsbG8K\"}}}" --home $VAL1HOME --from $DEL -y --fees 400000ubtsg --gas auto --gas-adjustment 1.3

hermes create channel --a-chain $CHAINID_B --b-chain $CHAINID_A\
    --a-port "wasm.$POLYONE_NOTE_ADDR_B"\
    --b-port "wasm.$POLYONE_VOICE_ADDR_B"\
    --order unordered\
    --chan-version polytone-1\
    --new-client-connection\
    --yes
$BIND tx wasm e $POLYONE_NOTE_ADDR_B "{\"execute\":{\"msgs\":[],\"timeout_seconds\":\"100\",\"callback\": {\"receiver\": \"$POLYONE_TESTER_ADDR_B\", \"msg\":\"aGVsbG8K\"}}}" --home $VAL2HOME --from $DEL -y --fees 400000ubtsg --gas auto --gas-adjustment 1.3


# wait for packet to relay.
sleep 60

# query callback history for test contract 
$BIND q wasm contract-state smart $POLYONE_NOTE_ADDR_A '{"history":{}}' -o json

## history should exist, and the callback initiator should equal the test addr
