package main

import (
	"api-workflow/internal/client/trello"
	"api-workflow/internal/service/github"
	trelloService "api-workflow/internal/service/trello"
	"api-workflow/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

// url: /sync_cards?listId=623ca78014f837744b044577&users=joaozinho&repos=https://api.github.com/repos/example/app-envelope-finalizer-service/pulls
func syncCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	boardId := r.URL.Query().Get("board_id")

	trelloService := trelloService.TrelloService{
		StoragePath:   "internal/storage/created_prs.txt",
		GithubBaseUrl: "https://api.github.com",
		TrelloClient: &trello.TrelloClient{
			BoardId: boardId,
			BaseUrl: "https://api.trello.com/1",
		},
	}

	listId := r.URL.Query().Get("list_id")
	users := r.URL.Query().Get("users")
	repos := r.URL.Query().Get("repos")

	if boardId == "" || listId == "" || users == "" || repos == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.EncodeJson(w, "Invalid request")
		return
	}

	trelloService.SyncCards(listId, strings.Split(users, ","), strings.Split(repos, ",")...)
	w.Header().Set("Content-Type", "application/json")
	utils.EncodeJson(w, "Cards synced")
}

func main() {
	const PORT = "12769"

	signals := make(chan os.Signal, 1)
	// closes server on siging -9
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	handler := http.HandlerFunc(buildMid(mux))

	mux.HandleFunc("/reviews", func(w http.ResponseWriter, r *http.Request) {
		// TODO: pass only the repo and name and pr id
		url := r.URL.Query().Get("url")
		result := github.ReviewsResult(url)
		w.Header().Set("Content-Type", "application/json")
		utils.EncodeJson(w, result)
	})

	mux.HandleFunc("/pull_requests", getPullRequests)
	mux.HandleFunc("/pull_requests/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pull_requests/" {
			http.NotFound(w, r)
			return
		}
		getPullRequests(w, r)
	})

	mux.HandleFunc("/sync_cards", syncCards)

	log.Println("Servidor rodando em https://localhost:", PORT)

	err := http.ListenAndServeTLS(":"+PORT, utils.GetFromRoot("certs/localhost+3.pem"), utils.GetFromRoot("certs/localhost+3-key.pem"), handler)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}
