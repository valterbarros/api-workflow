package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getPullRequests(w http.ResponseWriter, r *http.Request) {
	items := []Item{{ID: "1", Name: "Item One"}, {ID: "2", Name: "Item Two"}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
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
