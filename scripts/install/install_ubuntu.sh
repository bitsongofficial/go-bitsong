#!/usr/bin/env bash

# This script is intended to be run from
# a fresh Digital Ocean droplet with Ubuntu

# change this to a specific release or branch
BRANCH=develop
REPO=github.com/BitSongOfficial/go-bitsong

GO_VERSION=1.12.4

sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install -y make

# get and unpack golang
curl -O https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz
tar -xvf go$GO_VERSION.linux-amd64.tar.gz

# move go binary and add to path
mv go /usr/local
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile

# create the go directory, set GOPATH, and put it on PATH
mkdir go
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile
echo "export GO111MODULE=on" >> ~/.profile
echo "export LEDGER_ENABLED=false" >> ~/.profile
source ~/.profile

# get the code and move into repo
cd go-bitsong

# build & install master
git checkout $BRANCH
make install