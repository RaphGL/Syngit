package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphgl/syngit/config"
)

type CodebergRepo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`
}

func GetCodebergRepos(cfg *config.Config) ([]CodebergRepo, error) {
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
