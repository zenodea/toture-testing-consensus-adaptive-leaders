#!/bin/bash

NODES=$(seq 1 19)

GO_VERSION="1.19"

cmd="sudo kill -9 6223;sudo apt-mark hold grub-efi-amd64 grub-pc grub-common; export DEBIAN_FRONTEND=noninteractive; sudo dpkg --configure -a; sudo apt update ; sudo apt install -y git; sudo rm -rf /usr/local/go; wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz; sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz; rm go${GO_VERSION}.linux-amd64.tar.gz; echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc ; source ~/.bashrc; sudo rm -rf ~/.cargo ~/.rustup; curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y ; source ~/.cargo/env; sudo apt install -y iproute2; sudo chmod u+s /usr/sbin/tc; sudo apt install -y build-essential clang; sudo apt install -y pkg-config libssl-dev; sudo apt install -y python3-pip; pip3 install matplotlib; pip3 install --user --no-cache-dir --force-reinstall fabric invoke decorator; pip3 install --user boto3; sudo setcap cap_net_admin,cap_net_raw+ep $(which tc); getcap $(which tc); sudo apt install -y python3 python3-pip; pip3 install etcd3; pip3 install protobuf==3.19.6; sudo apt update ; sudo apt install -y openjdk-11-jdk; sudo apt install -y python3-pip; pip3 install kazoo"

for i in $NODES; do
    echo "Processing node_$i..."
    nohup ssh node"$i" "${cmd}" &
    sleep 5
done

ssh node20 "${cmd}"

sleep 100

echo "All nodes configured."