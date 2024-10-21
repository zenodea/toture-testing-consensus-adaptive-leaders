from json import load, JSONDecodeError


class SettingsError(Exception):
    pass


class Settings:
    def __init__(self, key_name, key_path, consensus_port, mempool_port,
                 front_port, repo_name, repo_url, branch):

        self.key_name = key_name
        self.key_path = key_path
        self.consensus_port = consensus_port
        self.mempool_port = mempool_port
        self.front_port = front_port
        self.repo_name = repo_name
        self.repo_url = repo_url
        self.branch = branch

    @classmethod
    def load(cls, filename):
        try:
            with open(filename, 'r') as f:
                data = load(f)
            print("settings: "+ str(data))
            return cls(
                data['key']['name'],
                data['key']['path'],
                data['ports']['consensus'],
                data['ports']['mempool'],
                data['ports']['front'],
                data['repo']['name'],
                data['repo']['url'],
                data['repo']['branch'],
            )
        except (OSError, JSONDecodeError) as e:
            raise SettingsError(str(e))

        except KeyError as e:
            raise SettingsError(f'Malformed settings: missing key {e}')
