<!-- PROJECT LOGO -->
<br />
<p align="center">
<!--  
  <a href="https://github.com/RaphGL/Syngit">
    <img src="logo.png" alt="Logo" height="80">
  </a> --->

  <h3 align="center">A simple to use repo synchronization tool</h3>
  <p align="center">
    <br />
    <a href="https://github.com/RaphGL/Syngit"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    ·
    <a href="https://github.com/RaphGL/Syngit/issues">Report Bug</a>
    ·
    <a href="https://github.com/RaphGL/Syngit/issues">Request Feature</a>
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

To get Syngit to work you need to use SSH login on the clients you wish to synchronize. Read [this article on how to set it up](https://docs.github.com/en/authentication/connecting-to-github-with-ssh) if you're uncertain.
You need a `syngit.toml` in your system's [default config directory](https://pkg.go.dev/os#UserConfigDir). The config file has the following structure:

```toml
# ignores all files that match the glob pattern, unimplemented!
glob_ignore = ["*cpp"]
# your main git client
main_client = "github"
# where the cache for syngit should be stored, defaults to https://pkg.go.dev/os#UserCacheDir
cache_dir = "~/Documents/Test"

[client.codeberg]
username = "RaphGL"
token = "my_token"

[client.github]
username = "RaphGL"
token = "my_token"
# temporarily disable client
disable = true

[client.gitlab]
username = "RaphGL"
token = "my_token"
# repositories to be ignored in client
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
3. Enable the syngit service:

```sh
$ systemctl enable syngit --now
```

<!-- LICENSE -->

## License

Distributed under GPLv3 License. See [`LICENSE`](https://github.com/RaphGL/Syngit/blob/main/LICENSE) for more information.
