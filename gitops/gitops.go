package gitops

import (
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

// returns repos in mirror cache
func GetLocalRepoPaths(cfg *config.Config) ([]string, error) {
	cachePath, err := cfg.GetRepoCachePath()
	if err != nil {
		return nil, err
	}

	repos, err := os.ReadDir(cachePath)
	if err != nil {
		err = os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	var paths []string
	for _, r := range repos {
		if r.IsDir() {
			repoPath := filepath.Join(cachePath, r.Name())
			paths = append(paths, repoPath)
		}
	}

	return paths, nil
}

func RepoIsOlderThanSpecified(repo clients.GitRepo, cfg *config.Config) (bool, error) {
	lastUpdated, err := repo.LastUpdated()
	if err != nil {
		return false, err
	}

	var days int
	if reflect.ValueOf(cfg.MaxAge).IsZero() {
		days = -30 * 6
	} else {
		days = -cfg.MaxAge
	}

	timeframe := time.Now().AddDate(0, 0, days)
	if lastUpdated.After(timeframe) {
		return false, nil
	} else {
		return true, nil
	}
}
