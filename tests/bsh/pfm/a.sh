#!/bin/bash

# PFM Test:
# https://hermes.informal.systems/documentation/forwarding/test.html
cargo install ibc-relayer-cli --bin hermes --locked

./hermes-init.sh
./start.sh bitsongd test-1 ./data 26657 26656 6060 9090 ubtsg
./start.sh bitsongd test-2 ./data 27657 27656 7060 10090 ubtsg

# start hermes
hermes start

# Send tokens
hermes tx ft-transfer \
 --timeout-seconds 10000 \
 --dst-chain test-2 \
 --src-chain test-1 \
 --src-port transfer \
 --src-channel channel-0 \
 --amount 1000000000 \
 --denom ubtsg

sleep 2

hermes tx ft-transfer \
 --denom ubtsg \
 --receiver bitsong1x3s7sdrq6zg7r8l8apt9pstjfsg5a8vndxjlum \
 --memo '{"forward": {"receiver": "bitsong1v6f23vwenxfd4s4wsaeww82gmryqd009603gtz", "port": "transfer", "channel": "channel-1"}}' \
 --timeout-seconds 120 \
 --dst-chain test-2 \
 --src-chain test-1 \
 --src-port transfer \
 --src-channel channel-0 \
 --amount 2000000000
