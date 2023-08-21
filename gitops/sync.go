package gitops

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

// TODO: use fetch to get commits for non main client remotes
// this will be used on an alternative syncing strategy which allows
// for people to commit from different clients and have their changes reflected

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

func SyncMirrors(rm clients.GitRepoMap, cfg *config.Config) {
	localRepos, err := GetLocalRepoPaths(cfg)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	for _, r := range localRepos {
		err := pullChangesFromMainRepo(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
