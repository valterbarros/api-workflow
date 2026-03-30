package main

import (
	"api-workflow/internal/service/github"
	"api-workflow/utils"
	"log"
	"net/http"
)

func buildMid(mux *http.ServeMux) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		mux.ServeHTTP(w, r)
	}
}

func getPullRequests(w http.ResponseWriter, r *http.Request) {
	items := github.GetPullRequests()
	w.Header().Set("Content-Type", "application/json")
	utils.EncodeJson(w, items)
}

func main() {
	const PORT = "12769"

	mux := http.NewServeMux()

	mux.HandleFunc("/reviews", func(w http.ResponseWriter, r *http.Request) {
		// TODO: pass only the repo and name and pr id
		url := r.URL.Query().Get("url")
		result := github.ReviewsResult(url)
		w.Header().Set("Content-Type", "application/json")
		utils.EncodeJson(w, result)
	})

	mux.HandleFunc("/pull_requests", getPullRequests)
	// Go 1.22+ ServeMux: só casa o path exato; barra no final dava 404.
	mux.HandleFunc("/pull_requests/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pull_requests/" {
			http.NotFound(w, r)
			return
		}
		getPullRequests(w, r)
	})

	handler := http.HandlerFunc(buildMid(mux))

	log.Println("Servidor rodando em https://localhost:", PORT)

	// Usando o seu ListenAndServeTLS
	// Certifique-se de que os caminhos para .crt e .key estão corretos
	err := http.ListenAndServeTLS(":"+PORT, utils.GetFromRoot("certs/localhost+3.pem"), utils.GetFromRoot("certs/localhost+3-key.pem"), handler)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}
