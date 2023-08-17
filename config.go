package main

type (
	Config struct {
		GlobIgnore []string          `toml:"glob_ignore"`
		Client     map[string]Client `toml:"client"`
	}

	Client struct {
		Token  string   `toml:"token"`
		Enable bool     `toml:"enable"`
		Ignore []string `toml:"ignore"`
	}
)
