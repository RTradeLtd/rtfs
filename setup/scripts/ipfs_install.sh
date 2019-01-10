#! /bin/bash

IPFS_PATH="/ipfs"
PRIVATE_NODE="no"
VERSION="v0.4.18"
export IPFS_PATH="/ipfs"

cd ~ || exit

sudo mkdir /ipfs
sudo chown -R rtrade:rtrade /ipfs
echo "[INFO] Downloading IPFS"
wget "https://dist.ipfs.io/go-ipfs/${VERSION}/go-ipfs_${VERSION}_linux-amd64.tar.gz"
tar zxvf go-ipfs*.gz
rm -- *gz
cd go-ipfs || exit
echo "[INFO] Running ipfs install script"
sudo ./install.sh