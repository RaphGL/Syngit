package clients

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/raphgl/syngit/config"
)

type GitRepo interface {
	// Returns the name of the repository
	GetName() string
	// Returns the name of the repository in the format "Username/Repo"
	GetFullName() string
	// Returns the Git URL for the repository
	GetURL() string
	IsPrivate() bool
	IsFork() bool
}

func CloneRepo(r GitRepo, cfg *config.Config) error {
	repoURL := r.GetURL()
	cachePath, err := cfg.GetRepoCachePath()
	if err != nil {
		return err
	}

	// appends directory for repo as otherwise the current directory is turned into a git repo
	_, err = git.PlainClone(cachePath+r.GetName(), false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return nil
}
