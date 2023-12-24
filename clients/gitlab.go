package clients

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/raphgl/syngit/config"
)

type GitlabRepo struct {
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility        string `json:"visibility"`
	HttpURLToRepo     string `json:"http_url_to_repo"`
	LastActivityAt    string `json:"last_activity_at"`
	ForkedFromProject *any   `json:"forked_from_project"`
}

func getGitlabRepos(cfg *config.Config) ([]GitlabRepo, error) {
	resultsPerPage := 100 //max result for gitlab
	page := 1
	resBody, err := getUserReposAPI("gitlab", resultsPerPage, page, cfg)
	if err != nil {
		return nil, err
	}
	var repos []GitlabRepo
	json.NewDecoder(resBody).Decode(&repos)

	// loop through for pagination
	for len(repos)%resultsPerPage == 0 {
		page++

		newResBody, err := getUserReposAPI("gitlab", resultsPerPage, page, cfg)
		if err != nil {
			return nil, err
		}
		var newRepos []GitlabRepo
		json.NewDecoder(newResBody).Decode(&newRepos)
		// prevent infinite loop if no results
		if len(newRepos) == 0 {
			break
		}
		repos = append(repos, newRepos...)
	}

	return repos, nil
}

func createRepoGitlab(repo GitRepo, cfg *config.Config) (GitlabRepo, error) {
	var newRepo GitlabRepo
	resBody, err := createRepoAPI("gitlab", repo, cfg)
	if err != nil {
		return newRepo, err
	}

	json.NewDecoder(resBody).Decode(&newRepo)

	fmt.Println(fmt.Sprintf("Gitlab repository name %s created successfully.", repo.GetName()))
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
	tm, err := time.Parse("2006-01-02T15:04:05Z", gl.LastActivityAt)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (gl *GitlabRepo) GetClientName() string {
	return "gitlab"
}
