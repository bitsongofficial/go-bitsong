<!-- 
order: 6
-->

# Client

## Transactions

The `transactions` commands allow users to issue, mint, burn and transfer `fantokens`.

```bash=
bitsongd tx fantoken --help
```
### issue

```bash=
bitsongd tx fantoken issue \
    --name "fantoken name" \
    --symbol "bitangel" \
    --uri "ipfs://...." \
    --max-supply 100000000000 \
    --from <key-name> -b block --chain-id <chain-id>
```

### mint

```bash=
bitsongd tx fantoken mint <denom> \
    --recipient <address> \
    --amount "1" \
    --from <key-name> -b block --chain-id <chain-id>
```

### burn

```bash=
bitsongd tx fantoken burn <denom> \
    --amount "1" \
    --from <key-name> -b block --chain-id <chain-id>
```

### disable-mint

```bash=
bitsongd tx fantoken disable-mint <denom> \
    --from <key-name> -b block --chain-id <chain-id>
```

### transfer-authority

```bash=
bitsongd tx fantoken transfer-authority <denom> \
    --dst-authority <address> \
    --from <key-name> -b block --chain-id <chain-id>
```

## Query

The `query` commands allow users to query the `fantoken` module

```bash=
bitsongd query fantoken --help
```

### denom

```bash=
bitsongd query fantoken denom <denom>
```

Example Output

```bash=
fantoken:
  authority: bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw
  denom: ftF0EA7AE2933E757BB3120E29A58FB63A54B3E726
  max_supply: "100000000000"
  meta_data:
    name: fantoken name
    symbol: angelo
    uri: ipfs://....
  mintable: true
```

### authority

```bash=
bitsongd q fantoken authority <address>
```

Example Output

```bash=
fantokens:
- authority: bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw
  denom: ftE1A74A564AF1AEBD79F0B9E45FA2D60E4E2241FF
  max_supply: "100000000000"
  meta_data:
    name: fantoken name
    symbol: angelo
    uri: ipfs://....
  mintable: true
pagination:
  next_key: null
  total: "0"
```

### params

```bash=
bitsongd q fantoken params
```

Example Output

```bash=
burn_fee:
  amount: "0"
  denom: ubtsg
issue_fee:
  amount: "1000000000"
  denom: ubtsg
mint_fee:
  amount: "0"
  denom: ubtsg
transfer_fee:
  amount: "0"
  denom: ubtsg

```

### total-burn

```bash=
bitsongd q fantoken total-burn
```

Example Output

```bash=
burned_coins:
- amount: "1"
  denom: ftE1A74A564AF1AEBD79F0B9E45FA2D60E4E2241FF
```