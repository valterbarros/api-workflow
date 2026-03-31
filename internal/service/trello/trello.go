package trello

import (
	"api-workflow/internal/client/github"
	"api-workflow/internal/client/trello"
	"api-workflow/utils"
	"slices"
	"strings"
)

type TrelloService struct {
	trelloClient trello.TrelloClientInterface
	storagePath  string `default:"internal/storage/created_prs.txt"`
}

func (t *TrelloService) CreatePrCard(users []string, repos ...string) {
	// List all prs from all repos selecteds
	// For each create a new trello card with title = gh title
	// TODO: comes from request
	// repos := []string{"tavola", "app-workspace", "nova-widget"}
	// pulls := []github.PullRequests{}
	// skipPrNames := "val, test"
	// usersFilter := []string{"joclicker"}
	usersFilter := users

	createdPrs := utils.GetFileData(utils.GetFromRoot(t.storagePath))

	merged := github.PullRequests{}

	for _, repo := range repos {
		// ghUrl := "https://api.github.com/repos/clicksign/" + repo + "/pulls"
		pulls := github.GetPullRequests(repo)
		// filters
		pulls = slices.DeleteFunc(pulls, func(p github.PullRequest) bool {
			return !slices.Contains(usersFilter, p.User.Login) || strings.Contains(createdPrs, p.Url)
		})
		merged = append(merged, pulls...)
	}

	var newPrs = strings.Builder{}
	newPrs.WriteString(createdPrs)

	for _, pull := range merged {
		// TODO: Save created cards to a file
		t.trelloClient.CreateCard(pull)
		newPrs.WriteString(pull.Url + "\n")
	}

	utils.SaveFile(utils.GetFromRoot(t.storagePath), newPrs.String())
}
