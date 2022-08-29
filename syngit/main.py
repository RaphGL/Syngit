from sync import RepoSynchronizer
import asyncio
import toml
import os


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
