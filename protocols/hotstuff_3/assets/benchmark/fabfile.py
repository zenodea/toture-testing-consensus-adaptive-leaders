from fabric import task

from aws.remote import Bench, BenchError
from benchmark.logs import ParseError, LogParser
from benchmark.utils import Print


@task
def install(ctx):
    ''' Install HotStuff on all machines '''
    try:
        Bench(ctx).install()
    except BenchError as e:
        Print.error(e)


@task
def remote(ctx, pid=1, attack_duration=2, num_replicas=3, param_load=4, param_size=5):
    ''' Run benchmarks on AWS '''
    bench_params = {
        'nodes': [int(num_replicas)],
        'rate': [int(param_load)],
        'tx_size': int(param_size),
        'faults': 0,
        'duration': int(attack_duration),
        'runs': 1,
    }
    node_params = {
        'consensus': {
            'timeout_delay': 5_000,
            'sync_retry_delay': 5_000,
            'max_payload_size': 1_000,
            'min_block_delay': 100
        },
        'mempool': {
            'queue_capacity': 100_000,
            'sync_retry_delay': 5_000,
            'max_payload_size': 500_000,
            'min_block_delay': 100
        }
    }
    try:
        Bench(ctx).run(bench_params, node_params, int(pid), debug=False)
    except BenchError as e:
        Print.error(e)


@task
def kill(ctx):
    ''' Stop any HotStuff execution on all machines '''
    try:
        Bench(ctx).kill()
    except BenchError as e:
        Print.error(e)


@task
def logs(ctx):
    ''' Print a summary of the logs '''
    try:
        print(LogParser.process('./logs').result())
    except ParseError as e:
        Print.error(BenchError('Failed to parse logs', e))
