#!/bin/bash

KEYRING=--keyring-backend="test"

BINARY=$1
CHAINID=$2
CHAINDIR=$3
RPCPORT=$4
P2PPORT=$5
PROFPORT=$6
GRPCPORT=$7
DENOM=$8

echo "Creating $BINARY instance: home=$CHAINDIR | chain-id=$CHAINID | p2p=:$P2PPORT | rpc=:$RPCPORT | profiling=:$PROFPORT | grpc=:$GRPCPORT"

if ! mkdir -p $CHAINDIR/$CHAINID 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

coins="100000000000$DENOM"
delegate="100000000000$DENOM"

$BINARY --home $CHAINDIR/$CHAINID --chain-id $CHAINID init $CHAINID
sleep 1

jq ".app_state.crisis.constant_fee.denom = \"$DENOM\" |
      .app_state.staking.params.bond_denom = \"$DENOM\" |
      .app_state.merkledrop.params.creation_fee.denom = \"$DENOM\" |
      .app_state.gov.deposit_params.min_deposit[0].denom = \"$DENOM\" |
      .app_state.fantoken.params.burn_fee.denom = \"$DENOM\" |
      .app_state.fantoken.params.issue_fee.denom = \"$DENOM\" |
      .app_state.fantoken.params.mint_fee.denom = \"$DENOM\"" $CHAINDIR/$CHAINID/config/genesis.json > tmp.json

mv tmp.json $CHAINDIR/$CHAINID/config/genesis.json

$BINARY --home $CHAINDIR/$CHAINID keys add validator $KEYRING --output json > $CHAINDIR/$CHAINID/validator_seed.json 2>&1 &&
$BINARY --home $CHAINDIR/$CHAINID keys add user $KEYRING --output json > $CHAINDIR/$CHAINID/key_seed.json 2>&1 &&
$BINARY --home $CHAINDIR/$CHAINID keys add relayer $KEYRING --output json > $CHAINDIR/$CHAINID/relayer_seed.json 2>&1 &&
$BINARY --home $CHAINDIR/$CHAINID genesis add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID keys $KEYRING show user -a) $coins &&
$BINARY --home $CHAINDIR/$CHAINID genesis add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID keys $KEYRING show validator -a) $coins &&
$BINARY --home $CHAINDIR/$CHAINID genesis add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID keys $KEYRING show relayer -a) $coins &&
$BINARY --home $CHAINDIR/$CHAINID genesis gentx validator $delegate $KEYRING --chain-id $CHAINID &&
$BINARY --home $CHAINDIR/$CHAINID genesis collect-gentxs
sleep 1

echo "Change settings in config.toml and genesis.json files..."
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
sed -i 's#"localhost:6060"#"localhost:'"$PROFPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/$CHAINID/config/config.toml
#sed -i 's/enable = false/enable = true/g' $CHAINDIR/$CHAINID/config/app.toml
#sed -i 's/swagger = false/swagger = true/g' $CHAINDIR/$CHAINID/config/app.toml
sed -i 's/"voting_period": "172800s"/"voting_period": "120s"/g' $CHAINDIR/$CHAINID/config/genesis.json
sed -i 's/"stake"/"ubtsg"/g' $CHAINDIR/$CHAINID/config/genesis.json

echo "Starting $CHAINID in $CHAINDIR..."
echo "Log file is located at $CHAINDIR/$CHAINID.log"

$BINARY --home $CHAINDIR/$CHAINID start --pruning=nothing --grpc-web.enable=false --grpc.address="0.0.0.0:$GRPCPORT"
