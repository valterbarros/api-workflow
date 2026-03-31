package trello

import (
	"api-workflow/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCard(t *testing.T) {
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

	labelIds := tc.GetLabelId("tavola", "nova-widget")

	ser1.Close()

	if len(labelIds) < 2 {
		t.Error("labelId is incomplete", labelIds)
	}

	labelIds2 := tc.GetLabelId("tavola")

	if len(labelIds2) < 1 {
		t.Error("not cached labelId", labelIds2)
	}
}
