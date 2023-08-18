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

// returns a map where the key is the name of the repo and
// the value is a slice of repos, where each repo belongs to a different client
func GetRepos(cfg *config.Config) map[string][]GitRepo {
	var (
		githubRepo       []GithubRepo
		gitlabRepo       []GitlabRepo
		codebergRepo     []CodebergRepo
		err1, err2, err3 error
	)

	if !cfg.Client["github"].Disable {
		githubRepo, err1 = getGithubRepos(cfg)
	}
	if !cfg.Client["gitlab"].Disable {
		gitlabRepo, err2 = getGitlabRepos(cfg)
	}
	if !cfg.Client["codeberg"].Disable {
		codebergRepo, err3 = getCodebergRepos(cfg)
	}

	erroredOut := func(errs ...error) bool {
		for _, err := range errs {
			if err != nil {
				return true
			}
		}
		return false
	}

	if erroredOut(err1, err2, err3) {
		return nil
	}

	// and a slice of all the repos that match said key
	// aka map[string][]GitClient
	repos := make(map[string][]GitRepo)

	// note: capturing is required here, otherwise the r will be mutated into the next value
	// resulting in all the items being the same
	for _, r := range githubRepo {
		func(r GithubRepo) {
			key := r.GetName()
			repos[key] = append(repos[key], &r)
		}(r)
	}
	for _, r := range gitlabRepo {
		func(r GitlabRepo) {
			key := r.GetName()
			repos[key] = append(repos[key], &r)
		}(r)
	}
	for _, r := range codebergRepo {
		func(r CodebergRepo) {
			key := r.GetName()
			repos[key] = append(repos[key], &r)
		}(r)
	}

	return repos
}
