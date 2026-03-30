package trello

import (
	"api-workflow/internal/client/github"
	"fmt"
	"slices"
)

func CreatePrCard(trelloClient func(string, string), repos ...string) {
	// List all prs from all repos selecteds
	// For each create a new trello card with title = gh title
	// TODO: comes from request
	// repos := []string{"tavola", "app-workspace", "nova-widget"}
	// pulls := []github.PullRequests{}
	usersFilter := []string{"spinelli-clk"}

	merged := github.PullRequests{}

	for _, repo := range repos {
		// ghUrl := "https://api.github.com/repos/clicksign/" + repo + "/pulls"
		pulls := github.GetPullRequests(repo)
		// filters
		pulls = slices.DeleteFunc(pulls, func(p github.PullRequest) bool {
			return !slices.Contains(usersFilter, p.User.Login)
		})
		merged = append(merged, pulls...)
	}

	fmt.Println("size: ", len(merged))

	for _, pull := range merged {
		trelloClient(pull.Url, pull.Title)
	}
}
