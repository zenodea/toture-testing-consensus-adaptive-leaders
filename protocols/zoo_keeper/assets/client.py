import random
import sys
import time

from kazoo.client import KazooClient

random.seed(time.time())

# Initialize ZooKeeper client
zk = KazooClient()
zk.start()

key_prefix = "/ark_key"  # ZooKeeper keys are paths (e.g., "/path/to/key")
value = b"k_value"  # ZooKeeper requires byte values

# Duration of the benchmark in seconds, passed as an argument
DURATION = int(sys.argv[1])

# Metrics
latencies = []


# Function to send a request to ZooKeeper
def send_request(i):
    start_time_r = time.time()
    try:
        zk.create(f"{key_prefix}_{i}_{str(random.randint(1, 10000000))}", value, makepath=True)
        latency = time.time() - start_time_r
        latencies.append(latency)
    except Exception as e:
        pass  # Handle or log error if necessary


# Function to benchmark requests
def benchmark():
    i = 1
    start_time_b = time.time()
    while time.time() - start_time_b < DURATION:
        send_request(i)
        i += 1


benchmark()

# Calculate results
average_latency = sum(latencies) / len(latencies) if latencies else 0
throughput = len(latencies) / DURATION

print(f"\nperf,{average_latency * 1000:.6f},{throughput:.2f}\n")

# Close the ZooKeeper client
zk.stop()
zk.close()
