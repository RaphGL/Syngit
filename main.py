from distutils.log import debug
from bs4 import BeautifulSoup
import asyncio
import aiohttp
import tempfile
import json
import os


async def main():
    async with aiohttp.ClientSession() as session:
        f = open('syngit.json')
        json_config = json.load(f)
        f.close()

        git_clients = dict()

        for git_client in json_config:
            username = json_config[git_client]
            match git_client:
                case 'github':
                    async with session.get(f'https://github.com/{username}?tab=repositories') as response:
                        git_clients[git_client] = get_github_repos(await response.text())
                case 'codeberg':
                    async with session.get(f'https://codeberg.org/{username}') as response:
                        git_clients[git_client] = get_codeberg_repos(await response.text())

        tempfolder = tempfile.mkdtemp(prefix='syngit_')
        for client in git_clients:
            client_path = os.path.join(tempfolder, client)
            os.mkdir(client_path)
            os.chdir(client_path)
            for repo in git_clients[client]:
                # TODO use Popen to run multiple git clones concurrently
                os.system(f"git clone {repo}")


def get_github_repos(repos_page):
    soup = BeautifulSoup(repos_page, 'html.parser')
    repos = []
    for repo in soup.find_all('a', attrs={'itemprop': 'name codeRepository'}):
        repos.append(f"https://github.com{repo['href']}")
    return repos


def get_codeberg_repos(repos_page):
    soup = BeautifulSoup(repos_page, 'html.parser')
    repos = []
    for anchor in soup.select('div.repo-title a.name'):
        repos.append(f"https://codeberg.org{anchor['href']}")
    return repos


if __name__ == '__main__':
    asyncio.run(main())
