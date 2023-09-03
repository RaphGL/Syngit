package clients

import (
	"time"

	"github.com/raphgl/syngit/config"
)

// TODO: URGENT: add authentication to be able to pull
// private repos and push changes to client
type GitRepo interface {
	// Returns the name of the repository
	GetName() string
	// Returns the name of the repository in the format "Username/Repo"
	GetFullName() string
	// Returns the Git URL for the repository
	GetURL() string
	IsPrivate() bool
	IsFork() bool
	LastUpdated() (*time.Time, error)
	GetClientName() string
}

type GitRepoMap = map[string][]GitRepo

func addRepoToMap(repos *GitRepoMap, r GitRepo, cfg *config.Config) {
	key := r.GetName()
	switch r.IsFork() {
	case true:
		if cfg.IncludeForks {
			(*repos)[key] = append((*repos)[key], r)
		}
	case false:
		if !cfg.IncludeForks {
			(*repos)[key] = append((*repos)[key], r)

		}
	}
}

// returns a map where the key is the name of the repo and
// the value is a slice of repos, where each repo belongs to a different client
// note: this function only returns repos not disabled by config
func GetRepos(cfg *config.Config) GitRepoMap {
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
	repos := make(GitRepoMap)

    // note: the r := r is required otherwise range will keep adding to the r pointer
    // resulting in all values in map pointing to the same repo
	for _, r := range githubRepo {
		r := r
		addRepoToMap(&repos, &r, cfg)
	}
	for _, r := range gitlabRepo {
		r := r
		addRepoToMap(&repos, &r, cfg)
	}
	for _, r := range codebergRepo {
		r := r
		addRepoToMap(&repos, &r, cfg)
	}


	return repos
}
