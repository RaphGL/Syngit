import asyncio
from asyncio import subprocess
from multiprocessing import synchronize
import aiohttp
import toml
import os
from pathlib import Path

MAIN_CLIENT = ""
CLIENT_URL = {
    'github': 'https://github.com',
    'codeberg': 'https://codeberg.org'
}


async def is_in_sync(session, clients, repo_name):
    """
    Checks if other clients match main_client's commit
    """
    commits = dict()
    for client in clients:
        match client:
            case 'github':
                async with session.get(f"https://api.github.com/repos/{repo_name}/commits") as response:
                    commits[client] = await response.json()
            case 'codeberg':
                async with session.get(f"https://codeberg.org/api/v1/repos/{repo_name}/commits/") as response:
                    commits[client] = await response.json()

    for client in commits:
        try:
            if commits[client][0]['sha'] != commits[MAIN_CLIENT][0]['sha']:
                return False
        # KeyError if repo does not exist in one of the clients
        except KeyError:
            return False
    return True


async def init_clients(session, toml_config):
    """
    Returns a dictionary with each of the clients repo names
    """
    git_clients = dict()

    # gets user name from toml config and uses it on API calls
    for git_client in toml_config:
        client = toml_config[git_client]
        if type(client) is dict:
            username = client['username']
            match git_client:
                case 'github':
                    async with session.get(f'https://api.github.com/users/{username}/repos') as response:
                        git_clients[git_client] = await response.json()
                case 'codeberg':
                    async with session.get(f'https://codeberg.org/api/v1/users/{username}/repos') as response:
                        git_clients[git_client] = await response.json()

    clients = dict()
    for client in git_clients:
        repos = []
        for i in range(len(git_clients[client])):
            repos.append(git_clients[client][i]['full_name'])
        clients[client] = repos
    return clients


async def clone_repo(repos_dir, clients, repo_name):
    """
    Clones repos with repo_name from all the referenced clients and adds their remotes
    """
    # TODO add check if repo folder already exists
    # TODO add pipe in case repo is private using toml's password
    for client in clients:
        client_path = f"{repos_dir}/{client}"
        Path(client_path).mkdir(parents=True, exist_ok=True)
        os.chdir(client_path)
        await asyncio.create_subprocess_shell(f"git clone {CLIENT_URL[client]}/{repo_name}", stdin=subprocess.PIPE, stdout=subprocess.DEVNULL, stderr=subprocess.STDOUT)
    os.chdir(os.environ['HOME'])

    return clients


async def synchronize_repos(session, clients, syngit_data_path):
    sync_info = []
    for repo in clients[MAIN_CLIENT]:
        is_synced_task = asyncio.create_task(
            is_in_sync(session, clients, repo))
        sync_info.append(is_synced_task)

    async_info = asyncio.gather(*sync_info)
    for index, repo in enumerate(clients[MAIN_CLIENT]):
        if async_info[index]:
            clone_repo(syngit_data_path, clients, repo)


async def main():
    # TODO read toml file from somewhere not inside this repo
    with open('syngit.toml') as f:
        toml_config = toml.load(f)

    global MAIN_CLIENT
    MAIN_CLIENT = toml_config['main_client']

    syngit_data_path = f"{os.environ['HOME']}/.local/share/syngit"
    if not os.path.exists(syngit_data_path):
        os.mkdir(syngit_data_path)

    async with aiohttp.ClientSession() as session:
        clients = await init_clients(session, toml_config)
        await synchronize_repos(session, clients, syngit_data_path)


if __name__ == '__main__':
    asyncio.run(main())
