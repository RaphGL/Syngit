package gitops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

func cloneRepo(r clients.GitRepo, cfg *config.Config, repoPath string) error {
	repoURL := r.GetURL()

	// appends directory for repo as otherwise the current directory is turned into a git repo
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return nil
}

func addMirrorAsRemote(m clients.GitRepo, repoPath string) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	remoteName := m.GetClientName()

	_, err = r.CreateRemote(&gitCfg.RemoteConfig{
		Name: remoteName,
		URLs: []string{m.GetURL()},
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to add remote:", remoteName)
		return
	}

	fmt.Println("INFO: Added", remoteName, "remote to", repoPath)
}

func CreateLocalMirrors(m clients.GitRepoMap, cfg *config.Config) {
	cachePath, err := cfg.GetRepoCachePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// clone repo from main client
	for _, v := range m {
		for _, r := range v {
			repoPath := filepath.Join(cachePath, r.GetName())

			if strings.Contains(r.GetURL(), cfg.MainClient) {
				cloneRepo(r, cfg, repoPath)

				// add remote for mirrors in the main client's repo
				for _, r := range v {
					repoPath := filepath.Join(cachePath, r.GetName())
					addMirrorAsRemote(r, repoPath)
				}
			}

		}
	}

}
