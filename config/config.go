package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		// a glob to ignore certain repositories through a pattern
		MainClient string `toml:"main_client"`
		// DO NOT USE DIRECTLY
		// GetRepoCachePath() provides a fallback, use that instead
		CacheDir   string            `toml:"cache_dir"`
		GlobIgnore []string          `toml:"glob_ignore"`
		Client     map[string]Client `toml:"client"`
		// oldest age allowed for repo before deletion
		// measured in days
		MaxAge       int  `toml:"max_age"`
		IncludeForks bool `toml:"include_forks"`
	}

	Client struct {
		Username string `toml:"username"`
		// authentication token or password for the client
		Token string `toml:"token"`
		// disable syncing this client
		Disable bool `toml:"disable"`
		// repositories to be ignored
		Ignore []string `toml:"ignore"`
	}
)

func LoadConfig() (*Config, error) {
	var cfg Config

	path, err := GetDefaultConfigPath()
	if err != nil {
		return nil, err
	}

	// uses absolute file paths to remove ambiguity
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	// making sure required fields are filled
	failedCheck := false
	if len(cfg.MainClient) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Did not specified what is the `main_client`")
		failedCheck = true
	}

	for clientName, client := range cfg.Client {
		if len(client.Username) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Missing username for", clientName)
			failedCheck = true
		}
	}

	if failedCheck {
		os.Exit(1)
	}

	return &cfg, nil
}

func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	homeConfig := filepath.Join(homeDir, "syngit.toml")
	cfgConfig := filepath.Join(cfgDir, "syngit.toml")

	// returns the first config file that exists
	// the cfgConfig is checked first, this means it complies with XDG_CONFIG
	// but falls back to home config in case the user wants it there
	for _, f := range [...]string{cfgConfig, homeConfig} {
		_, err := os.Stat(f)
		if err == nil {
			return f, nil
		}
	}

	return "", err
}

func (c *Config) GetRepoCachePath() (string, error) {
	// TODO: expand tilde and make paths absolute
	if c.CacheDir == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(cacheDir, "syngit"), nil
	}

	return c.CacheDir, nil
}
