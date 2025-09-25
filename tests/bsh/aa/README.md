# Interchain Bitsong Abstract Accounts

Here we test the deployment and various functions of the Abstract Account Framework on Bitsong.

Once deployed, it will confirm contracts can be made use of by create a default monarch account

## TODO

- ~~implement support for bs-account deployment & test integration~~
- simulate use of various workflows:
  - `account modules`: propose, approve, reject, yank modules
  - `claiming namespaces`
  - `ans-host`: using ans-host to register:
    - common assets on blockchain
    - contracts related to specific protocols
    - ibc-channel data
- `x/smart-accounts integration`

## Usage

This project supports two deployment scenarios:

### With AuthZ (Currently broken. Issue tracked [here](https://github.com/AbstractSDK/abstract/issues/569))

```bash
# Run deployment with AuthZ grants
sh a.sh --enable-authz
```

### Without AuthZ (works as expected)

```bash
# Run deployment without AuthZ grants
sh a.sh --disable-authz
```

## Manual Rust Scripts

You can also run the Rust scripts directly:

### With AuthZ

```bash
cargo run --bin init_contracts -- --authz-granter <GRANTER_ADDRESS>
cargo run --bin full_deploy -- --authz-granter <GRANTER_ADDRESS>
```

### Without AuthZ

```bash
cargo run --bin init_contracts
cargo run --bin full_deploy
```

be sure to `pkill -f bitsongd` once complete (make sure you dont kill any production services!!)
