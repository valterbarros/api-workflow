package github

import (
	"api-workflow/utils"
	"fmt"
	"net/http"
	"os"
)

type PullRequest struct {
	// Keys should be in Capitalized style
	Url        string           `json:"url"`
	Title      string           `json:"title"`
	Labels     []map[string]any `json:"labels"`
	Updated_at string           `json:"updated_at"`
	Draft      bool             `json:"draft"`
	Number     int              `json:"number"`
	Html_url   string           `json:"html_url"`
	Links      struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Additions int `json:"additions"`
	Base      struct {
		Repo struct {
			Name string `json:"name"`
		} `json:"repo"`
	} `json:"base"`
	User struct {
		Login      string `json:"login"`
		Avatar_url string `json:"avatar_url"`
	} `json:"user"`
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
	if err != nil {
		fmt.Println("error doing the request", err)
	}
	defer resp.Body.Close()

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
	if err != nil {
		fmt.Println("error doing the request", err)
	}
	defer resp.Body.Close()

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
	if err != nil {
		fmt.Println("error doing the request", err)
	}
	defer resp.Body.Close()

	var pullRequests PullRequests
	utils.DecodeJson(resp.Body, &pullRequests)

	return pullRequests
}
