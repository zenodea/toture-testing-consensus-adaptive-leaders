import sys

if len(sys.argv) < 3:
    exit("not enough arguments")

numReplicas = int(sys.argv[1])
numClients = int(sys.argv[2])
replicaIPs = []
clientIPs = []

if len(sys.argv) < 3 + numReplicas + numClients:
    exit("not enough arguments, enter the ips of replicas and clients")

argC = 3
for i in range(argC, argC + numReplicas, 1):
    replicaIPs.append(sys.argv[i])

argC = argC + numReplicas

for i in range(argC, argC + numClients, 1):
    clientIPs.append(sys.argv[i])

# generate the .yml of peers
def print_peers(replicaIPs):
    print("peers:")
    for j in range(1, 1 + numReplicas, 1):
        print("   - name: " + str(j))
        print("     ip: " + str(replicaIPs[j - 1]))
        print("     proxyport: " + str(9000 ))
        print("     recorderport: " + str(10000))

# generate the .yml of clients
def print_clients(clientIPs):
    print("clients:")
    for j in range(1, 1 + numClients, 1):
        print("   - name: " + str(j + 50)) # first client is 21
        print("     ip: " + str(clientIPs[j - 1]))
        print("     clientport: " + str(11000))


print_peers(replicaIPs)
print()
print_clients(clientIPs)
