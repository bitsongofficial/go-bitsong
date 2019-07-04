# Guide to Joining The BitSong Testnet

This article will provide a step-by-step tutorial for setting up a server and getting started as a validator on the BitSong Network testnet. You can view the code and see the latest release at our [GitHub](https://github.com/bitsongofficial/go-bitsong).

## Getting Started

First off, you’ll need to setup a server. Having a dedicated server helps ensure that your validator is highly available and doesn’t go offline. BitSong Network uses Tendermint consensus, which selects a  _leader_  for each block. If your validator is offline when it gets chosen as a leader, consensus will take longer, and you could even get  [_slashed_](https://medium.com/coinmonks/cosmos-atom-staking-guide-4a4e703c998a)!

For this guide, we’ll be using a server with the following specifications:

-   Ubuntu 18.04 OS
-   2 CPUs
-   4GB RAM
-   24GB SSD
-   Allow incoming connections on ports 26656
-   Static IP address (Elastic IP for AWS, floating IP for DigitalOcean,  _etc_)

You can get a server with these specifications on most cloud service providers (AWS, DigitalOcean, Google Cloud Platform, Linode, etc).

After logging into your server, we’ll install security updates and the required packages to run BitSong:

```bash
# Update Ubuntu
sudo apt update

# Installs packages necessary to run go
sudo apt upgrade -y 

# Installs go
sudo apt install build-essential libleveldb1v5 git unzip -y
wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
sudo tar -xvf go1.12.5.linux-amd64.tar.gz
sudo mv go /usr/local

# Updates environmental variables to include go
cat <<EOF >> ~/.profile
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

source ~/.profile
```
To verify that go is installed:

```bash
# Should return go version go1.12.5 linux/amd64
go version
```

## Install the Go-BitSong daemon

Next, we’ll install the software needed to run the BitSong blockchain.
1) First, create your own fork so that you can submit a genesis transaction pull request if necessary.
2) Head over to [GitHub](https://github.com/bitsongofficial/go-bitsong) and click “Fork.”

```bash
# Clone your repository
git clone https://github.com/<YOUR-USERNAME>/go-bitsong.git

# install binaries
cd go-bitsong
make install
```

Now, we’ll setup the `bitsongd` software to run the current BitSong testnet:

```bash
# Replace <your-moniker> with the publicly viewable name for your validator.
# bitsong-testnet-1 is the name of the current testnet
bitsongd init --chain-id bitsong-testnet-1 <your-moniker>
```

**Note:** `bitsongd init` sets the `node-id` of your validator. You can see this value by doing `bitsongd tendermint show-node-id`. The `node-id`value will become part of your genesis transaction, so if you are planning on submitting a genesis transaction, don’t reset your `node-id` by doing `bitsongd unsafe-reset-all` or changing your public IP address.

```bash
# Create a wallet for your node. <your-wallet-name> is just a human readable name you can use to remember your wallet. It can be the same or different than your moniker.

bitsongcli keys add <your_wallet_name>
```
This will spit out your recovery mnemonic. 
**Be sure to back up your mnemonic before proceeding to the next step!**

## Submit a genesis transaction

If you are planning on participating in the genesis of the BitSong testnet, you can follow along here and create a genesis transaction that you can submit as a pull request before launch. Otherwise, skip to the section about obtaining some coins from the faucet. If you are participating in genesis, it is expected that your validator will be up and available at all times during the testnet. If you can’t commit to this, we recommend joining via the faucet after the testnet is live.

```bash
# Create an account in genesis with 1000000000 ubtsg (1000 btsg) tokens. Don't change the amount of ubtsg tokens so that we can have equal distribution among genesis participants.
# 1 btsg = 1000000ubtsg
bitsongd add-genesis-account $(bitsongcli keys show <your_wallet_name> -a) 1000000000ubtsg

# Sign a gentx that creates a validator in the genesis file for your account. Note to pass your public ip to the --ip flag.
bitsongd gentx --name <your_wallet_name> --amount 1000000000ubtsg --ip <your-public-ip>
```

This will write your genesis transaction to `$HOME/.bitsongd/config/gentx/gentx-<gen-tx-hash>.json`. This should be the only file in your `gentx` directory. If you have more than one, delete them and repeat the `gentx` command above.

Now we will submit the transaction as a PR to be included in the genesis block:

```bash
# create a branch for your pr submission
cd $HOME/go-bitsong
git checkout -b genesis-<your-moniker>

# check that there's only one gentx
ls $HOME/.bitsongd/config/gentx

# copy your gentx
cp $HOME/.bitsongd/config/gentx/* $HOME/go-bitsong/testnet-1/gentx/

# Add and commit your changes
git add testnet-1/gentx/*
git commit -m "feat: gentx for <your-moniker>"

# Push your branch to the remote repositor
git push -u origin genesis-<your-moniker>
```

Now go to BitSong's [GitHub repo](https://github.com/bitsongofficial/go-bitsong/pulls) and select **New Pull Request**

Create a pull request for `<github-username>/go-bitsong:genesis-<your-moniker>` against the `develop` branch of the BitSong repo.

1) New Pull Request
2) Click on "compare across forks"

```
base repository: bitsongofficial/go-bitsong
base: develop

head repository: <YOUR-USERNAME>/go-bitsong
compare: genesis-<your-moniker>
```

We’ll make sure to promptly review your PR, let you know if there are any issues, and merge it in!

## Launching the testnet

On **July 10, 2019 at 15:00 UTC**, BitSong will release the proposed genesis block. All validators are invited to review the included transactions, agree on the final genesis block, and start validating on **July 11, 2019 at 15:00 UTC.**

To start validating the testnet after the genesis has been released, run the following commands:

```bash
# Delete the old genesis
rm -f ~/.bitsongd/config/genesis.json

# Copy the genesis file to the bitsongd directory
wget https://raw.githubusercontent.com/bitsongofficial/go-bitsong/develop/testnet-1/genesis.json -P ~/.bitsongd/config

# Create log files for bitsongd
sudo mkdir -p /var/log/bitsongd && sudo touch /var/log/bitsongd/bitsongd.log && sudo touch /var/log/bitsongd/bitsongd_error.log

# create a systemd file to run the bitsongd daemon
# replace <your_user> where necessary
sudo tee /etc/systemd/system/bitsongd.service > /dev/null <<EOF
[Unit]
Description=BitSong Network Daemon
After=network-online.target

[Service]
User=<your_user>
ExecStart=/home/<your_user>/go/bin/bitsongd start
StandardOutput=file:/var/log/bitsongd/bitsongd.log
StandardError=file:/var/log/bitsongd/bitsongd_error.log
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF

# Start the node
sudo systemctl enable bitsongd
sudo systemctl start bitsongd
```

To check on the status of the node

```bash
bitsongcli status
sudo journalctl -u bitsongd -f
```

To view the logs

```bash
# Standard output of bitsongd
tail -f /var/log/bitsongd/bitsongd.log

# Standard error of bitsongd
tail -f /var/log/bitsongd/bitsongd_error.log
```

After the BitSong Network blockchain reaches a quorum, the testnet will be officially launched!

If a quorum is not reached by 15:00 UTC on July 11, we will coordinate further communication through the BitSong validator [Discord Validator Chat](https://discord.gg/qSFUps6).