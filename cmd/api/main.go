package main

import (
	"api-workflow/internal/service/github"
	"api-workflow/utils"
	"fmt"
	"log"
	"net/http"
)

func getPullRequests(w http.ResponseWriter, r *http.Request) {
	items := github.GetPullRequests()
	w.Header().Set("Content-Type", "application/json")
	utils.EncodeJson(w, items)
}

func main() {
	http.HandleFunc("/pull_requests", getPullRequests)
	const PORT = "12769"

	go func() {
		fmt.Printf("is running on port %v", PORT)
	}()

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
}
