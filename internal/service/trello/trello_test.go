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

func (m *mockTrelloClient) CreateCard(pr github.PullRequest) {
	arguments = append(arguments, argm{"title": pr.Title, "url": pr.Url})
}

func (m *mockTrelloClient) GetLabelId(searches ...string) []string {
	return []string{}
}

var tcm = mockTrelloClient{
	BaseUrl: "",
	BoardId: "",
}

func TestCreatePrCard(t *testing.T) {
	t.Run("should create cards", func(t *testing.T) {
		arguments = nil

		ser1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			data := utils.GetFileData(utils.GetFromRoot("fixtures/pull-requests.json"))
			w.Write([]byte(data))
		}))

		defer ser1.Close()

		ser2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			data := utils.GetFileData(utils.GetFromRoot("fixtures/pull-requests2.json"))
			w.Write([]byte(data))
		}))

		defer ser2.Close()

		CreatePrCard(&tcm, []string{"spinelli-clk"}, ser1.URL, ser2.URL)

		expected := []argm{
			{
				"title": "[VOP-752] (T2) Configurar fontes do repositório ao subir o container",
				"url":   "https://api.github.com/repos/clicksign/app-stamped/pulls/48",
			},
			{
				"title": "Add unified Terraform setup for multi-environment support",
				"url":   "https://api.github.com/repos/clicksign/app-envelope-finalizer-service/pulls/43",
			},
			{
				"title": "[VOP-1129] Add table prefixing support in configuration",
				"url":   "https://api.github.com/repos/clicksign/app-envelope-finalizer-service/pulls/42",
			},
			{
				"title": "VOP-1129 Provisionar ambiente de staging com DynamoDB, SQS e Lambda",
				"url":   "https://api.github.com/repos/clicksign/app-envelope-finalizer-service/pulls/38",
			},
		}

		if !reflect.DeepEqual(expected, arguments) {
			t.Errorf("not equal expected: %s arguments %s", expected, arguments)
		}
	})
	t.Run("should create real cards", func(t *testing.T) {
		// real test
		var tc = trello.TrelloClient{
			BoardId: "623ca78014f837744b044573",
			BaseUrl: "https://api.trello.com/1",
		}
		CreatePrCard(&tc, []string{"spinelli-clk"}, "https://api.github.com/repos/clicksign/app-envelope-finalizer-service/pulls")

		t.Error("not implemented")
	})
	// t.Run("Check for not pull same pr twice")
}
