package trello

import (
	"api-workflow/internal/client/github"
	"api-workflow/internal/client/trello"
	"api-workflow/utils"
	"slices"
	"strings"
)

type TrelloService struct {
	TrelloClient trello.TrelloClientInterface
	StoragePath  string `default:"internal/storage/created_prs.txt"`
	// TODO: move to github client
	GithubBaseUrl string `default:"https://api.github.com"`
}

func (t *TrelloService) SyncCards(listId string, users []string, reposNames ...string) {
	// List all prs from all repos selecteds
	// For each create a new trello card with title = gh title
	// TODO: comes from request
	// repos := []string{"tavola", "app-workspace", "nova-widget"}
	// pulls := []github.PullRequests{}
	// skipPrNames := "val, test"
	// usersFilter := []string{"joclicker"}
	usersFilter := users

	createdPrs := utils.GetFileData(utils.GetFromRoot(t.StoragePath))

	merged := github.PullRequests{}

	for _, repo := range reposNames {
		pulls := github.GetPullRequests(t.urlFromRepo(repo))
		// filters
		pulls = slices.DeleteFunc(pulls, func(p github.PullRequest) bool {
			return !slices.Contains(usersFilter, p.User.Login) || strings.Contains(createdPrs, p.Url) || p.Draft
		})
		merged = append(merged, pulls...)
	}

	var newPrs = strings.Builder{}
	newPrs.WriteString(createdPrs)

	for _, pull := range merged {
		// TODO: Save created cards to a file
		t.TrelloClient.CreateCard(listId, pull)
		newPrs.WriteString(pull.Url + "\n")
	}

	utils.SaveFile(utils.GetFromRoot(t.StoragePath), newPrs.String())
}

func (t *TrelloService) urlFromRepo(repo string) string {
	return t.GithubBaseUrl + "/repos/" + repo + "/pulls"
}
