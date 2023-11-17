package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
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
	client := &http.Client{}

	req, err := http.NewRequest("GET", APIPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Client["codeberg"].Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []CodebergRepo
	json.NewDecoder(res.Body).Decode(&repos)

	return repos, nil
}

func createRepoCodeberg(repo GitRepo, cfg *config.Config) error {
	APIPoint := "https://codeberg.org/api/v1/users/repos"
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
	req.Header.Set("Authorization", "Bearer "+cfg.Client["codeberg"].Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		slog.Error("Failed to create Codeberg repository. Status code: %d", res.StatusCode)
		return err
	}

	fmt.Println(fmt.Sprintf("Codeberg repository name %s created successfully.", repo.GetName()))
	return nil

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
