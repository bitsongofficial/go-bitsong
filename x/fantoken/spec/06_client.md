<!-- 
order: 6
-->

# Client

## Transactions

The `transactions` commands allow users to `issue`, `mint`, `burn`, `disable minting`, `transfer minting and editing capabilities` for _fan tokens_.

```bash=
bitsongd tx fantoken --help
```
### issue

```bash=
bitsongd tx fantoken issue \
    --name "fantoken name" \
    --symbol "bitangel" \
    --max-supply 100000000000 \
    --uri "ipfs://...." \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### mint

```bash=
bitsongd tx fantoken mint [amount][denom] \
    --recipient <address> \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### burn

```bash=
bitsongd tx fantoken burn [amount][denom] \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### set-authority

```bash=
bitsongd tx fantoken set-authority [denom] \
    --new-authority <address> \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### set-minter

```bash=
bitsongd tx fantoken set-minter [denom] \
    --new-minter <address> \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### set-uri

```bash=
bitsongd tx fantoken set-uri [denom] \
    --uri <uri> \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

### disable-mint

```bash=
bitsongd tx fantoken disable-mint [denom] \
    --from <key-name> -b block --chain-id <chain-id> --fees <fee>
```

## Query

The `query` commands allow users to query the `fantoken` module.

```bash=
bitsongd q fantoken --help
```

### denom

```bash=
bitsongd q fantoken denom <denom>
```

### authority

```bash=
bitsongd q fantoken authority <address>
```

### params

```bash=
bitsongd q fantoken params
```