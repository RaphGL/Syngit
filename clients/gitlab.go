package clients

import (
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
	APIPoint := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/projects", cfg.Client["gitlab"].Username)
	res, err := http.Get(APIPoint)
	if err != nil {
		return nil, err
	}

	var repos []GitlabRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil
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
