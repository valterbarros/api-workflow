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

func getRepos() {
	utils.SetEnvFromFile()
	client := &http.Client{}
	url := "https://api.github.com/repos/clicksign/app-workspaces/pulls/162"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error 1", err)
	}
	token := "Bearer " + os.Getenv("GH_TOKEN")

	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("error 2", err)
	}

	var pr pullRequest
	utils.DecodeJson(resp.Body, &pr)

	fmt.Println(pr.Labels)
}
