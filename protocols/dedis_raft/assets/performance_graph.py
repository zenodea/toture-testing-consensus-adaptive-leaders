import sys
import matplotlib.pyplot as plt
from collections import defaultdict


def extract_start_end_pairs(filenames):
    """Extract start and end times from the given list of filenames."""
    start_end_pairs = []

    for filename in filenames:
        try:
            with open(filename, 'r') as f:
                for line in f:
                    parts = line.strip().split(',')
                    if len(parts) >= 3:
                        start_time = int(parts[1])
                        end_time = int(parts[2])
                        start_end_pairs.append((start_time, end_time))
        except FileNotFoundError:
            print(f"File not found: {filename}")
        except Exception as e:
            print(f"Error reading {filename}: {e}")

    return start_end_pairs


def calculate_throughput_latency(start_end_pairs):
    """Calculate per second throughput and average latency from start and end times."""
    throughput = defaultdict(int)
    latency_sum = defaultdict(int)
    latency_count = defaultdict(int)

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

    return throughput, average_latency


def plot_throughput(throughput):
    """Plot per second throughput."""
    seconds = sorted(throughput.keys())
    throughput_values = [throughput[sec] for sec in seconds]

    plt.figure(figsize=(10, 6))
    plt.plot(seconds, throughput_values, label='Raft', color='black')
    plt.xlabel('Time (s)')
    plt.ylabel('Transactions/sec')
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/"+sys.argv[1]+"-throughput.pdf")


def plot_latency(average_latency):
    """Plot per second average latency."""
    seconds = sorted(average_latency.keys())
    latency_values = [average_latency[sec] for sec in seconds]

    plt.figure(figsize=(10, 6))
    plt.plot(seconds, latency_values, label='Raft', color='black')
    plt.xlabel('Time (s)')
    plt.ylabel('Latency (ms)')
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/"+sys.argv[1]+"-latency.pdf")


def main():
    if len(sys.argv) < 3:
        print("Usage: python throughput_latency.py <name> <file1> <file2> ...")
        return

    filenames = sys.argv[2:]
    start_end_pairs = extract_start_end_pairs(filenames)
    if not start_end_pairs:
        print("No valid data found.")
        return

    throughput, average_latency = calculate_throughput_latency(start_end_pairs)

    plot_throughput(throughput)
    plot_latency(average_latency)


if __name__ == "__main__":
    main()