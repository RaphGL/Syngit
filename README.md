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
    <a href="https://github.com/RaphGL/ProjectName"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    ·
    <a href="https://github.com/RaphGL/ProjectName/issues">Report Bug</a>
    ·
    <a href="https://github.com/RaphGL/ProjectName/issues">Request Feature</a>
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

Syngit is a CLI and init service that let's you synchronize repositories across different clients (Github, Codeberg, Gitlab, etc) with a very simple interface and configuration file.

### Built With

- [Python](https://www.python.org/)
- [Poetry](https://python-poetry.org/)
- [AioHttp](https://docs.aiohttp.org/en/stable/)
- [TOML](https://github.com/uiri/toml)

<!-- GETTING STARTED -->

## Getting Started

To get Syngit to work you need to use SSH login on the clients you wish to synchronize. Read [this article on how to set it up](https://docs.github.com/en/authentication/connecting-to-github-with-ssh) if you're uncertain.

### Installation

TODO

<!-- USAGE EXAMPLES -->

## Usage

1. Create a `$HOME/.config/syngit.toml` file
2. Make something akin to this:

```toml
main_client = "github"

[github]
username = "RaphGL"

[codeberg]
username = "RaphGL"

[gitlab]
username = "RaphGL"
```

3. Enable the syngit service:

```sh
$ systemctl enable syngit --now
```

<!-- LICENSE -->

## License

Distributed under LICENSE License. See [`LICENSE`](https://github.com/RaphGL/Syngit/blob/main/LICENSE) for more information.

<!-- ACKNOWLEDGEMENTS -->
<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/othneildrew/Best-README-Template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/othneildrew/Best-README-Template/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/othneildrew/Best-README-Template/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/othneildrew/Best-README-Template/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png

