package trello

import (
	"api-workflow/internal/client/github"
	"api-workflow/utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type _pullRequest struct {
	// Keys should be in Capitalized style
	Url    string
	Labels []map[string]any
}

type TrelloClientInterface interface {
	CreateCard(pr github.PullRequest)
	GetLabelId(search string) string
}

type TrelloClient struct {
	BoardId string `default:"623ca78014f837744b044573"`
	BaseUrl string `default:"https://api.trello.com/1"`
}

const BASE_URL = "https://api.trello.com/1"
const DOING_LIST_ID = "623ca78014f837744b044577"
const BOARD_ID = "623ca78014f837744b044573"
const TAVOLA_LABEL_ID = "649cdc411b86084babfb4605"

func (t *TrelloClient) CreateCard(pr github.PullRequest) {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := t.BaseUrl + "/cards"

	token := os.Getenv("TRELLO_TOKEN")
	key := os.Getenv("TRELLO_KEY")

	data := url.Values{}
	data.Set("key", key)
	data.Set("token", token)
	// auth end

	pullRequest := github.GetPullRequest(pr.Links.Self.Href)

	name := fmt.Sprintf("%s [%s] #%d[%s] [%s] [a: %d]", pullRequest.Title, pullRequest.Base.Repo.Name, pullRequest.Number, pullRequest.User.Login, pullRequest.Updated_at, pullRequest.Additions)
	data.Set("name", name)
	data.Set("idList", DOING_LIST_ID)
	data.Set("pos", "top")
	labels := []string{t.GetLabelId(pullRequest.User.Login), t.GetLabelId(pullRequest.Base.Repo.Name)}
	data.Set("idLabels", strings.Join(labels, ","))

	req, err := http.NewRequest("POST", currentUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("error creating the request", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error doing the request", err)
	}

	defer resp.Body.Close()
}

type Label struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// simple cache label
var labels []Label

func (t *TrelloClient) GetLabelId(search string) string {
	if len(labels) > 0 {
		return findLabel(search).Id
	}

	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := t.BaseUrl + "/boards/" + t.BoardId + "/labels"

	data := url.Values{}
	data.Set("key", os.Getenv("TRELLO_KEY"))
	data.Set("token", os.Getenv("TRELLO_TOKEN"))

	req, err := http.NewRequest("GET", currentUrl+"?"+data.Encode(), nil)
	if err != nil {
		fmt.Println("error creating the request", err)
		return ""
	}

	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error doing the request", err)
		return ""
	}

	defer resp.Body.Close()

	utils.DecodeJson(resp.Body, &labels)

	return findLabel(search).Id
}

func findLabel(search string) Label {
	for _, label := range labels {
		if strings.Contains(strings.ToLower(label.Name), strings.ToLower(search)) {
			return label
		}
	}
	return Label{}
}

// TODO: Create attachment to card
// func CreateAttachment(urlPath string, pr github.PullRequest) {}
