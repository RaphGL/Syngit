package clients

import (
	"bytes"
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
	APIPoint := "https://api.github.com/user/repos"
	client := &http.Client{}

	req, err := http.NewRequest("GET", APIPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+cfg.Client["github"].Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []GithubRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil
}

func createRepoGitHub(repo GitRepo, cfg *config.Config) error {
	APIPoint := "https://api.github.com/user/repos"
	client := &http.Client{}

	payload := map[string]interface{}{
		"name":    repo.GetName(),
		"private": repo.IsPrivate(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", APIPoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+cfg.Client["github"].Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("Failed to create GitHub repository. Status code: %d", res.StatusCode)
	}

	fmt.Println(fmt.Sprintf("GitHub repository name %s created successfully.", repo.GetName()))

	return nil
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
