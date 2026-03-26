package requests

import (
	"api-workflow/utils"
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

	result := getPullRequest(ser.URL)

	if len(result.Labels) < 1 {
		t.Errorf("No labels found")
	}
}

func TestGetReviews(t *testing.T) {
	// Mock server
	ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := utils.GetFileData("../fixtures/reviews.json")
		w.Write([]byte(data))
	}))

	defer ser.Close()

	result := getReviews(ser.URL)

	if len(result) < 1 {
		t.Errorf("No results found")
	}

	if result[0].User["avatar_url"] == nil {
		t.Errorf("No avata found in 0 position")
	}
}

func TestGetRepos(t *testing.T) {
	// Mock server
	ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := utils.GetFileData("../fixtures/repos.json")
		w.Write([]byte(data))
	}))

	defer ser.Close()

	result := getRepos(ser.URL)

	if len(result) < 1 {
		t.Errorf("No results found")
	}

	if result[0].Name == "" {
		t.Errorf("No name found in 0 position")
	}
}

func TestGetPullRequests(t *testing.T) {
	// Mock server
	ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := utils.GetFileData("../fixtures/pull-requests.json")
		w.Write([]byte(data))
	}))

	defer ser.Close()

	result := getPullRequests(ser.URL)

	if len(result) < 1 {
		t.Errorf("No results found")
	}

	if result[0].Url == "" {
		t.Errorf("No url found in 0 position")
	}
}
