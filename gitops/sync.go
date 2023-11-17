package gitops

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

func pullChangesFromRepo(repoPath string, cfg *config.Config) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: cfg.Client[cfg.MainClient].Username,
			Password: cfg.Client[cfg.MainClient].Token,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

//// TODO: use fetch to get commits for non main client remotes
//// this will be used on an alternative syncing strategy which allows
//// for people to commit from different clients and have their changes reflected
// func fetchInfoForRemotes(repoPath string) error {
// 	r, err := git.PlainOpen(repoPath)
// 	if err != nil {
// 		return err
// 	}
//
// 	remotes, err := r.Remotes()
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, remote := range remotes {
// 		remoteCfg := remote.Config()
// 		if remoteCfg.Name == "origin" {
// 			continue
// 		}
//
// 		err = r.Fetch(&git.FetchOptions{
// 			RemoteName: remoteCfg.Name,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

func pushToClientRepo(repoPath string, cfg *config.Config) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return err
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
			return err
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
		if err = pullChangesFromRepo(r, cfg); err != nil {
            slog.Error(err.Error())
		}
		if err = pushToClientRepo(r, cfg); err != nil {
            slog.Error(err.Error())
		}
	}

	return nil
}
