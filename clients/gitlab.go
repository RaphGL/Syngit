package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raphgl/syngit/config"
)

type GitlabRepo struct {
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility        string `json:"visibility"`
	HttpURLToRepo     string `json:"http_url_to_repo"`
	UpdatedAt         string `json:"last_activity_at"`
	ForkedFromProject *any   `json:"forked_from_project"`
}

func getGitlabRepos(cfg *config.Config) ([]GitlabRepo, error) {
	resultsPerPage := 100 //max result for gitlab
	page := 1
	repos, err := requestGitlabRepos(resultsPerPage, page, cfg)
	if err != nil {
		return nil, err
	}

	// loop through for pagination
	for len(repos)%resultsPerPage == 0 {
		page++
		newRepos, err := requestGitlabRepos(resultsPerPage, page, cfg)
		if err != nil {
			return nil, err
		}
		// prevent infinite loop if no results
		if len(newRepos) == 0 {
			break
		}
		repos = append(repos, newRepos...)
	}

	return repos, nil
}

func requestGitlabRepos(resultsPerPage int, pageNumber int, cfg *config.Config) ([]GitlabRepo, error) {
	APIPoint := fmt.Sprintf(
		"https://gitlab.com/api/v4/users/%s/projects?per_page=%d&page=%d",
		cfg.Client["gitlab"].Username,
		resultsPerPage,
		pageNumber,
	)
	client := &http.Client{}

	req, err := http.NewRequest("GET", APIPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Client["gitlab"].Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []GitlabRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil

}

func createRepoGitLab(repo GitRepo, cfg *config.Config) (GitlabRepo, error) {
	var newRepo GitlabRepo
	APIPoint := "https://gitlab.com/api/v4/projects"
	client := &http.Client{}

	visibility := "public"
	if repo.IsPrivate() {
		visibility = "private"
	}

	payload := map[string]any{
		"name":       repo.GetName(),
		"visibility": visibility,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return newRepo, err
	}

	req, err := http.NewRequest("POST", APIPoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return newRepo, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Client["gitlab"].Token)

	res, err := client.Do(req)
	if err != nil {
		return newRepo, err
	}

	if res.StatusCode != http.StatusCreated {
		return newRepo, fmt.Errorf("Failed to create GitLab repository. Status code: %d", res.StatusCode)
	}

	json.NewDecoder(res.Body).Decode(&newRepo)

	fmt.Println(fmt.Sprintf("GitLab repository name %s created successfully.", repo.GetName()))
	return newRepo, nil

}

func (gl *GitlabRepo) GetName() string {
	return gl.Name
}

func (gl *GitlabRepo) GetFullName() string {
	return gl.PathWithNamespace
}

func (gl *GitlabRepo) GetURL() string {
	return gl.HttpURLToRepo
}

func (gl *GitlabRepo) IsPrivate() bool {
	return gl.Visibility == "private"
}

func (gl *GitlabRepo) IsFork() bool {
	return gl.ForkedFromProject != nil
}

func (gl *GitlabRepo) LastUpdated() (*time.Time, error) {
	tm, err := time.Parse("2006-01-02T15:04:05Z", gl.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (gl *GitlabRepo) GetClientName() string {
	return "gitlab"
}
