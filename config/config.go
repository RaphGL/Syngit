package config

import (
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		// a glob to ignore certain repositories through a pattern
		GlobIgnore []string          `toml:"glob_ignore"`
		MainClient string            `toml:"main_client"`
		Client     map[string]Client `toml:"client"`
	}

	Client struct {
		Username string `toml:"username"`
		// authentication token or password for the client
		Token string `toml:"token"`
		// disable syncing this client
		Enable bool `toml:"enable"`
		// repositories to be ignored
		Ignore []string `toml:"ignore"`
	}
)

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
