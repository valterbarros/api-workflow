package requests

import (
	"api-workflow/utils"
	"fmt"
	"net/http"
	"os"
)

type pullRequest struct {
	// Keys should be in Capitalized style
	Url    string
	Labels []map[string]any
}

func getPullRequest(url string) pullRequest {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := "https://api.github.com/repos/clicksign/app-workspaces/pulls/162"

	if url != "" {
		currentUrl = url
	}

	req, err := http.NewRequest("GET", currentUrl, nil)
	if err != nil {
		fmt.Println("error creating the request", err)
	}
	token := "Bearer " + os.Getenv("GH_TOKEN")

	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("error doing the request", err)
	}

	var pr pullRequest
	utils.DecodeJson(resp.Body, &pr)

	return pr
}
