import sys
from collections import defaultdict
import matplotlib.pyplot as plt
from concurrent.futures import ThreadPoolExecutor


def extract_start_end_pairs(filename):
    """Extract start and end times from a single file."""
    pairs = []
    try:
        with open(filename, 'r') as f:
            for line in f:
                parts = line.strip().split(',')
                if len(parts) >= 2:
                    start_time = int(parts[0].strip())
                    end_time = int(parts[1].strip())
                    if (40000000 - 1) < start_time and end_time < 100000001:
                        pairs.append((start_time, end_time))
    except FileNotFoundError:
        print(f"File not found: {filename}")
    except Exception as e:
        print(f"Error reading {filename}: {e}")
    return filename, pairs


def calculate_throughput_latency(start_end_pairs, duration):
    """Calculate per second throughput and average latency from start and end times."""
    throughput = defaultdict(int)
    latency_sum = defaultdict(int)
    latency_count = defaultdict(int)

    for start, end in start_end_pairs:
        end_sec = end // 1_000_000
        latency = end - start
        throughput[end_sec] += 1
        latency_sum[end_sec] += latency
        latency_count[end_sec] += 1

    average_latency = {sec: latency_sum[sec] / latency_count[sec] / 1_000 for sec in latency_sum}

    total_requests = sum(throughput.values())
    total_latency = sum(latency_sum.values())
    total_seconds = int(duration)
    overall_throughput = total_requests / total_seconds
    overall_average_latency = total_latency / total_requests / 1_000

    return {
        "throughput": throughput,
        "average_latency": average_latency,
        "overall_throughput": overall_throughput,
        "overall_avg_latency": overall_average_latency
    }


def plot_throughput(throughput, name):
    seconds = [sec - 40 for sec in sorted(throughput.keys())]
    throughput_values = [throughput[sec+40] for sec in seconds]
    plt.figure(figsize=(10, 6))
    plt.plot(seconds, throughput_values, label='Mysticeti', color="black")
    plt.xlabel("Time (s)")
    plt.ylabel("Transactions/sec")
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/" + name + "_throughput.pdf")


def plot_latency(average_latency, name):
    seconds = [sec - 40 for sec in sorted(average_latency.keys())]
    latency_values = [average_latency[sec+40] for sec in seconds]
    plt.figure(figsize=(10, 6))
    plt.plot(seconds, latency_values, label='Mysticeti', color='black')
    plt.xlabel("Time (s)")
    plt.ylabel("Latency (ms)")
    plt.grid(True)
    plt.legend()
    plt.savefig("logs/" + name + "_latency.pdf")


def main():
    if len(sys.argv) < 4:
        print("Usage: python3 performance_graph.py <name> <duration> <file1> <file2> ...")
        return

    output_name = sys.argv[1]
    duration = sys.argv[2]
    filenames = sys.argv[3:]

    with ThreadPoolExecutor() as executor:
        results = executor.map(extract_start_end_pairs, filenames)

    file_data = [(fname, pairs) for fname, pairs in results if pairs]

    if not file_data:
        print("No valid data found.")
        return

    metrics = {}
    with ThreadPoolExecutor() as executor:
        futures = {executor.submit(calculate_throughput_latency, pairs, duration): fname
                   for fname, pairs in file_data}

        for future in futures:
            fname = futures[future]
            metrics[fname] = future.result()


    best_file = max(metrics.items(), key=lambda x: x[1]['overall_throughput'])

    best_name = best_file[0]
    best_metrics = best_file[1]

    # print(f"Best file: {best_name}")
    print(f"{best_metrics['overall_throughput']:.2f},{best_metrics['overall_avg_latency']:.2f} ms")

    plot_throughput(best_metrics["throughput"], output_name)
    plot_latency(best_metrics["average_latency"], output_name)


if __name__ == "__main__":
    main()
