# Copyright(C) Facebook, Inc. and its affiliates.
from fabric import task

from benchmark.logs import ParseError, LogParser
from benchmark.remote import Bench, BenchError
from benchmark.utils import Print


@task
def install(ctx):
    ''' Install the codebase on all machines '''
    try:
        Bench(ctx).install()
    except BenchError as e:
        Print.error(e)


@task
def remote(ctx, pid=1, attack_duration=2, num_replicas=3, param_load=4, param_size=5, debug=False):
    ''' Run benchmarks on AWS '''
    bench_params = {
        'faults': 0,
        'nodes': [int(num_replicas)],
        'workers': 1,
        'collocate': True,
        'rate': [int(param_load)],
        'tx_size': int(param_size),
        'duration': int(attack_duration),
        'runs': 1,
    }
    node_params = {
        'header_size': 1_000,  # bytes
        'max_header_delay': 200,  # ms
        'gc_depth': 50,  # rounds
        'sync_retry_delay': 10_000,  # ms
        'sync_retry_nodes': 3,  # number of nodes
        'batch_size': 500_000,  # bytes
        'max_batch_delay': 200  # ms
    }
    try:
        Bench(ctx).run(bench_params, node_params, pid, debug)
    except BenchError as e:
        Print.error(e)


@task
def kill(ctx):
    ''' Stop execution on all machines '''
    try:
        Bench(ctx).kill()
    except BenchError as e:
        Print.error(e)


@task
def logs(ctx):
    ''' Print a summary of the logs '''
    try:
        print(LogParser.process('./logs', faults='?').result())
    except ParseError as e:
        Print.error(BenchError('Failed to parse logs', e))
