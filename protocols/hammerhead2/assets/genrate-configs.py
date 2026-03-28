import os
import yaml
import argparse

# Set up command-line argument parsing
parser = argparse.ArgumentParser(description='Generate node-parameters.yml and client-parameters.yml files.')

parser.add_argument('--wave_length', type=int, default=5, help='Wave length (default: 3)')
parser.add_argument('--number_of_leaders', type=int, default=3, help='Number of leaders (default: 3)')
parser.add_argument('--enable_pipelining', type=str, default="True", help='Enable pipelining (default: True)')
parser.add_argument('--consensus_only', type=str, default="True", help='Consensus only (default: True)')
parser.add_argument('--enable_synchronizer', type=str, default="True", help='Enable synchronizer (default: true)')
parser.add_argument('--initial_delay_secs', type=int, default=5, help='Initial delay in seconds (default: 1)')
parser.add_argument('--initial_delay_nanos', type=int, default=500, help='Initial delay in nanoseconds (default: 500)')
parser.add_argument('--load', type=int, default=1000, help='Load (default: 1000)')
parser.add_argument('--transaction_size', type=int, default=512, help='Transaction size (default: 512)')
parser.add_argument('--output_dir', type=str, default="protocols/hammerhead2/assets/", help='configuration file output directory')
parser.add_argument('--max_leaders_per_round', type=int, default=5, help='Maximum number of leaders (default: 5)')
parser.add_argument('--scheduler_interval', type=int, default=100, help='Interval at which to recalc the number of leaders (default: 100)')

args = parser.parse_args()

# Convert string arguments to boolean
def str_to_bool(value):
    return value.lower() == 'true'

# File 1: node-parameters.yml
node_parameters = {
    'leader_timeout': {
        'secs': 0,
        'nanos': 250000000
    },
    'wave_length': args.wave_length,
    'number_of_leaders': args.number_of_leaders,
    'enable_pipelining': str_to_bool(args.enable_pipelining),
    'consensus_only': str_to_bool(args.consensus_only),
    'enable_synchronizer': str_to_bool(args.enable_synchronizer),
    'max_leaders_per_round': args.max_leaders_per_round,
    'scheduler_interval': args.scheduler_interval
}

# Write node-parameters.yml
with open(os.path.join(args.output_dir, 'node-parameters.yml'), 'w') as file:
    yaml.dump(node_parameters, file, default_flow_style=False)

# File 2: client-parameters.yml
client_parameters = {
    'initial_delay': {
        'secs': args.initial_delay_secs,
        'nanos': args.initial_delay_nanos
    },
    'load': args.load,
    'transaction_size': args.transaction_size
}

# Write client-parameters.yml
with open(os.path.join(args.output_dir, 'client-parameters.yml'), 'w') as file:
    yaml.dump(client_parameters, file, default_flow_style=False)

print(f"Files generated in {args.output_dir}:")
print("- node-parameters.yml")
print("- client-parameters.yml")
