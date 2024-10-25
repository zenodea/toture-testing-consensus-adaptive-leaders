import sys
import time

import etcd3

etcd = etcd3.client()


key_prefix = "benchmark_key"
value = "benchmark_value"

# Metrics
latencies = []


# Function to send a request
def send_request(i):
    start_time = time.time()
    try:
        etcd.put(f"{key_prefix}_{i}", value)
        latency = time.time() - start_time
        latencies.append(latency)
    except Exception as e:
        var = None

DURATION = sys.argv[1]

# Function to benchmark requests
def benchmark():
    i = 1
    start_time = time.time()
    while time.time() - start_time < int(DURATION):
        send_request(i)
        i = i + 1


benchmark()


# Calculate results
total_time = DURATION
average_latency = sum(latencies) / len(latencies) if latencies else 0
throughput = len(latencies) / int(DURATION)


print(f"\nperf,{average_latency * 1000:.6f},{throughput:.2f}\n")
