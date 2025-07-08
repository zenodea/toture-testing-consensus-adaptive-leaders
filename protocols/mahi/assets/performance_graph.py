import sys
from collections import defaultdict

import matplotlib.pyplot as plt


def extract_start_end_pairs(filenames):
    """Extract start and end times from the given list of filenames."""
    start_end_pairs = []

    for filename in filenames:
        try:
            with open(filename, 'r') as f:
                for line in f:
                    parts = line.strip().split(',')
                    if len(parts) >= 2:
                        start_time = int(parts[0].strip())
                        end_time = int(parts[1].strip())
                        if (40000000-1) < start_time and end_time < 100000001:
                            start_end_pairs.append((start_time, end_time))
        except FileNotFoundError:
            print(f"File not found: {filename}")
        except Exception as e:
            print(f"Error reading {filename}: {e}")

    return start_end_pairs


def calculate_throughput_latency(start_end_pairs, n,duration):
    """Calculate per second throughput and average latency from start and end times."""
    throughput = defaultdict(int)
    latency_sum = defaultdict(int)
    latency_count = defaultdict(int)

    min_time = min(start for start, _ in start_end_pairs)
    max_time = max(end for _, end in start_end_pairs)

    for start, end in start_end_pairs:
        end_sec = end // 1_000_000  # convert end time to seconds
        latency = end - start

        # Count completed requests per second
        throughput[end_sec] += 1

        # Calculate latency per second
        latency_sum[end_sec] += latency
        latency_count[end_sec] += 1

    # Calculate average latency per second in ms
    average_latency = {sec: latency_sum[sec] / latency_count[sec] / 1_000 for sec in latency_sum}

    # print the overall throughput and average_latency

    total_requests = sum(throughput.values())
    total_latency = sum(latency_sum.values())
    total_seconds = int(duration)
    overall_throughput = total_requests / total_seconds / n
    overall_average_latency = total_latency / total_requests / 1000
    print(f"{overall_throughput:.2f} {overall_average_latency:.2f} ")

    return throughput, average_latency


def plot_throughput(throughput, n):
    """Plot per second throughput."""
    seconds = [sec - 40 for sec in sorted(throughput.keys())]
    throughput_values = [throughput[sec+40]/n for sec in seconds]

    plt.figure(figsize=(10, 6))
    plt.plot(seconds, throughput_values, label='Mahi-Mahi', color="red")
    plt.xlabel("Time (s)")
    plt.ylabel("Transactions/sec")
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/" + sys.argv[1] + "_throughput.pdf")


def plot_latency(average_latency):
    """Plot per second average latency."""
    seconds = [sec - 40 for sec in sorted(average_latency.keys())]
    latency_values = [average_latency[sec+40] for sec in seconds]

    plt.figure(figsize=(10, 6))
    plt.plot(seconds, latency_values, label='Mahi-Mahi', color='red')
    plt.xlabel("Time (s)")
    plt.ylabel("Latency (ms)")
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/" + sys.argv[1] + "_latency.pdf")


def main():
    if len(sys.argv) < 3:
        print("Usage: python throughput_latency.py <name> <file1> <file2> ...")
        return
    duration = sys.argv[2]
    filenames = sys.argv[3:]
    start_end_pairs = extract_start_end_pairs(filenames)
    if not start_end_pairs:
        print("No valid data found.")
        return

    throughput, average_latency = calculate_throughput_latency(start_end_pairs, len(filenames),duration)

    plot_throughput(throughput, len(filenames))
    plot_latency(average_latency)


if __name__ == "__main__":
    main()
