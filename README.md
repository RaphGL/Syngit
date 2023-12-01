<!-- PROJECT LOGO -->
<br />
<p align="center">
<!--  
  <a href="https://github.com/RaphGL/Syngit">
    <img src="logo.png" alt="Logo" height="80">
  </a> --->

  <h1 align="center">Syngit</h3>
  <p align="center">Synchronize repositories across Git clients</p>
  <p align="center">
    <br />
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

Syngit is a CLI and daemon that lets you synchronize repositories across different clients (Github, Codeberg, Gitlab, etc) with a very simple configuration file.
Syngit is a tool with no external dependencies (not even git). It contains everything you need, just compile it and run it.

### Built With

- [Go](https://go.dev/)
- [go-git](https://github.com/go-git/go-git)

<!-- GETTING STARTED -->

## Getting Started

To use Syngit, make sure you got the tokens/passwords needed to authenticate for your target Git clients (ie Github, GitLab, Codeberg, etc).
After you've gotten your tokens/passwords you need a `syngit.toml` in your system's [default config directory](https://pkg.go.dev/os#UserConfigDir). The config file has the following structure:

```toml
# your main git client
main_client = "github"
# ignores all files that match the glob pattern, unimplemented!
glob_ignore = ["*cpp"]
# where the cache for syngit should be stored, defaults to https://pkg.go.dev/os#UserCacheDir
cache_dir = "~/Documents/Test"

[client.codeberg]
username = "RaphGL"
token = "my_token"

[client.github]
username = "RaphGL"
token = "my_token"
# temporarily disable synchronization to this client
disable = true

[client.gitlab]
username = "RaphGL"
token = "my_token"
# repositories to be ignored on this client
ignore = ["repo1", "repo2"]
```


### Installation

```sh
$ git clone https://github.com/RaphGL/Syngit
$ cd Syngit
$ go build
```

<!-- USAGE EXAMPLES -->

## Usage

1. Create a `$HOME/.config/syngit.toml` file
2. Fill out the configuration file 
3. Enable the syngit service (WIP):

```sh
$ systemctl enable syngit --now
```

<!-- LICENSE -->

## License

Distributed under GPLv3 License. See [`LICENSE`](https://github.com/RaphGL/Syngit/blob/main/LICENSE) for more information.
