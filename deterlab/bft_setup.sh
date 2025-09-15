#!/bin/bash

NODES=$(seq 1 11)
GO_VERSION="1.19"

cmd="set -e;sudo apt -y update;sudo apt-mark hold grub-efi-amd64 grub-pc grub-common; sudo apt -y upgrade;sudo apt-get -y autoremove; sudo apt install -y git iproute2 cmake build-essential clang pkg-config libssl-dev python3 python3-pip openjdk-11-jdk;sudo rm -rf /usr/local/go; wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz;sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz;rm go${GO_VERSION}.linux-amd64.tar.gz;echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc;. ~/.bashrc;export PATH=$PATH:/usr/local/go/bin:$HOME/.cargo/bin; sudo rm -rf ~/.cargo ~/.rustup;curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y;. $HOME/.cargo/env;rustup default stable; sudo chmod u+s /usr/sbin/tc;sudo setcap cap_net_admin,cap_net_raw+ep $(which tc);getcap $(which tc);pip3 install --user --no-cache-dir --force-reinstall fabric invoke decorator boto3 etcd3 kazoo protobuf==3.19.6; sudo apt install -y python3-matplotlib"

for i in $NODES; do
    echo "Processing node_$i..."
    ssh node"$i" "$cmd"
done

echo "All nodes configured."
