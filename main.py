import asyncio
import aiohttp
import tempfile
import toml
import os


async def main():
    f = open('syngit.toml')
    toml_config = toml.load(f)
    f.close()

    async with aiohttp.ClientSession() as session:
        await init(session, toml_config)
            

async def is_updated(session):
    """
    Checks if other clients match main_client's commit
    """
    pass

async def init(session, toml_config):
    """
    Creates the needed repository for syncing
    TODO: use .local instead of /tmp and only clone if repo's are non existant
    # TODO use Popen to run multiple git clones concurrently
    """
    git_clients = dict()

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

    tempfolder = tempfile.mkdtemp(prefix='syngit_')

    for client in git_clients:
        client_path = os.path.join(tempfolder, client)
        os.mkdir(client_path)
        os.chdir(client_path)
        for i in range(len(git_clients[client])):
            repo = git_clients[client][i]['html_url']
            # TODO use Popen to run multiple git clones concurrently
            os.system(f"git clone {repo}")


if __name__ == '__main__':
    asyncio.run(main())
