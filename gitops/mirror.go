package gitops

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

func cloneRepo(r clients.GitRepo, cfg *config.Config, repoPath string) {
	repoURL := r.GetURL()

	// appends directory for repo as otherwise the current directory is turned into a git repo
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: cfg.Client[cfg.MainClient].Username,
			Password: cfg.Client[cfg.MainClient].Token,
		},
	})

	if err != nil {
		slog.Error(r.GetName(), err)
	}

	slog.Info("Cloning " + r.GetURL())
}

// creates a new remote for client in repoPath
func AddMirrorAsRemote(m clients.GitRepo, repoPath string) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	remoteName := m.GetClientName()

	_, err = r.Remote(remoteName)
	// don't do anything if remote exists already
	if err == nil {
		return
	}

	_, err = r.CreateRemote(&gitCfg.RemoteConfig{
		Name: remoteName,
		URLs: []string{m.GetURL()},
	})

	if err != nil {
		slog.Warn("Failed to add remote " + remoteName)
		return
	}

	slog.Info(fmt.Sprintf("Added %s remote to %s", remoteName, repoPath))
}

func CreateLocalMirrors(m clients.GitRepoMap, cfg *config.Config) error {
	cachePath, err := cfg.GetRepoCachePath()
	if err != nil {
		return err
	}

	// clone repo from main client
	for _, v := range m {
		for _, r := range v {
			repoPath := filepath.Join(cachePath, r.GetName())

			if strings.Contains(r.GetURL(), cfg.MainClient) {
				isOlder, err := RepoIsOlderThanSpecified(r, cfg)
				if err != nil {
					return err
				}

				if isOlder {
					continue
				}

				cloneRepo(r, cfg, repoPath)

				// add remote for mirrors in the main client's repo
				for _, r := range v {
					repoPath := filepath.Join(cachePath, r.GetName())
					AddMirrorAsRemote(r, repoPath)
				}
			}
		}
	}

	return nil
}
