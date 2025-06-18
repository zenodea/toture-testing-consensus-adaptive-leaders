import glob
import os
import sys
import matplotlib.pyplot as plt
from collections import defaultdict

INPUT_DIR = sys.argv[1]
OUTPUT_FILE = sys.argv[2]

all_requests = []

for file_path in glob.glob(os.path.join(INPUT_DIR, "*.log")):
    with open(file_path, "r") as file:
        for line in file:
            start, end = map(float, line.strip().split(","))
            all_requests.append((start, end))

if not all_requests:
    print("No data found in files.")
    exit()

print("total requests: " + str(len(all_requests)) + "\n")

min_start_time = min(start for start, _ in all_requests)

print("Min Start time: " + str(min_start_time) + "\n")

max_end_time = max(end for _, end in all_requests)
test_time = int(max_end_time - min_start_time)

# Throughput: count of requests per second
throughput_by_sec = defaultdict(int)
latency_by_sec = defaultdict(list)

for start, end in all_requests:
    sec = int(end - min_start_time)
    latency = (end - start) * 1000  # convert to milliseconds
    throughput_by_sec[sec] += 1
    latency_by_sec[sec].append(latency)

# Sort by time
time_series = sorted(throughput_by_sec.keys())
throughput_values = [throughput_by_sec[t] for t in time_series]
avg_latency_values = [sum(latency_by_sec[t]) / len(latency_by_sec[t]) for t in time_series]

# Plot throughput
plt.figure(figsize=(10, 5))
plt.scatter(time_series, throughput_values, marker='o', label='Throughput (req/s)')
plt.xlabel('Time (s)')
plt.ylabel('Throughput (requests/sec)')
plt.title('Time vs Throughput')
plt.grid(True)
plt.tight_layout()
plt.savefig(OUTPUT_FILE + "/zoo_throughput.pdf")
print("Generated throughput plot in " + OUTPUT_FILE + "/zoo_throughput.pdf")

# Plot latency
plt.figure(figsize=(10, 5))
plt.scatter(time_series, avg_latency_values, marker='o', color='orange', label='Avg Latency (ms)')
plt.xlabel('Time (s)')
plt.ylabel('Latency (ms)')
plt.title('Time vs Latency')
plt.grid(True)
plt.tight_layout()
plt.savefig(OUTPUT_FILE + "/zoo_latency.pdf")
print("Generated latency plot in " + OUTPUT_FILE + "/zoo_latency.pdf")