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

const BASE_URL = "https://api.trello.com/1"
const DOING_LIST_ID = "623ca78014f837744b044577"
const TAVOLA_LABEL_ID = "649cdc411b86084babfb4605"

func CreateCard(urlPath string, pr github.PullRequest) {
	utils.SetEnvFromFile()
	client := &http.Client{}
	currentUrl := BASE_URL + "/cards"

	if urlPath != "" {
		currentUrl = urlPath
	}

	token := os.Getenv("TRELLO_TOKEN")
	key := os.Getenv("TRELLO_KEY")

	// const pullRequestApiUrl = pullRequest.url;
	// const userName = pullRequest.user.login;
	// const updatedPr = pullRequest.updated_at;
	// const splittedUrl = pullRequestUrl.split('/');
	// const prNumber = splittedUrl[splittedUrl.length - 1];
	// const cardTitle = pullRequest.title;
	// const repoName = pullRequest.base.repo.name;
	// const hasReadtyToReview = pullRequest.labels.find(({ name }) => name.includes('review'));
	// const isDraft = pullRequest.draft;

	// const { additions } = await gitHubRequest(pullRequestApiUrl, githubUserInfo.ghToken);

	// const { id: boardId } = await t.board('id');

	// const userLabelId = await getLabelId(boardId, userName);
	// const repoLabelId = await getLabelId(boardId, repoName);
	// const readyToReviewLabelId = await getLabelId(boardId, 'ready to review');
	// const idLabels = [userLabelId, repoLabelId].filter(Boolean);

	// if (hasReadtyToReview && readyToReviewLabelId && isDraft) idLabels.push(readyToReviewLabelId);

	// const createCardParams = {
	// 	name: `${cardTitle} [${repoName}] [${userName}] #${prNumber} [${formatDate(updatedPr)}] [a: ${additions}]`,
	// 	idList: listBoardId,
	// 	...( idLabels.length && { idLabels: idLabels.join(',') }),
	// 	pos: "top",
	// };

	data := url.Values{}
	data.Set("key", key)
	data.Set("token", token)
	// auth end

	pullRequest := github.GetPullRequest(pr.Links.Self.Href)

	name := fmt.Sprintf("%s [%s] #%d[%s] [%s] [a: %d]", pullRequest.Title, pullRequest.Base.Repo.Name, pullRequest.Number, pullRequest.User.Login, pullRequest.Updated_at, pullRequest.Additions)
	data.Set("name", name)
	data.Set("idList", DOING_LIST_ID)
	data.Set("pos", "top")

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
