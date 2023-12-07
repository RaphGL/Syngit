package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/raphgl/syngit/config"
)

func RequestReposAPI(clientName string, resultsPerPage int, pageNumber int, cfg *config.Config) (io.ReadCloser, error) {
	APIPoint := getReposAPIURL(clientName, resultsPerPage, pageNumber, cfg)
	client := &http.Client{}

	req, err := http.NewRequest("GET", APIPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Client[clientName].Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err

	}

	return res.Body, nil
}

func CreateRepoAPI(clientName string, repo GitRepo, cfg *config.Config) (io.ReadCloser, error) {
	APIPoint := getCreateRepoAPIURL(clientName)
	client := &http.Client{}

	payload := getPayload(clientName, repo)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", APIPoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Client[clientName].Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create GitHub repository. Status code: %d", res.StatusCode)
	}

	return res.Body, nil
}

func getReposAPIURL(clientName string, resultsPerPage int, pageNumber int, cfg *config.Config) string {
	var URL string
	switch clientName {
	case "github":
		URL = fmt.Sprintf(
			"https://api.github.com/user/repos?type=owner&per_page=%d&page=%d",
			resultsPerPage,
			pageNumber,
		)
	case "gitlab":
		URL = fmt.Sprintf(
			"https://gitlab.com/api/v4/users/%s/projects?per_page=%d&page=%d",
			cfg.Client["gitlab"].Username,
			resultsPerPage,
			pageNumber,
		)
	case "codeberg":
		URL = fmt.Sprintf("https://codeberg.org/api/v1/users/%s/repos?limit=%d&page=%d",
			cfg.Client["codeberg"].Username,
			resultsPerPage,
			pageNumber,
		)
	}

	return URL
}

func getCreateRepoAPIURL(clientName string) string {
	var URL string
	switch clientName {
	case "github":
		URL = "https://api.github.com/user/repos"
	case "gitlab":
		URL = "https://gitlab.com/api/v4/projects"
	case "codeberg":
		URL = "https://codeberg.org/api/v1/user/repos"
	}

	return URL
}

func getPayload(clientName string, repo GitRepo) map[string]any {
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

	return payload
}
