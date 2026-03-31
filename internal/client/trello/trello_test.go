package trello

import (
	"api-workflow/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePrCards(t *testing.T) {
	// TODO
}

func TestGetLabelId(t *testing.T) {
	ser1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := utils.GetFileData(utils.GetFromRoot("fixtures/trello-labels.json"))
		w.Write([]byte(data))
	}))

	tc := TrelloClient{
		BoardId: "623ca78014f837744b044573",
		BaseUrl: ser1.URL,
	}

	labelId1 := tc.GetLabelId("3 - tavola")

	ser1.Close()

	if labelId1 == "" {
		t.Error("labelId is empty", labelId1)
	}

	labelId2 := tc.GetLabelId("3 - tavola")

	if labelId2 == "" {
		t.Error("not cached labelId", labelId2)
	}
}
