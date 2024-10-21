from fabric import task

from benchmark.local import LocalBench
from benchmark.logs import ParseError, LogParser
from benchmark.utils import Print
from benchmark.plot import Ploter, PlotError
from benchmark.instance import InstanceManager
from benchmark.remote import Bench, BenchError


@task
def local(ctx):
    """Run benchmarks on localhost"""
    bench_params = {
        "faults": 0,
        "nodes": 4,
        "rate": 1_000,
        "tx_size": 512,
        "duration": 20,
    }
    node_params = {
        "consensus": {
            "timeout_delay": 1_000,
            "sync_retry_delay": 10_000,
        },
        "mempool": {
            "gc_depth": 50,
            "sync_retry_delay": 5_000,
            "sync_retry_nodes": 3,
            "batch_size": 15_000,
            "max_batch_delay": 10,
        },
    }
    try:
        ret = LocalBench(bench_params, node_params).run(debug=True).result()
        print(ret)
    except BenchError as e:
        Print.error(e)

@task
def install(ctx):
    """Install the codebase on all machines"""
    try:
        Bench(ctx).install()
    except BenchError as e:
        Print.error(e)


@task
def remote(ctx):
    """Run benchmarks on AWS"""
    bench_params = {
        "faults": 0,
        "nodes": [4],
        "rate": [10_000],
        "tx_size": 512,
        "duration": 300,
        "runs": 1,
    }
    node_params = {
        "consensus": {
            "timeout_delay": 5_000,
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
        Bench(ctx).run(bench_params, node_params, debug=False)
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
