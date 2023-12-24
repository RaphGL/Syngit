package clients

import (
	"log/slog"
	"time"

	"github.com/raphgl/syngit/config"
)

type GitRepo interface {
	// Returns the name of the repository
	GetName() string
	// Returns the name of the repository in the format "Username/Repo"
	GetFullName() string
	// Returns the Git URL for the repository
	GetURL() string
	// Returns whether the repository is private or not
	IsPrivate() bool
	// Returns whether the repository is a fork of another one
	IsFork() bool
	// Returns the last time the repository was committed to
	LastUpdated() (*time.Time, error)
	// Returns the name of the Git client/frontend
	GetClientName() string
}

// a map where the key is the name of the repo and the value is the slice of all repos with that name (one for each git client)
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

func CreateRepo(repo GitRepo, clientName string, cfg *config.Config, repoMap *GitRepoMap) {
	switch clientName {
	case "github":
		newGithubRepo, err := createRepoGitHub(repo, cfg)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		addRepoToMap(repoMap, &newGithubRepo, cfg)
	case "gitlab":
		newGitlabRepo, err := createRepoGitlab(repo, cfg)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		addRepoToMap(repoMap, &newGitlabRepo, cfg)
	case "codeberg":
		newCodebergRepo, err := createRepoCodeberg(repo, cfg)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		addRepoToMap(repoMap, &newCodebergRepo, cfg)
	}

}
