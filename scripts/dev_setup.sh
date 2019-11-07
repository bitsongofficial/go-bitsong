#!/bin/bash
GO_VERSION=1.12.13

echo "Update Ubuntu..."
sudo apt update && sudo apt upgrade -y

echo "Installs packages necessary to run go..."
sudo apt install build-essential libleveldb1v5 git unzip -y

echo "Install go${GO_VERSION}..."
sudo wget -c https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local

echo "Updates environmental variables to include go..."
cat <<EOF >> $HOME/.profile
export GOPATH=\$HOME/go
export GO111MODULE=on
export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin
EOF
source $HOME/.profile

echo "Compile go-bitsong"
make install

echo "Finish!"
