<!-- 
order: 7
-->

# Client

## Transactions

The `transactions` commands allow users to `create` and `claim` for _merkledrops_.

```bash=
bitsongd tx merkledrop --help
```
### create

```bash=
bitsongd tx merkledrop create [account-file] [output-file] \
	--denom=ubtsg \
	--start-height=1 \
	--end-height=10 \
	--from=<key-name> -b block --chain-id <chain-id>
```

### claim

```bash=
bitsongd tx merkledrop claim [merkledrop-id] \
	--proofs=[proofs-list] \
	--amount=[amount-to-claim] \
	--index=[level-index]
	--from=<key-name> -b block --chain-id <chain-id>
```

## Query

The `query` commands allow users to query the _merkledrop_ module.

```bash=
bitsongd q merkledrop --help
```

### detail by id

```bash=
bitsongd q merkledrop detail [id]
```

### if index and id have been claimed

```bash=
bitsongd q merkledrop index-claimed [id] [index]
```

### params

```bash=
bitsongd q merkledrop params
```