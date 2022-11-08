from pathlib import Path
from asyncio import subprocess
import asyncio
import aiohttp
import os
import toml

CLIENT_URL = {
    'github': 'github.com',
    'codeberg': 'codeberg.org'
}


class RepoSynchronizer:
    def __init__(self, main_client, data_path, config_toml):
        self.data_path = data_path
        self.main_client = main_client
        self.config_toml = config_toml

    async def __aenter__(self):
        self.session = aiohttp.ClientSession()
        self.clients = await self.__init_clients()
        return self

    async def __aexit__(self, *args):
        if not self.session.closed:
            await self.session.close()

    async def is_in_sync(self, repo_name):
        """
        Checks if other clients match main_client's commit
        """
        commits = dict()
        for client in self.clients:
            match client:
                case 'github':
                    async with self.session.get(f"https://api.github.com/repos/{repo_name}/commits") as response:
                        commits[client] = await response.json()
                case 'codeberg':
                    async with self.session.get(f"https://codeberg.org/api/v1/repos/{repo_name}/commits/") as response:
                        commits[client] = await response.json()

        for client in commits:
            try:
                if commits[client][0]['sha'] != commits[self.main_client][0]['sha']:
                    return False
            # KeyError if repo does not exist in one of the clients
            except KeyError:
                return False
        return True

    async def __init_clients(self):
        """
        Returns a dictionary with each of the clients repo names
        """
        git_clients = dict()

        # gets user name from toml config and uses it on API calls
        for git_client in self.config_toml:
            client = self.config_toml[git_client]
            if type(client) is dict:
                username = client['username']
                match git_client:
                    case 'github':
                        async with self.session.get(f'https://api.github.com/users/{username}/repos') as response:
                            git_clients[git_client] = await response.json()
                    case 'codeberg':
                        async with self.session.get(f'https://codeberg.org/api/v1/users/{username}/repos') as response:
                            git_clients[git_client] = await response.json()

        clients = dict()
        try:
            for client in git_clients:
                repos = []
                if len(git_clients[client]) > 0:
                    for i in range(len(git_clients[client])):
                        repos.append(git_clients[client][i]['full_name'])
                    clients[client] = repos
        except KeyError:
            print("You're being rate limited. Please try again later.")
        return clients

    async def clone_repo(self, repo_name):
        """
        Clones repos with repo_name from all the referenced clients and adds their remotes
        """
        Path(self.data_path).mkdir(parents=True, exist_ok=True)
        os.chdir(self.data_path)
        await asyncio.create_subprocess_shell(
            f"git clone git@{CLIENT_URL[self.main_client]}:{repo_name}.git",
            stdin=subprocess.PIPE, stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT
        )
        os.chdir(os.environ['HOME'])

        return self.clients

    async def synchronize(self):
        """
        Synchronizes all the known client with the repos in the main client
        """
        sync_info = []
        try:
            for repo in self.clients[self.main_client]:
                is_synced_task = asyncio.create_task(
                    self.is_in_sync(repo))
                sync_info.append(is_synced_task)

            # retrieve and clone unsynced repos
            async_info = await asyncio.gather(*sync_info)

            for index, repo in enumerate(self.clients[self.main_client]):
                if not async_info[index]:
                    await self.clone_repo(repo)

            for index, repo in enumerate(self.clients[self.main_client]):
                if not async_info[index]:
                    # only push to repo if repo exists in main client and on target client
                    self.__push_to_repo(repo)
        except:
            return

    def __push_to_repo(self, repo):
        os.chdir(f"{self.data_path}/{repo.split('/')[1]}")
        os.system(f"git pull")
        for client in self.clients:
            if client != self.main_client and repo in self.clients[client]:
                print(f"pushing repo {repo}")
                os.system(
                    f"git remote add {client} git@{CLIENT_URL[client]}:{repo}")
                os.system(f"git push -f --all {client}")


async def main():
    # TODO read toml file from somewhere not inside this repo
    with open('syngit.toml') as f:
        toml_config = toml.load(f)

    syngit_data_path = f"{os.environ['HOME']}/.local/share/syngit"
    if not os.path.exists(syngit_data_path):
        os.mkdir(syngit_data_path)

    async with RepoSynchronizer(
            main_client=toml_config['main_client'],
            config_toml=toml_config,
            data_path=syngit_data_path) as repo_sync:
        await repo_sync.synchronize()


if __name__ == '__main__':
    asyncio.run(main())
