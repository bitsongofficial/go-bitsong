To test pool creation, swap, deposit, withdraw, it requires to do following steps.

1. First run the old binary.(Install old binary using fantoken branch)

2. Submit a upgrade proposal.
e.g. `bitsongd tx gov submit-proposal software-upgrade "Gravity-DEX" --title "Gravity-DEX" --description "upgrade" --from validator --upgrade-height=100 --deposit 10000000ubtsg `--chain-id localnet --keyring-backend test -y

3. Vote the proposal and make it pass.
e.g. `bitsongd tx gov vote 1 yes --from=validator --chain-id=localnet --keyring-backend=test -y`
Then the proposal will be pased and the chain will be halted in block 100

4. Run the new binary which has liquidity module integrated in it.(Install new binary using liquidity branch)

5. Run `bitsongd --home .data/localnet start --pruning=nothing`

6. To test pool creation, run `sh scripts/sample-create-pool.sh`

7. To test deposit, run `sh scripts/sample-create-pool.sh`

8. To test swap, run `sh scripts/sample-create-pool.sh`

9. To test withdraw, run `sh scripts/sample-withdraw.sh`
