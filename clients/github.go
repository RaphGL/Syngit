package clients

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/raphgl/syngit/config"
)

type GithubRepo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	PushedAt string `json:"pushed_at"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`
}

func getGithubRepos(cfg *config.Config) ([]GithubRepo, error) {
	resultsPerPage := 100 //max result for github
	page := 1
	resBody, err := getUserReposAPI("github", resultsPerPage, page, cfg)
	if err != nil {
		return nil, err
	}
	var repos []GithubRepo
	json.NewDecoder(resBody).Decode(&repos)

	// loop through for pagination
	for len(repos)%resultsPerPage == 0 {
		page++

		newResBody, err := getUserReposAPI("github", resultsPerPage, page, cfg)
		if err != nil {
			return nil, err
		}
		var newRepos []GithubRepo
		json.NewDecoder(newResBody).Decode(&newRepos)
		// prevent infinite loop if no results
		if len(newRepos) == 0 {
			break
		}
		repos = append(repos, newRepos...)
	}

	return repos, nil
}

func createRepoGitHub(repo GitRepo, cfg *config.Config) (GithubRepo, error) {
	var newRepo GithubRepo
	resBody, err := createRepoAPI("github", repo, cfg)
	if err != nil {
		return newRepo, err
	}

	json.NewDecoder(resBody).Decode(&newRepo)

	fmt.Println(fmt.Sprintf("Github repository name %s created successfully.", repo.GetName()))
	return newRepo, nil
}

func (gr *GithubRepo) GetName() string {
	return gr.Name
}

func (gr *GithubRepo) GetFullName() string {
	return gr.FullName
}

func (gr *GithubRepo) GetURL() string {
	return gr.CloneURL
}

func (gr *GithubRepo) IsPrivate() bool {
	return gr.Private
}

func (gr *GithubRepo) IsFork() bool {
	return gr.Fork
}

func (gr *GithubRepo) LastUpdated() (*time.Time, error) {
	tm, err := time.Parse("2006-01-02T15:04:05Z", gr.PushedAt)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (gr *GithubRepo) GetClientName() string {
	return "github"
}
