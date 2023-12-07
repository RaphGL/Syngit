package gitops

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

func pullChangesFromRepo(repoPath string, cfg *config.Config) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []gitCfg.RefSpec{"refs/*:refs/*"},
	})
	if err != nil {
		return err
	}

	// retrieve branches in repo
	refs, err := r.Branches()
	if err != nil {
		return err
	}
	defer refs.Close()
	branches := make([]plumbing.ReferenceName, 3)
	refs.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, r.Name())
		return nil
	})

	// pull changes from every branch
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	for _, branchRef := range branches {
		err = w.Checkout(&git.CheckoutOptions{Branch: branchRef})
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{
			RemoteName: cfg.MainClient,
			Auth: &http.BasicAuth{
				Username: cfg.Client[cfg.MainClient].Username,
				Password: cfg.Client[cfg.MainClient].Token,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func pushToClientRepo(repoPath string, cfg *config.Config) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return err
	}

	refs, err := r.Branches()
	if err != nil {
		return err
	}
	defer refs.Close()
	branches := make([]plumbing.ReferenceName, 3)
	refs.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, r.Name())
		return nil
	})

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	for _, branchRef := range branches {
		err = w.Checkout(&git.CheckoutOptions{Branch: branchRef})
		if err != nil {
			// TODO: find out why checkout doesn't find certain references
			continue
		}

		for _, remote := range remotes {
			remoteName := remote.Config().Name
			currClient := cfg.Client[remoteName]
			if remoteName == cfg.MainClient || remoteName == "origin" || currClient.Disable {
				continue
			}

			slog.Info(fmt.Sprintf("Updating %s for %s", repoPath, remoteName))
			err = r.Push(&git.PushOptions{
				RemoteName: remoteName,
				Auth: &http.BasicAuth{
					Username: currClient.Username,
					Password: currClient.Token,
				},
			})
			if err != nil {
				slog.Error(err.Error())
			}
		}

	}

	return nil
}

func repoExists(repo []clients.GitRepo, clientName string) bool {
	repoExists := false
	for _, client := range repo {
		if client.GetClientName() == clientName {
			repoExists = true
		}
	}

	return repoExists
}

func CreateRepos(repoPath string, cfg *config.Config, m clients.GitRepoMap) error {
	repoName := filepath.Base(repoPath)
	repos := m[repoName]
	var mainRepo clients.GitRepo
	for _, repo := range repos {
		if repo.GetClientName() == cfg.MainClient {
			mainRepo = repo
		}
	}

	if mainRepo == nil {
		slog.Error("Failed to find a main client")
		return nil
	}

	for clientName := range cfg.Client {
		if !repoExists(repos, clientName) && !cfg.Client[clientName].Disable && cfg.Client[clientName].Create &&
			cfg.MainClient != clientName {
			clients.CreateRepo(mainRepo, clientName, cfg, &m)
			AddMirrorAsRemote(m[repoName][len(m[repoName])-1], repoPath)
		}
	}

	return nil
}

func SyncMirrors(rm clients.GitRepoMap, cfg *config.Config) error {
	localRepos, err := GetLocalRepoPaths(cfg)
	if err != nil {
		return err
	}

	for _, r := range localRepos {
		if err = CreateRepos(r, cfg, rm); err != nil {
			slog.Error(err.Error())
		}
		if err = pullChangesFromRepo(r, cfg); err != nil {
			slog.Error(err.Error())
		}
		if err = pushToClientRepo(r, cfg); err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}
