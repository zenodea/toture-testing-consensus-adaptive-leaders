#!/bin/bash

NODES=$(seq 1 11)

for i in $NODES; do
    echo "Reobooting node_$i..."
    ssh node"$i" "sudo reboot"
done

echo "All nodes reboot."
