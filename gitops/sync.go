package gitops

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

func pullChangesFromMainRepo(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
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

		fmt.Println("INFO: Updating", repoPath, "for", remoteName)
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

func SyncMirrors(rm clients.GitRepoMap, cfg *config.Config) {
	localRepos, err := GetLocalRepoPaths(cfg)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	for _, r := range localRepos {
		if err = pullChangesFromMainRepo(r); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if err = pushToClientRepo(r, cfg); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
