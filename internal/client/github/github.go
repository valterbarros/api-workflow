package github

import (
	"api-workflow/utils"
	"fmt"
	"net/http"
	"os"
)

type PullRequest struct {
	// Keys should be in Capitalized style
	Url    string
	Title  string
	Labels []map[string]any
	User   struct {
		Login      string
		Avatar_url string
	}
}

func GetPullRequest(url string) PullRequest {
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

	if err != nil {
		fmt.Println("error doing the request", err)
	}

	defer resp.Body.Close()

	var pr PullRequest
	utils.DecodeJson(resp.Body, &pr)

	return pr
}

type User struct {
	User struct {
		Login      string
		Avatar_url string
	}
	State string
}

type Reviews []User

// returns raw reviews from on pr
func GetReviews(url string) Reviews {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := "https://api.github.com/repos/clicksign/tavola/pulls/14781/reviews"

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

	var review Reviews
	utils.DecodeJson(resp.Body, &review)

	return review
}

type repos []struct {
	Id   int
	Name string
}

func GetRepos(url string) repos {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := "https://api.github.com/user/repos?sort=pushed&per_page=30"

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

	var repos repos
	utils.DecodeJson(resp.Body, &repos)

	return repos
}

type PullRequests []PullRequest

func GetPullRequests(url string) PullRequests {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := "https://api.github.com/repos/clicksign/tavola/pulls"

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

	var pullRequests PullRequests
	utils.DecodeJson(resp.Body, &pullRequests)

	return pullRequests
}
