package trello

import (
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

const BASE_URL = "https://api.trello.com/1"
const DOING_LIST_ID = "623ca78014f837744b044577"
const TAVOLA_LABEL_ID = "649cdc411b86084babfb4605"

func CreateCard(urlPath string, name string) {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := BASE_URL + "/cards"

	if urlPath != "" {
		currentUrl = urlPath
	}

	token := os.Getenv("TRELLO_TOKEN")
	key := os.Getenv("TRELLO_KEY")

	data := url.Values{}
	data.Set("key", key)
	data.Set("token", token)
	data.Set("name", name)
	data.Set("idList", DOING_LIST_ID)

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
