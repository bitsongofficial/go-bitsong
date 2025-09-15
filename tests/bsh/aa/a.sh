#!/bin/bash
BIND=bitsongd
CHAINID_A=test-1
CHAINID_B=test-2

# Authz feature flag - can be overridden with --enable-authz or --disable-authz
USE_AUTHZ=true

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --enable-authz)
      USE_AUTHZ=true
      shift
      ;;
    --disable-authz)
      USE_AUTHZ=false
      shift
      ;;
    --help)
      echo "Usage: $0 [--enable-authz|--disable-authz] [--help]"
      echo ""
      echo "Options:"
      echo "  --enable-authz    Enable authz grants and usage (default)"
      echo "  --disable-authz   Disable authz grants and usage"
      echo "  --help            Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

echo "AuthZ feature: $([ "$USE_AUTHZ" = "true" ] && echo "ENABLED" || echo "DISABLED")"

# setup test keys.
VAL=val
RELAYER=relayer
DEL=del
USER=user
VALFILE="test-keys/$VAL.json"
RELAYERFILE="test-keys/$RELAYER.json"
USERFILE="test-keys/$USER.json"

# file paths
CHAINDIR=./data
VAL1HOME=$CHAINDIR/$CHAINID_A/val1
VAL2HOME=$CHAINDIR/$CHAINID_B/val1
HERMES=~/.hermes
# abstract paths
ABSTRACT_DIR="./abstract"
ARTIFACTS_DIR="$ABSTRACT_DIR/framework/artifacts"
SCRIPTS_DIR="./ibaa-scripts"
 
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
# trap 'pkill -f '"$BIND" EXIT


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
$BIND init $CHAINID_B --overwrite --home $VAL2HOME --chain-id $CHAINID_B &&

mkdir $VAL1HOME/test-keys
mkdir $VAL2HOME/test-keys

# cli config
$BIND --home $VAL1HOME config keyring-backend test
$BIND --home $VAL2HOME config keyring-backend test
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
# set correct chain-id for second network

jq ".chain_id = \"$CHAINID_B\"" $VAL2HOME/config/genesis.json > $VAL2HOME/config/tmp.json
mv $VAL2HOME/config/tmp.json $VAL2HOME/config/genesis.json

yes | $BIND  --home $VAL1HOME keys add $VAL --output json > $VAL1HOME/$VALFILE 2>&1 &&
yes | $BIND  --home $VAL2HOME keys add $VAL --output json > $VAL2HOME/$VALFILE 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add $USER --output json > $VAL1HOME/$USERFILE 2>&1 &&
yes | $BIND  --home $VAL1HOME keys add $RELAYER  --output json > $VAL1HOME/$RELAYERFILE 2>&1 &&
RELAYERADDR=$(jq -r '.address' $VAL1HOME/$RELAYERFILE)
VAL1A_ADDR=$(jq -r '.address'  $VAL1HOME/$VALFILE)
VAL1B_ADDR=$(jq -r '.address'  $VAL2HOME/$VALFILE)
USERAADDR=$(jq -r '.address' $VAL1HOME/$USERFILE)


echo "RELAYERADDR: $RELAYERADDR"
echo "VAL1A_ADDR: $VAL1A_ADDR"
echo "VAL1B_ADDR: $VAL1B_ADDR"
echo "USERAADDR: $USERAADDR"

$BIND --home $VAL1HOME genesis add-genesis-account "$USERAADDR" $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account "$RELAYERADDR" $defaultCoins &&
$BIND --home $VAL1HOME genesis add-genesis-account "$VAL1A_ADDR" $defaultCoins &&
$BIND genesis gentx val $delegate --chain-id $CHAINID_A --home $VAL1HOME &&
$BIND genesis collect-gentxs --home $VAL1HOME
sleep 1

# setup second chain 
$BIND genesis add-genesis-account "$VAL1B_ADDR" $defaultCoins --home $VAL2HOME &&
$BIND genesis add-genesis-account "$VAL1A_ADDR" $defaultCoins --home $VAL2HOME &&
$BIND genesis add-genesis-account "$RELAYERADDR" $defaultCoins --home $VAL2HOME &&
$BIND genesis gentx $VAL $defaultCoins --home $VAL2HOME --chain-id $CHAINID_B &&
$BIND genesis collect-gentxs --home $VAL2HOME 

# app & config modiifications
# config.toml
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL1_PROXY_APP_PORT\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL1_RPC_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL1_P2P_PORT\"/" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL1HOME/config/config.toml &&
sed -i.bak -e "s/^pprof_laddr *=.*/pprof_laddr = \"localhost:6060\"/g" $VAL1HOME/config/config.toml &&
# val2
sed -i.bak -e "s/^proxy_app *=.*/proxy_app = \"tcp:\/\/127.0.0.1:$VAL2_PROXY_APP_PORT\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[rpc\]/,/^\[/ s/address.*/address = \"tcp:\/\/127.0.0.1:$VAL2_RPC_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak "/^\[p2p\]/,/^\[/ s/laddr.*/laddr = \"tcp:\/\/0.0.0.0:$VAL2_P2P_PORT\"/" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^grpc_laddr *=.*/grpc_laddr = \"\"/g" $VAL2HOME/config/config.toml &&
sed -i.bak -e "s/^pprof_laddr *=.*/pprof_laddr = \"localhost:6070\"/g" $VAL2HOME/config/config.toml &&
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
cp ../pfm/hermes.toml $HERMES/config.toml

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

## create authz deployment 
 
#####################################################################
# A. DEPLOY BITSONG ACCOUNTS 
####################################################################

# TODO: deploy the bitsong account ownership tokens to test using tokenized ownership-

#####################################################################
# B. DEPLOY ABSTRACT ON BOTH CHAINS
####################################################################

## check if abstract folder exists (./abstract), download if not
if [ ! -d "$ABSTRACT_DIR" ]; then
    echo "Downloading abstract repository..."
    git clone https://github.com/AbstractSDK/abstract.git "$ABSTRACT_DIR"
else
    echo "Abstract repository already exists."
fi
## if exists, check if artifacts exists for the framework (./abstract/framework/artifacts), if doesnt exist, compile wasm (just wasm-all)
if [ ! -d "$ARTIFACTS_DIR" ]; then
    echo "Building wasm artifacts..."
    cd "$ABSTRACT_DIR/framework" || exit
    just wasm
    cd - || exit
  else
    echo "Wasm artifacts already exist."
  fi

## configure deploy scripts in ./ibaa-scripts (.env file w/ test mnemonic, and state file path )
## add val1 mnemonic for use to deploy ($VAL1HOME/mnemonic.txt)
env_file="$SCRIPTS_DIR/.env"
mnemonic=$(jq -r '.mnemonic' "$VAL1HOME/test-keys/val.json")

# Create the .env
rm -rf $env_file
 
cat > "$env_file" <<EOF
LOCAL_MNEMONIC="$mnemonic"
STATE_FILE=./state.json
ARTIFACTS_DIR=../abstract/framework/artifacts
LOGGING=debug
CW_ORCH_SERIALIZE_JSON=true
USE_AUTHZ=$USE_AUTHZ
EOF
 
 

## grant authz from val1 to user for (upload,init,execute,migrate)
if [ "$USE_AUTHZ" = "true" ]; then
  echo "Setting up AuthZ grants..."
  $BIND tx authz grant $VAL1A_ADDR  generic --msg-type=/cosmwasm.wasm.v1.MsgExecuteContract --from $USERAADDR --fees 1000ubtsg --chain-id $CHAINID_A --home $VAL1HOME -y
  sleep 3
  $BIND tx authz grant $VAL1A_ADDR  generic --msg-type=/cosmwasm.wasm.v1.MsgMigrateContract --from $USERAADDR --fees 1000ubtsg --chain-id $CHAINID_A --home $VAL1HOME -y
  sleep 3
  $BIND tx authz grant $VAL1A_ADDR  generic --msg-type=/cosmwasm.wasm.v1.MsgStoreCode --from $USERAADDR --fees 1000ubtsg --chain-id $CHAINID_A --home $VAL1HOME -y
  sleep 3
  $BIND tx authz grant $VAL1A_ADDR  generic --msg-type=/cosmwasm.wasm.v1.MsgInstantiateContract --from $USERAADDR --fees 1000ubtsg --chain-id $CHAINID_A --home $VAL1HOME -y
  sleep 3
else
  echo "Skipping AuthZ grants (disabled)"
fi
 


# # Deploy on both chains
echo "Deploying on both chains..."
cd "$SCRIPTS_DIR" || exit  
cat .env
if [ "$USE_AUTHZ" = "true" ]; then
  echo "Running with AuthZ granter: $USERAADDR"
  RUST_LOG=info cargo run --bin full_deploy -- --authz-granter "$USERAADDR"
else
  echo "Running without AuthZ"
  RUST_LOG=info cargo run --bin full_deploy
fi

echo "Preparation and deployment complete."


#####################################################################
# C. CREATE INTERCHAIN ABSTRACT ACCOUNT 
####################################################################


####################################################################
# D. INSTALL MODULE ON ACCOUNT
####################################################################

####################################################################
# E. USE MODULE
####################################################################