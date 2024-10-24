import sys
import time

import etcd3

etcd = etcd3.client()

# Variables
requests_per_second = int(sys.argv[2])  # Requests per second
key_prefix = "benchmark_key"
value = "benchmark_value"

# Metrics
latencies = []
successful_requests = 0


# Function to send a request
def send_request(i):
    global successful_requests
    start_time = time.time()
    try:
        etcd.put(f"{key_prefix}_{i}", value)
        latency = time.time() - start_time
        latencies.append(latency)
        successful_requests += 1
    except Exception as e:
        var = None

DURATION = sys.argv[1]

# Function to benchmark requests
def benchmark():
    i = 1
    while time.time() - start_time < int(DURATION):
        send_request(i)
        i = i + 1


# Start benchmark
start_time = time.time()
benchmark()
end_time = time.time()

# Calculate results
total_time = end_time - start_time
average_latency = sum(latencies) / len(latencies) if latencies else 0
throughput = successful_requests / total_time


print(f"{average_latency * 1000:.6f},{throughput:.2f}")
