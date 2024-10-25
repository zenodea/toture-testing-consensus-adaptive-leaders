import sys
import time
import random

import etcd3

etcd = etcd3.client()
random.seed(time.time())

key_prefix = "ark_key"
value = "k_value"

# Metrics
latencies = []


# Function to send a request
def send_request(i):
    start_time_r = time.time()
    try:
        etcd.put(f"{key_prefix}_{i}_{str(random.randint(1, 10000000))}", value)
        latency = time.time() - start_time_r
        latencies.append(latency)
    except Exception as e:
        var = None


DURATION = sys.argv[1]


# Function to benchmark requests
def benchmark():
    i = 1
    start_time_b = time.time()
    while time.time() - start_time_b < int(DURATION):
        send_request(i)
        i = i + 1


benchmark()

# Calculate results
total_time = DURATION
average_latency = sum(latencies) / len(latencies) if latencies else 0
throughput = len(latencies) / int(DURATION)

print(f"\nperf,{average_latency * 1000:.6f},{throughput:.2f}\n")
