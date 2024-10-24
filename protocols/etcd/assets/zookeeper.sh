#!/bin/bash

# Variables
ETCD_VERSION="v3.5.9"
ETCD_DOWNLOAD_URL="https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz"
USERNAME="pasindu"
SSH_KEY="/home/tennage/Pictures/pasindu"
NODES=("10.156.33.141" "10.156.33.142" "10.156.33.143")

install() {
    local node_ip=$1

    echo "Installing and running Etcd on ${node_ip}..."


    # Stop any running Etcd instance on the machine
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo pkill -f etcd"
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo pkill -f etcd"
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo rm -f /usr/local/bin/etcd;sudo rm -f /usr/local/bin/etcdctl; sudo rm -f /usr/local/bin/etcdutl"
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo systemctl stop etcd; sudo systemctl disable etcd; sudo rm -f /etc/systemd/system/etcd.service"

    # Clean up the Etcd data directory to ensure no previous bootstrapping exists
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo rm -rf /home/pasindu/etcd/; sudo rm -r /var/lib/etcd; sudo rm -r /tmp/etcd"

    echo "Etcd is stopped on ${node_ip}"

  # Create the necessary directories and download Etcd if not already downloaded
    ssh -i $SSH_KEY $USERNAME@$node_ip "mkdir -p /home/pasindu/etcd/data"
    ssh -i $SSH_KEY $USERNAME@$node_ip "wget -q $ETCD_DOWNLOAD_URL && tar xzvf etcd-${ETCD_VERSION}-linux-amd64.tar.gz"

    # Move binaries to /home/pasindu/etcd
    ssh -i $SSH_KEY $USERNAME@$node_ip "mv etcd-${ETCD_VERSION}-linux-amd64/etcd* /home/pasindu/etcd/"
    ssh -i $SSH_KEY $USERNAME@$node_ip "rm -rf etcd-${ETCD_VERSION}-linux-amd64" # Clean up extracted files

    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo apt update && sudo apt install -y python3 python3-pip"
    ssh -i $SSH_KEY $USERNAME@$node_ip "pip3 install etcd3"
    ssh -i $SSH_KEY $USERNAME@$node_ip "pip3 install protobuf==3.19.6"

    scp -i $SSH_KEY client.py $USERNAME@$node_ip:/home/pasindu/etcd/

    echo "Etcd is installed on ${node_ip}"
}

run() {
    local node_ip=$1
    local node_name=$2

    nohup ssh -i $SSH_KEY $USERNAME@$node_ip "/home/pasindu/etcd/etcd --name infra${node_name} --initial-advertise-peer-urls http://${node_ip}:2380 --listen-peer-urls http://${node_ip}:2380 --listen-client-urls http://${node_ip}:2379,http://127.0.0.1:2379 --advertise-client-urls http://${node_ip}:2379 --initial-cluster-token etcd-cluster-1 --initial-cluster infra0=http://10.156.33.141:2380,infra1=http://10.156.33.142:2380,infra2=http://10.156.33.143:2380 --initial-cluster-state new"  > logs/etcd${node_name}.log &
    echo "Etcd is running on ${node_ip}"
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_1.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_2.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_3.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_4.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_5.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_6.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_7.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_8.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_9.log &
    nohup ssh -i $SSH_KEY $USERNAME@${node_ip} "python3 /home/pasindu/etcd/client.py 60" > logs/client${node_name}_10.log &
}

kill() {
    local node_ip=$1

    echo "killing Etcd on ${node_ip}..."

    # Stop any running Etcd instance on the machine
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo pkill -f etcd"
    ssh -i $SSH_KEY $USERNAME@$node_ip "sudo pkill -f etcd"
}

rm -r logs/
mkdir logs/
# Loop over all nodes and install Etcd
for i in "${!NODES[@]}"; do
    install "${NODES[$i]}"
done

echo "Etcd installation completed on all nodes."

for i in "${!NODES[@]}"; do
    run "${NODES[$i]}" "$i"
done

echo "Etcd is running on all nodes."

sleep 120

for i in "${!NODES[@]}"; do
    kill "${NODES[$i]}"
done

echo "Etcd killed."

