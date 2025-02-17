import concurrent.futures
import random
import sys
import time

from kazoo.client import KazooClient

DURATION = int(sys.argv[1])
NUM_INSTANCES = int(sys.argv[2])
ID = sys.argv[3]
LOG_PATH = sys.argv[4]


def run_instance(instance_id):
    zk = KazooClient()
    zk.start()

    random.seed(time.time() + instance_id)
    key_prefix = f"/ark_key_{instance_id}"
    value = b"k_value"

    latencies = []
    logs = []

    def send_request(i):
        start_time_r = time.time()
        try:
            zk.create(f"{key_prefix}_{i}_{random.randint(1, 10000000)}", value, makepath=True)
            latency = time.time() - start_time_r
            latencies.append(latency)
            logs.append(f"{start_time_r},{time.time()}\n")
        except Exception:
            pass

    i = 1
    start_time_b = time.time()
    while time.time() - start_time_b < DURATION:
        send_request(i)
        i += 1

    with open(f"{LOG_PATH}{ID}_{instance_id}.log", "w") as file:
        file.writelines(logs)

    zk.stop()
    zk.close()

    return latencies


def main():
    all_latencies = []

    with concurrent.futures.ThreadPoolExecutor(max_workers=NUM_INSTANCES) as executor:
        futures = {executor.submit(run_instance, i): i for i in range(NUM_INSTANCES)}
        for future in concurrent.futures.as_completed(futures):
            all_latencies.extend(future.result())

    total_requests = len(all_latencies)
    average_latency = sum(all_latencies) / total_requests if total_requests else 0
    throughput = total_requests / DURATION

    print(f"\nperf,{average_latency * 1000:.6f},{throughput:.2f}\n")


if __name__ == "__main__":
    main()
