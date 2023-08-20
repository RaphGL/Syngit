package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raphgl/syngit/config"
)

type CodebergRepo struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	CloneURL  string `json:"clone_url"`
	UpdatedAt string `json:"updated_at"`
	Private   bool   `json:"private"`
	Fork      bool   `json:"fork"`
}

func getCodebergRepos(cfg *config.Config) ([]CodebergRepo, error) {
	APIPoint := fmt.Sprintf("https://codeberg.org/api/v1/users/%s/repos", cfg.Client["codeberg"].Username)
	res, err := http.Get(APIPoint)
	if err != nil {
		return nil, err
	}

	var repos []CodebergRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil
}

func (cb *CodebergRepo) GetName() string {
	return cb.Name
}

func (cb *CodebergRepo) GetFullName() string {
	return cb.FullName
}

func (cb *CodebergRepo) GetURL() string {
	return cb.CloneURL
}

func (cb *CodebergRepo) IsPrivate() bool {
	return cb.Private
}

func (cb *CodebergRepo) IsFork() bool {
	return cb.Fork
}

func (cb *CodebergRepo) LastUpdated() (*time.Time, error) {
	tm, err := time.Parse("2006-01-02T15:04:05Z", cb.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (cb *CodebergRepo) GetClientName() string {
	return "codeberg"
}
