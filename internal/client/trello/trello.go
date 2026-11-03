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

type Card struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TrelloClientInterface interface {
	CreateCard(listId string, pr github.PullRequest)
	GetLabelId(searches ...string) []string
}

type TrelloClient struct {
	BoardId string `default:"623ca78014f837744b044573"`
	BaseUrl string `default:"https://api.trello.com/1"`
}

const BASE_URL = "https://api.trello.com/1"

// const DOING_LIST_ID = "623ca78014f837744b044577"
const BOARD_ID = "623ca78014f837744b044573"

func (t *TrelloClient) CreateCard(listId string, pr github.PullRequest) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error creating the card", r)
		}
	}()

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
	data.Set("idList", listId)
	data.Set("pos", "top")
	labels := t.GetLabelId(pullRequest.User.Login, pullRequest.Base.Repo.Name)
	data.Set("idLabels", strings.Join(labels, ","))
	// add pr as attachment

	req, err := http.NewRequest("POST", currentUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("error creating the request", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	var card Card
	utils.DecodeJson(resp.Body, &card)
	go t.CreateAttachment(pr.Html_url, card.Id)

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

func (t *TrelloClient) GetLabelId(searches ...string) []string {
	if len(labels) > 0 {
		return findLabels(labels, searches)
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
		return []string{}
	}

	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error doing the request", err)
		return []string{}
	}

	defer resp.Body.Close()

	utils.DecodeJson(resp.Body, &labels)

	return findLabels(labels, searches)
}

func findLabels(labelsToSearch []Label, searches []string) []string {
	foundLabels := []string{}
	for _, label := range labelsToSearch {
		for _, search := range searches {
			if strings.Contains(strings.ToLower(label.Name), strings.ToLower(search)) {
				foundLabels = append(foundLabels, label.Id)
			}
		}
	}
	return foundLabels
}

func (t *TrelloClient) CreateAttachment(urlPr string, cardId string) {
	fmt.Println("creating attachment for", urlPr, "and", cardId)
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := t.BaseUrl + "/card/" + cardId + "/attachments"

	data := url.Values{
		"key":   {os.Getenv("TRELLO_KEY")},
		"token": {os.Getenv("TRELLO_TOKEN")},
		"name":  {"Pull Request"},
		"url":   {urlPr},
	}

	req, err := http.NewRequest("POST", currentUrl, strings.NewReader(data.Encode()))

	if err != nil {
		fmt.Println("error creating the request", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Go-HTTP-Client")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error doing the request", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		fmt.Println("error creating the attachment", resp.StatusCode)
	}
}
