# Copyright(C) Facebook, Inc. and its affiliates.
import boto3
from botocore.exceptions import ClientError
from collections import defaultdict, OrderedDict
from time import sleep

from benchmark.utils import Print, BenchError, progress_bar
from benchmark.settings import Settings, SettingsError
import yaml
import os


class AWSError(Exception):
    def __init__(self, error):
        assert isinstance(error, ClientError)
        self.message = error.response["Error"]["Message"]
        self.code = error.response["Error"]["Code"]
        super().__init__(self.message)


class InstanceManager:
    INSTANCE_NAME = "dag-node"
    SECURITY_GROUP_NAME = "dag"

    def __init__(self, settings):
        assert isinstance(settings, Settings)
        self.settings = settings
        self.clients = OrderedDict()

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
        user = nodes[1]['Username']
        print("user" + str(user))
        return user

