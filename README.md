<p align="center" background="black"><img src="bitsong-logo.png" width="398"></p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/BitSongOfficial/go-bitsong/blob/master/LICENSE)

**BitSong** is a new music streaming platform based on [Tendermint](https://github.com/tendermint/tendermint) consensus algorythm and the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits. Please make sure you study these projects as well if you are not already familiar.

**BitSong** is a project dedicated to musicians and their fans, which aims to overcome the bureaucratic and economic obstacles within this industry and reward artists and users for simply using the platform.

On the **BitSong** platform you (artist) will be able to produce songs in which an advertiser can attach advertisements and users can access from any device. Funds via the Bitsong token `$BTSG` will be credited to the artist wallet immediately and they will be able to withdraw or convert as they see fit.

**Artists** need no longer to wait several months before a record label sends various reports, they can check the progress in real time directly within the Wallet.

_NOTE: This is alpha software. Please contact us if you aim to run it in production._

**TODO**
- [x] ~~Tendermint/ABCI~~
- [x] ~~Cosmos-SDK integration~~
- [x] ~~Token BTSG (1btsg = 100000000ubtsg)~~
- [x] ~~Auth module~~
- [x] ~~Bank module~~
- [x] ~~Crisis module~~
- [x] ~~Params module~~
- [x] ~~Slashing module~~
- [x] ~~Staking module~~
- [ ] Terraform / ansible
- [ ] Terraform / ansible for AWS
- [ ] Terraform / ansible for DO
- [ ] Governance module
- [ ] Proposal module
- [ ] Params managed by our validator
- [ ] Community pool
- [ ] Song store
- [ ] Song listing by proposal and voted by our validator
- [ ] Artist store
- [ ] Playlist store

**Note**: Requires [Go 1.12.4+](https://golang.org/dl/)

# Install BitSong Blockchain

There are many ways you can install BitSong Blockchain Testnet node on your machine.

## From Source
1. **Install Go** by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:
	```bash
	mkdir -p $HOME/go/bin
	echo  "export GOPATH=$HOME/go" >> ~/.bash_profile
	echo  "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
	echo  "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
	echo  "export GO111MODULE=on" >> ~/.bash_profile
	source ~/.bash_profile
	```
2. **Clone BitSong source code to your machine**
	```bash
	mkdir -p $GOPATH/src/github.com/BitSongOfficial
	cd $GOPATH/src/github.com/BitSongOfficial
	git clone https://github.com/BitSongOfficial/go-bitsong.git
	cd go-bitsong
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

## Install on Digital Ocean
1. **Clone repository**
    ```bash
	git clone https://github.com/BitSongOfficial/go-bitsong.git
    chmod +x go-bitsong/scripts/install/install_ubuntu.sh
	```
2. **Run the script**
    ```bash
    go-bitsong/scripts/install/install_ubuntu.sh
    source ~/.profile
	```
3. Now you should be able to run the following commands:
	```bash
	bitsongd help
	bitsongcli help
	```
    The latest `go-bitsongd version` is now installed.

## Running the test network and using the commands

To initialize configuration and a `genesis.json` file for your application and an account for the transactions, start by running:

>  _*NOTE*_: In the below commands addresses are are pulled using terminal utilities. You can also just input the raw strings saved from creating keys, shown below. The commands require [`jq`](https://stedolan.github.io/jq/download/) to be installed on your machine.

>  _*NOTE*_: If you have run the tutorial before, you can start from scratch with a `bitsongd unsafe-reset-all` or by deleting both of the home folders `rm -rf ~/.bitsong*`

>  _*NOTE*_: If you have the Cosmos app for ledger and you want to use it, when you create the key with `bitsongcli keys add jack` just add `--ledger` at the end. That's all you need. When you sign, `jack` will be recognized as a Ledger key and will require a device.

```bash
# Initialize configuration files and genesis file
bitsongd init --chain-id bitsong-test-network-1

# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
bitsongcli keys add jack

# Copy the `Address` output here and save it for later use
bitsongcli keys add alice

# Add both accounts, with coins to the genesis file
bitsongd add-genesis-account $(bitsongcli keys show jack -a) 1000btsg
bitsongd add-genesis-account $(bitsongcli keys show alice -a) 1000btsg

# Configure your CLI to eliminate need for chain-id flag
bitsongcli config chain-id bitsong-test-network-1
bitsongcli config output json
bitsongcli config indent true
bitsongcli config trust-node true
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
bitsongcli tx send --from=$(bitsongcli keys show jack -a)  $(bitsongcli keys show alice -a) 10btsg
```

# Query
Query an account

```bash
bitsongcli query account $(bitsongcli keys show jack -a)
```

## Resources
- [Official Website](https://bitsong.io)

### Community
- [Telegram Channel (English)](https://t.me/BitSongOfficial)
- [Facebook](https://www.facebook.com/BitSongOfficial)
- [Twitter](https://twitter.com/BitSongOfficial)
- [Medium](https://medium.com/@BitSongOfficial)
- [Reddit](https://www.reddit.com/r/bitsong/)
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
