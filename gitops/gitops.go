package gitops

import (
	"os"
	"path/filepath"

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
		return nil, err
	}

	var paths []string
	for _, r := range repos {
		repoPath := filepath.Join(cachePath, r.Name())
		paths = append(paths, repoPath)
	}

	return paths, nil
}
