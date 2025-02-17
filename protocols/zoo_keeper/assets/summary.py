import glob
import os
import sys

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

print("throughput: " + str(len(all_requests) / test_time) + "\n")

binned_requests = [0] * (test_time + 1)

for _, end in all_requests:
    bin_index = int(end - min_start_time)
    if 0 <= bin_index and bin_index < len(binned_requests):
        binned_requests[bin_index] += 1

with open(OUTPUT_FILE, "w") as output:
    for second, count in enumerate(binned_requests):
        output.write(f"{second}, {count}\n")
