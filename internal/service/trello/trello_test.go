package trello

import (
	"api-workflow/internal/client/github"
	"api-workflow/internal/client/trello"
	"api-workflow/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type argm map[string]string

var arguments []argm

type mockTrelloClient struct {
	BaseUrl string
	BoardId string
}

func (m *mockTrelloClient) CreateCard(listId string, pr github.PullRequest) {
	arguments = append(arguments, argm{"title": pr.Title, "url": pr.Url})
}

func (m *mockTrelloClient) GetLabelId(searches ...string) []string {
	return []string{}
}

var trelloService = TrelloService{
	StoragePath: "tmp/created_prs.txt",
	TrelloClient: &mockTrelloClient{
		BaseUrl: "",
		BoardId: "",
	},
}

func TestSyncCards(t *testing.T) {
	t.Run("should create cards", func(t *testing.T) {
		utils.ResetFile(utils.GetFromRoot(trelloService.StoragePath))
		arguments = nil

		ser1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			data := utils.GetFileData(utils.GetFromRoot("fixtures/pull-requests.json"))
			data2 := utils.GetFileData(utils.GetFromRoot("fixtures/pull-requests2.json"))
			switch r.URL.Path {
			case "/repos/example/app-stamped/pulls":
				w.Write([]byte(data))
			case "/repos/example/app-envelope-finalizer-service/pulls":
				w.Write([]byte(data2))
			default:
				t.Errorf("unexpected url: %s", r.URL.Path)
			}
		}))

		defer ser1.Close()

		trelloService.GithubBaseUrl = ser1.URL

		trelloService.SyncCards("123", []string{"spinelli-example"}, "example/app-stamped", "example/app-envelope-finalizer-service")

		expected := []argm{
			{
				"title": "[VOP-752] (T2) Configurar fontes do repositório ao subir o container",
				"url":   "https://api.github.com/repos/example/app-stamped/pulls/48",
			},
			{
				"title": "Add unified Terraform setup for multi-environment support",
				"url":   "https://api.github.com/repos/example/app-envelope-finalizer-service/pulls/43",
			},
			{
				"title": "[VOP-1129] Add table prefixing support in configuration",
				"url":   "https://api.github.com/repos/example/app-envelope-finalizer-service/pulls/42",
			},
			{
				"title": "VOP-1129 Provisionar ambiente de staging com DynamoDB, SQS e Lambda",
				"url":   "https://api.github.com/repos/example/app-envelope-finalizer-service/pulls/38",
			},
		}

		if !reflect.DeepEqual(expected, arguments) {
			t.Errorf("not equal expected(%d): %s arguments(%d) %s", len(expected), expected, len(arguments), arguments)
		}
	})
	t.Run("should not create card twice when already has on storage", func(t *testing.T) {
		utils.ResetFile(utils.GetFromRoot(trelloService.StoragePath))
		// one pr of fixtures/pull-requests2.json
		utils.SaveFile(utils.GetFromRoot(trelloService.StoragePath), "https://api.github.com/repos/example/app-envelope-finalizer-service/pulls/43\n")

		ser1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			data := utils.GetFileData(utils.GetFromRoot("fixtures/pull-requests2.json"))
			w.Write([]byte(data))
		}))

		trelloService.GithubBaseUrl = ser1.URL

		defer ser1.Close()

		arguments = nil
		trelloService.SyncCards("123", []string{"spinelli-example"}, "app-envelope-finalizer-service")
		if len(arguments) != 2 {
			t.Errorf("expected 1 card, got %d", len(arguments))
		}
	})
	t.Run("should create real cards", func(t *testing.T) {
		// real hit api test
		t.Skip()
		var tc = trello.TrelloClient{
			BoardId: "623ca78014f837744b044573",
			BaseUrl: "https://api.trello.com/1",
		}
		trelloService.TrelloClient = &tc
		trelloService.SyncCards("123", []string{"spinelli-example"}, "app-stamped", "app-envelope-finalizer-service")

		t.Error("not implemented")
	})
}
