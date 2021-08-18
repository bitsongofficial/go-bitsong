<p align="center" background="black"><img src="bitsong-logo.png" width="398"></p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/BitSongOfficial/go-bitsong/blob/master/LICENSE)

**BitSong** is a new music streaming platform based on [Tendermint](https://github.com/tendermint/tendermint) consensus BFT, the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits and the [IPFS](https://ipfs.io/) distribuited filesystem. Please make sure you study these projects as well if you are not already familiar.

**BitSong** is a project dedicated to musicians and their fans, which aims to overcome the bureaucratic and economic obstacles within this industry and reward artists and users for simply using the platform.

**Artists** need no longer to wait several months before a record label sends various reports, they can check the progress in real time directly within the Wallet.

_NOTE: This is alpha software. Please contact us if you aim to run it in production._

**Note**: Requires [Go 1.13.6+](https://golang.org/dl/)

# Install BitSong Blockchain (mpeg21 test)

There are many ways you can install BitSong Blockchain Testnet node on your machine.

## From Source
1. **Install Go** by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH` and `$PATH` environment variables, for example:
    ```bash
    wget https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
    sudo tar -xvzf go1.13.6.linux-amd64.tar.gz
    sudo mv go /usr/local
     
    cat <<EOF >> ~/.profile  
    export GOPATH=$HOME/go  
    export GO111MODULE=on  
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin  
    EOF
    ```
2. **Clone BitSong source code to your machine**
    ```bash
    mkdir -p $GOPATH/src/github.com/BitSongOfficial
    cd $GOPATH/src/github.com/BitSongOfficial
    git clone https://github.com/BitSongOfficial/go-bitsong.git
    cd go-bitsong
    git checkout mpeg21
    ```
  3. **Compile**
		```bash
		# Install the app into your $GOBIN
		make install
		# Now you should be able to run the following commands:
		bitsongd help
		bitsongcli help
		```
		The latest `go-bitsong version` is now installed.
3. **Run BitSong**
	```bash
	bitsongd start
	```

## Running the test network and using the commands

To initialize configuration and a `genesis.json` file for your application and an account for the transactions, start by running:

>  _*NOTE*_: In the below commands addresses are are pulled using terminal utilities. You can also just input the raw strings saved from creating keys, shown below. The commands require [`jq`](https://stedolan.github.io/jq/download/) to be installed on your machine.

>  _*NOTE*_: If you have run the tutorial before, you can start from scratch with a `bitsongd unsafe-reset-all` or by deleting both of the home folders `rm -rf ~/.bitsong*`

>  _*NOTE*_: If you have the Cosmos app for ledger and you want to use it, when you create the key with `bitsongcli keys add jack` just add `--ledger` at the end. That's all you need. When you sign, `jack` will be recognized as a Ledger key and will require a device.

```bash
# Initialize configuration files and genesis file
bitsongd init MyValidator --chain-id bitsong-localnet

# Configure your CLI to eliminate need for chain-id flag
bitsongcli config chain-id bitsong-localnet
bitsongcli config output json
bitsongcli config indent true
bitsongcli config trust-node true
bitsongcli config keyring-backend test

# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
bitsongcli keys add jack

# Copy the `Address` output here and save it for later use
bitsongcli keys add alice

# Add both accounts, with coins to the genesis file
bitsongd add-genesis-account jack 100000000000ubtsg --keyring-backend test
bitsongd add-genesis-account alice 100000000000ubtsg --keyring-backend test

# Generate the transaction that creates your validator
bitsongd gentx --name jack --amount=10000000ubtsg --keyring-backend test

# Add the generated bonding transaction to the genesis file
bitsongd collect-gentxs
bitsongd validate-genesis

# Now its safe to start `bitsongd`
bitsongd start
```

You can now start `bitsongd` by calling `bitsongd start`. You will see logs begin streaming that represent blocks being produced, this will take a couple of seconds.

Open another terminal to run commands against the network you have just created:

```bash
# First check the accounts to ensure they have funds
bitsongcli query account $(bitsongcli keys show jack -a)
bitsongcli query account $(bitsongcli keys show alice -a)
```

# Transactions
You can now start the first transaction

```bash
bitsongcli tx send jack $(bitsongcli keys show alice -a) 10ubtsg --from jack -b block
```

# Query
## Query an account

```bash
bitsongcli query account $(bitsongcli keys show jack -a)
```

## Query total supply

```bash
bitsongcli query supply total
```

# Branch Details

There are several branches in the repo for now.

master = the most stable branch which was released and already running.

develop = the branch which has the latest dev features.

v0.42.x = the branch where only ibc module was enabled on the default cosmos-sdk chain.

module/fantoken = v0.42.x + fantoken module

module/liquidity = module/fantoken + liquidity module

module/nft = module/liquidity + nft module

The reason to have v0.42.x, module/fantoken, module/liquidity, module/nft separately is because we decided to activate every module one by one by governance proposal and we want to build binary easily based on the separated branches.

ex:) Let's assume module/fantoken was activated and now want to activate the module/liquidity

the bitsong binary of module/fantoken is already running now. Using this binary, need to submit chain upgrade proposal with target block. In the target block, if the proposal passes, the chain will be halted. After that, need to get new bitsong binary using module/liquidity and run it.

# Module MPEG21

## Store MPEG21 MCO Contract
```bash
bitsongcli tx mpeg21 store-mco ./testdata/mpeg21/download-no-label.json \
  --from jack \
  --gas 400000 \
  -b block
```

## Query MPEG21 MCO Contract by ID
```bash
bitsongcli query mpeg21 mco-id contract6
```t

## Query all MPEG21 MCO Contracts
```bash
bitsongcli query mpeg21 mco-all
```


## Resources
- [Official Website](https://bitsong.io)

## Buy/Sell BTSG ERC20
- [Uniswap BTSG/ETH](https://app.uniswap.org/#/swap?outputCurrency=0x05079687d35b93538cbd59fe5596380cae9054a9)
- [Uniswap Liquidity Pool BTSG/ETH](https://uniswap.info/pair/0x98195a436C46EE53803eDA921dC3932073Ed7f4d)

### Community
- [Discord](https://discord.gg/mZC9Yk3)
- [Twitter](https://twitter.com/BitSongOfficial)
- [Telegram Channel (English)](https://t.me/BitSongOfficial)
- [Medium](https://medium.com/@BitSongOfficial)
- [Reddit](https://www.reddit.com/r/bitsong/)
- [Facebook](https://www.facebook.com/BitSongOfficial)
- [BitcoinTalk ANN](https://bitcointalk.org/index.php?topic=2850943)
- [Linkedin](https://www.linkedin.com/company/bitsong)
- [Instagram](https://www.instagram.com/bitsong_official/)

## License

MIT License

## Versioning

### SemVer

BitSong uses [SemVer](http://semver.org/) to determine when and how the version changes.
According to SemVer, anything in the public API can change at any time before version 1.0.0

To provide some stability to BitSong users in these 0.X.X days, the MINOR version is used
to signal breaking changes across a subset of the total public API. This subset includes all
interfaces exposed to other processes, but does not include the in-process Go APIs.
