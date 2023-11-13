package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raphgl/syngit/config"
)

type GithubRepo struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	CloneURL  string `json:"clone_url"`
	UpdatedAt string `json:"pushed_at"`
	Private   bool   `json:"private"`
	Fork      bool   `json:"fork"`
}

func getGithubRepos(cfg *config.Config) ([]GithubRepo, error) {
	APIPoint := fmt.Sprintf("https://api.github.com/users/%s/repos", cfg.Client["github"].Username)
	res, err := http.Get(APIPoint)
	if err != nil {
		return nil, err
	}

	var repos []GithubRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil
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
	tm, err := time.Parse("2006-01-02T15:04:05Z", gr.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (gr *GithubRepo) GetClientName() string {
	return "github"
}
