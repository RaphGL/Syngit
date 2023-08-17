package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GithubRepo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`
}

func GetGithubRepos() ([]GithubRepo, error) {
	APIPoint := fmt.Sprintf("https://api.github.com/users/%s/repos", user)
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
