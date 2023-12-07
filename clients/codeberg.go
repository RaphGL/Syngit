package clients

import (
	"encoding/json"
	"fmt"
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
	resultsPerPage := 100 //max result for gitlab
	page := 1
	resBody, err := RequestReposAPI("codeberg", resultsPerPage, page, cfg)
	if err != nil {
		return nil, err
	}
	var repos []CodebergRepo
	json.NewDecoder(resBody).Decode(&repos)

	// loop through for pagination
	for len(repos)%resultsPerPage == 0 {
		page++

		newResBody, err := RequestReposAPI("codeberg", resultsPerPage, page, cfg)
		if err != nil {
			return nil, err
		}
		var newRepos []CodebergRepo
		json.NewDecoder(newResBody).Decode(&newRepos)
		// prevent infinite loop if no results
		if len(newRepos) == 0 {
			break
		}
		repos = append(repos, newRepos...)
	}

	return repos, nil
}

func createRepoCodeberg(repo GitRepo, cfg *config.Config) (CodebergRepo, error) {
	var newRepo CodebergRepo
	resBody, err := CreateRepoAPI("codeberg", repo, cfg)
	if err != nil {
		return newRepo, err
	}

	json.NewDecoder(resBody).Decode(&newRepo)

	fmt.Println(fmt.Sprintf("Codeberg repository name %s created successfully.", repo.GetName()))
	return newRepo, nil
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
