from mergexp import *
net = Network('torturefeb2025', addressing==ipv4, routing==static)
nodes = []
for i in range(20):
    node_id = i + 1
    node_name = f'node{node_id}'
    node = net.node(node_name, memory.capacity==gb(16), proc.cores==16, image=="2004", disk.capacity==gb(20))
    nodes.append(node)
net.connect(nodes, latency==ms(10))
experiment(net)