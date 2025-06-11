#!/bin/bash

NODES=$(seq 1 19)
GO_VERSION="1.19"

cmd="set -e;sudo apt update;sudo apt upgrade;sudo apt-get -y autoremove; sudo apt install -y git iproute2 cmake build-essential clang pkg-config libssl-dev python3 python3-pip openjdk-11-jdk;sudo rm -rf /usr/local/go; wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz;sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz;rm go${GO_VERSION}.linux-amd64.tar.gz;echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc;. ~/.bashrc; sudo rm -rf ~/.cargo ~/.rustup;curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y;. $HOME/.cargo/env;rustup default stable; sudo chmod u+s /usr/sbin/tc;sudo setcap cap_net_admin,cap_net_raw+ep $(which tc);getcap $(which tc);pip3 install --user --no-cache-dir --force-reinstall fabric invoke decorator matplotlib boto3 etcd3 kazoo protobuf==3.19.6"

for i in $NODES; do
    echo "Processing node_$i..."
    nohup ssh node"$i" "$cmd" > "setup_node${i}.log" 2>&1 &
    sleep 5
done

ssh node20 "$cmd"

sleep 100
echo "All nodes configured."
