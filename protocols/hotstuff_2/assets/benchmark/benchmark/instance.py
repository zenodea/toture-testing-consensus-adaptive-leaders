import os
import yaml
from collections import OrderedDict

import boto3
from benchmark.settings import Settings, SettingsError
from benchmark.utils import BenchError
from botocore.exceptions import ClientError


class InstanceManager:
    def __init__(self, settings):
        assert isinstance(settings, Settings)
        self.settings = settings

    @classmethod
    def make(cls, settings_file="settings.json"):
        try:
            return cls(Settings.load(settings_file))
        except SettingsError as e:
            raise BenchError("Failed to load settings", e)


    def hosts(self):
        # Construct the relative path to the ip.yaml file
        current_dir = os.path.dirname(os.path.abspath(__file__))
        yaml_file = os.path.join(current_dir, '../../../../../consenbench/assets/ip.yaml')

        # Read the YAML file
        with open(yaml_file, 'r') as file:
            data = yaml.safe_load(file)

        # Extract nodes' information
        nodes = data.get('nodes', [])
        ips = [node['Ip'] for node in nodes[1:]]
        print("available machines" + str(ips))
        return ips

    def user(self):
        # Construct the relative path to the ip.yaml file
        current_dir = os.path.dirname(os.path.abspath(__file__))
        yaml_file = os.path.join(current_dir, '../../../../../consenbench/assets/ip.yaml')

        # Read the YAML file
        with open(yaml_file, 'r') as file:
            data = yaml.safe_load(file)

        # Extract nodes' information
        nodes = data.get('nodes', [])

        return nodes[1]['Username']