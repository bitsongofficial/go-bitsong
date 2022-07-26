#!/bin/sh

rm -rf ~/.bitsongd

set -o errexit -o nounset

# Build genesis file incl account for passed address
bitsongd init --chain-id test test
bitsongd keys add validator --keyring-backend="test"
bitsongd add-genesis-account $(bitsongd keys show validator -a --keyring-backend="test") 100000000000000stake
bitsongd gentx validator 100000000stake --keyring-backend="test" --chain-id test
bitsongd collect-gentxs

# Set proper defaults and change ports
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bitsongd/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bitsongd/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bitsongd/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' ~/.bitsongd/config/config.toml

# Start bitsong
bitsongd start --pruning=nothing
