#!/bin/sh

CHAIN_ID=localbitsong
BITSONG_HOME=$HOME/.bitsongd
CONFIG_FOLDER=$BITSONG_HOME/config
MONIKER=val
STATE='false'

MNEMONIC="bottom loan skill merry east cradle onion journey palm apology verb edit desert impose absurd oil bubble sweet glove shallow size build burst effort"

while getopts s flag
do
    case "${flag}" in
        s) STATE='true';;
    esac
done

install_prerequisites () {
    apk add dasel lz4
}

edit_genesis () {

    GENESIS=$CONFIG_FOLDER/genesis.json

    # Update staking module
    dasel put string -f $GENESIS '.app_state.staking.params.bond_denom' -v 'ubtsg'
    dasel put string -f $GENESIS '.app_state.staking.params.unbonding_time' -v '240s'

    # Update bank module
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].description' -v 'Registered denom ubtsg for localbitsong testing'
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].denom_units.[].denom' -v 'ubtsg'
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].denom_units.[0].exponent' -v 0
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].base' -v 'ubtsg'
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].display' -v 'ubtsg'
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].name' -v 'ubtsg'
    dasel put string -f $GENESIS '.app_state.bank.denom_metadata.[0].symbol' -v 'ubtsg'

    # Update crisis module
    dasel put string -f $GENESIS '.app_state.crisis.constant_fee.denom' -v 'ubtsg'

    # Update gov module
    dasel put string -f $GENESIS '.app_state.gov.voting_params.voting_period' -v '60s'
    dasel put string -f $GENESIS '.app_state.gov.params.min_deposit.[0].denom' -v 'ubtsg'

    # Update mint module
    dasel put string -f $GENESIS '.app_state.mint.params.mint_denom' -v "ubtsg"
}

add_genesis_accounts () {

    # val
    bitsongd genesis add-genesis-account bitsong1gws6wz8q5kyyu4gqze48fwlmm4m0mdjz0620gw 10000000000000ubtsg --home $BITSONG_HOME
    
    # wallets
    bitsongd genesis add-genesis-account bitsong1regz7kj3ylg2dn9rl8vwrhclkgz528mf0tfsck 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong1hvrhhex6wfxh7r7nnc3y39p0qlmff6v9t5rc25 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong175vgzztymvvcxvqun54nlu9dq6856thgvyl5sa 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong1t8nznzj4sd6zzutwdmslgy4dcxyd2jafz7822x 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong14vdrvstsffj8mq5e4fhm6y2hpfxtedajczsj5d 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong1vwe5hay74v0vhuzdhadteyqfasu5d7tdf83pyy 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong16866dezn6ez2qpmpcrrv9cyud8v8c7ufnzwhhh 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong1tlwh75lvu35nw9vcg557mxhspz5s88t6vzscd8 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong16z9wj8n5f3zgzwspw0r9sj9v7k7hdasqj95us9 10000000000000ubtsg --home $BITSONG_HOME
    bitsongd genesis add-genesis-account bitsong1gulaxnca7rped0grw0lz4h4zy0xn3ttvmlad8x 10000000000000ubtsg --home $BITSONG_HOME

    echo $MNEMONIC | bitsongd keys add $MONIKER --recover --keyring-backend=test --home $BITSONG_HOME
    bitsongd genesis gentx $MONIKER 5000000000000ubtsg --keyring-backend=test --chain-id=$CHAIN_ID --home $BITSONG_HOME

    bitsongd genesis collect-gentxs --home $BITSONG_HOME

    bitsongd genesis validate-genesis --home $BITSONG_HOME
}

edit_config () {

    # Remove seeds
    dasel put string -f $CONFIG_FOLDER/config.toml '.p2p.seeds' -v ''

    # Expose the rpc
    dasel put string -f $CONFIG_FOLDER/config.toml '.rpc.laddr' -v "tcp://0.0.0.0:26657"
    
    # Expose pprof for debugging
    # To make the change enabled locally, make sure to add 'EXPOSE 6060' to the root Dockerfile
    # and rebuild the image.
    dasel put string -f $CONFIG_FOLDER/config.toml '.rpc.pprof_laddr' -v "0.0.0.0:6060"
}

enable_cors () {

    # Enable cors on RPC
    dasel put string -f $CONFIG_FOLDER/config.toml -v "*" '.rpc.cors_allowed_origins.[]'
    dasel put string -f $CONFIG_FOLDER/config.toml -v "Accept-Encoding" '.rpc.cors_allowed_headers.[]'
    dasel put string -f $CONFIG_FOLDER/config.toml -v "DELETE" '.rpc.cors_allowed_methods.[]'
    dasel put string -f $CONFIG_FOLDER/config.toml -v "OPTIONS" '.rpc.cors_allowed_methods.[]'
    dasel put string -f $CONFIG_FOLDER/config.toml -v "PATCH" '.rpc.cors_allowed_methods.[]'
    dasel put string -f $CONFIG_FOLDER/config.toml -v "PUT" '.rpc.cors_allowed_methods.[]'

    # Enable unsafe cors and swagger on the api
    dasel put bool -f $CONFIG_FOLDER/app.toml -v "true" '.api.swagger'
    dasel put bool -f $CONFIG_FOLDER/app.toml -v "true" '.api.enabled-unsafe-cors'

    # Enable cors on gRPC Web
    dasel put bool -f $CONFIG_FOLDER/app.toml -v "true" '.grpc-web.enable-unsafe-cors'

    # Enable SQS & route caching
    dasel put string -f $CONFIG_FOLDER/app.toml -v "true" '.bitsong-sqs.is-enabled'
    dasel put string -f $CONFIG_FOLDER/app.toml -v "true" '.bitsong-sqs.route-cache-enabled'

    dasel put string -f $CONFIG_FOLDER/app.toml -v "redis" '.bitsong-sqs.db-host'
}

run_with_retries() {
  cmd=$1
  success_msg=$2

  substring='code: 0'
  COUNTER=0

  while [ $COUNTER -lt 15 ]; do
    string=$(eval $cmd 2>&1)
    echo $string

    if [ "$string" != "${string%"$substring"*}" ]; then
      echo "$success_msg"
      break
    else
      COUNTER=$((COUNTER+1))
      sleep 0.5
    fi
  done
}

if [[ ! -d $CONFIG_FOLDER ]]
then
    echo $MNEMONIC | bitsongd init -o --chain-id=$CHAIN_ID --home $BITSONG_HOME --recover $MONIKER
    install_prerequisites
    edit_genesis
    add_genesis_accounts
    edit_config
    enable_cors
fi

bitsongd start --home $BITSONG_HOME &

wait
