package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/raphgl/syngit/config"
)

func getUserReposAPI(clientName string, resultsPerPage int, pageNumber int, cfg *config.Config) (io.ReadCloser, error) {
	var endpoint string
	switch clientName {
	case "github":
		endpoint = fmt.Sprintf(
			"https://api.github.com/user/repos?type=owner&per_page=%d&page=%d",
			resultsPerPage,
			pageNumber,
		)
	case "gitlab":
		endpoint = fmt.Sprintf(
			"https://gitlab.com/api/v4/users/%s/projects?per_page=%d&page=%d",
			cfg.Client["gitlab"].Username,
			resultsPerPage,
			pageNumber,
		)
	case "codeberg":
		endpoint = fmt.Sprintf("https://codeberg.org/api/v1/users/%s/repos?limit=%d&page=%d",
			cfg.Client["codeberg"].Username,
			resultsPerPage,
			pageNumber,
		)
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Client[clientName].Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err

	}

	return res.Body, nil
}

func createRepoAPI(clientName string, repo GitRepo, cfg *config.Config) (io.ReadCloser, error) {
    // -- get client specific endpoint
	var endpoint string
	switch clientName {
	case "github":
		endpoint = "https://api.github.com/user/repos"
	case "gitlab":
		endpoint = "https://gitlab.com/api/v4/projects"
	case "codeberg":
		endpoint = "https://codeberg.org/api/v1/user/repos"
	}

    // -- create client payload
	var payload map[string]any
	switch clientName {
	case "gitlab":
		visibility := "public"
		if repo.IsPrivate() {
			visibility = "private"
		}

		payload = map[string]any{
			"name":       repo.GetName(),
			"visibility": visibility,
		}
	default:
		payload = map[string]any{
			"name":    repo.GetName(),
			"private": repo.IsPrivate(),
		}
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

    // --- create repo in client
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Client[clientName].Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create %s repository. Status code: %d", clientName, res.StatusCode)
	}

	return res.Body, nil
}
