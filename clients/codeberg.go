package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CodebergRepo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`
}

func GetCodebergRepos() ([]CodebergRepo, error) {
	APIPoint := fmt.Sprintf("https://codeberg.org/api/v1/users/%s/repos", user)
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
