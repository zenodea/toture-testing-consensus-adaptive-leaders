from fabric import task

from benchmark.logs import ParseError, LogParser
from benchmark.remote import Bench, BenchError
from benchmark.utils import Print


@task
def install(ctx):
    """Install the codebase on all machines"""
    try:
        Bench(ctx).install()
    except BenchError as e:
        Print.error(e)


@task
def remote(ctx, pid=1, attack_duration=2, num_replicas=3):
    """Run benchmarks on AWS"""
    bench_params = {
        "nodes": [int(num_replicas)],
        "rate": [5000],
        "tx_size": 18,
        "duration": int(attack_duration),
        "runs": 1,
    }
    node_params = {
        "consensus": {
            "timeout_delay": 2_000,
            "sync_retry_delay": 5_000,
        },
        "mempool": {
            "gc_depth": 50,
            "sync_retry_delay": 5_000,
            "sync_retry_nodes": 3,
            "batch_size": 500_000,
            "max_batch_delay": 100,
        },
    }
    try:
        Bench(ctx).run(bench_params, node_params, int(pid), debug=False)
    except BenchError as e:
        Print.error(e)


@task
def kill(ctx):
    """Stop any HotStuff execution on all machines"""
    try:
        Bench(ctx).kill()
    except BenchError as e:
        Print.error(e)


@task
def logs(ctx):
    """Print a summary of the logs"""
    try:
        print(LogParser.process("./logs", faults="?").result())
    except ParseError as e:
        Print.error(BenchError("Failed to parse logs", e))
