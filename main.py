import asyncio
import aiohttp
import toml

MAIN_CLIENT = ""


async def main():
    with open('syngit.toml') as f:
        toml_config = toml.load(f)

    global MAIN_CLIENT
    MAIN_CLIENT = toml_config['main_client']

    async with aiohttp.ClientSession() as session:
        clients = await init_clients(session, toml_config)
        for repo in clients[MAIN_CLIENT]:
            if not is_in_sync(session, clients, repo):
                clone_repo(clients, repo)

def clone_repo(clients, repo_name):
    """
    Clones repos with repo_name from all the referenced clients and adds their remotes
    """
    pass


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
                    commits['codeberg'] = await response.json()

    for client in commits:
        if commits[client][0]['sha'] != commits[MAIN_CLIENT][0]['sha']:
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


if __name__ == '__main__':
    asyncio.run(main())
