from asyncio import subprocess
from config import CLIENT_URL
from pathlib import Path
import asyncio
import aiohttp
import os


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
        for client in git_clients:
            repos = []
            for i in range(len(git_clients[client])):
                repos.append(git_clients[client][i]['full_name'])
            clients[client] = repos
        return clients

    async def clone_repo(self, repo_name):
        """
        Clones repos with repo_name from all the referenced clients and adds their remotes
        """
        # TODO add check if repo folder already exists
        # TODO add pipe in case repo is private using toml's password
        Path(self.data_path).mkdir(parents=True, exist_ok=True)
        os.chdir(self.data_path)
        await asyncio.create_subprocess_shell(
            f"git clone {CLIENT_URL[self.main_client]}/{repo_name}",
            stdin=subprocess.PIPE, stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT
        )
        os.chdir(os.environ['HOME'])

        return self.clients

    async def synchronize(self):
        sync_info = []
        for repo in self.clients[self.main_client]:
            is_synced_task = asyncio.create_task(
                self.is_in_sync(repo))
            sync_info.append(is_synced_task)

        # retrieve and clone unsynced repos
        async_info = await asyncio.gather(*sync_info)

        for index, repo in enumerate(self.clients[self.main_client]):
            if async_info[index]:
                await self.clone_repo(repo)

        for index, repo in enumerate(self.clients[self.main_client]):
            if async_info[index]:
                await self.__push_to_repo(repo)

    async def __push_to_repo(self, repo):
        # TODO add authentication to be able to pull to repo
        os.chdir(f"{self.data_path}/{repo.split('/')[1]}")
        os.system(f"git pull")
        for client in self.clients:
            if client != self.main_client:
                os.system(
                    f"git remote add {client} {CLIENT_URL[client]}/{repo}")

                push_cmd = await subprocess.create_subprocess_shell(
                    "git push {client} -f",
                    stdin=subprocess.PIPE, stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL)
