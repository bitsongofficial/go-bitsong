# LocalBitSong

LocalBitSong is a complete BitSong testnet containerized with Docker and orchestrated with a simple docker-compose file. LocalBitSong comes preconfigured with opinionated, sensible defaults for a standard testing environment.

LocalBitSong comes in two flavors:

1. No initial state: brand new testnet with no initial state. 
2. TBD - With mainnet state: creates a testnet from a mainnet state export

Both ways, the chain-id for LocalBitSong is set to 'LocalBitSong'.

## Prerequisites

Ensure you have docker and docker-compose installed:

```sh
# Docker
sudo apt-get remove docker docker-engine docker.io
sudo apt-get update
sudo apt install docker.io -y

# Docker compose
sudo apt install docker-compose -y
```

## 1. LocalBitSong - No Initial State

The following commands must be executed from the root folder of the BitSong repository.

1. Make any change to the bitsong code that you want to test

2. Initialize LocalBitSong:

```bash
make localnet-init
```

The command:

- Builds a local docker image with the latest changes
- Cleans the `$HOME/.bitsongd-local` folder

3. Start LocalBitSong:

```bash
make localnet-start
```

> Note
>
> You can also start LocalBitSong in detach mode with:
>
> `make localnet-startd`

4. (optional) Add your validator wallet and 10 other preloaded wallets automatically:

```bash
make localnet-keys
```

- These keys are added to your `--keyring-backend test`
- If the keys are already on your keyring, you will get an `"Error: aborted"`
- Ensure you use the name of the account as listed in the table below, as well as ensure you append the `--keyring-backend test` to your txs
- Example: `bitsongd tx bank send lo-test2 bitsong1cyyzpxplxdzkeea7kwsydadg87357qnahakaks --keyring-backend test --chain-id LocalBitSong`

5. You can stop chain, keeping the state with

```bash
make localnet-stop
```

6. When you are done you can clean up the environment with:

```bash
make localnet-clean
```

## LocalBitSong Accounts and Keys

LocalBitSong is pre-configured with one validator and 9 accounts with BTSG balances.

| Account   |  | |
|-----------|--------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| val   |  |   |
| del1  |  |   |
| del2  |  |   |
| del3  |  |   |
To list all keys in the keyring named `test`
```bash
bitsongd keys list --keyring-backend test
```

To import an account into the keyring `test`. NOTE: replace the address with any of the above user accounts. 
```bash
bitsongd keys add bitsong1regz7kj3ylg2dn9rl8vwrhclkgz528mf0tfsck --keyring-backend test --recover
```

## FAQ

Q: How do I enable pprof server in LocalBitSong?

A: everything but the Dockerfile is already configured. Since we use a production Dockerfile in LocalBitSong, we don't want to expose the pprof server there by default. As a result, if you would like to use pprof, make sure to add `EXPOSE 6060` to the Dockerfile and rebuild the LocalBitSong image.
