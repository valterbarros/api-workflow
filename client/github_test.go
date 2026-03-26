package requests

import (
	"api-workflow/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRepo(t *testing.T) {
	// Mock server
	ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := utils.GetFileData("../fixtures/json-return.json")
		w.Write([]byte(data))
	}))

	defer ser.Close()

	fmt.Println(ser.URL)

	// TODO
	result := getPullRequest(ser.URL)

	if len(result.Labels) < 1 {
		t.Errorf("No labels found")
	}
}
