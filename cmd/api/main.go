package main

import (
	"api-workflow/internal/service/github"
	"api-workflow/utils"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getPullRequests(w http.ResponseWriter, r *http.Request) {
	items := github.GetPullRequests()
	w.Header().Set("Content-Type", "application/json")
	utils.EncodeJson(w, items)
}

func main() {
	const PORT = "12769"

	mux := http.NewServeMux()

	// Sua rota de exemplo
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Conteúdo do Iframe via HTTPS e Go"))
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// O navegador exige isso para 'loopback' também!
		w.Header().Set("Access-Control-Allow-Private-Network", "true")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		mux.ServeHTTP(w, r)
	})

	log.Println("Servidor rodando em https://localhost:", PORT)

	// Usando o seu ListenAndServeTLS
	// Certifique-se de que os caminhos para .crt e .key estão corretos
	err := http.ListenAndServeTLS(":"+PORT, utils.GetFromRoot("certs/localhost+3.pem"), utils.GetFromRoot("certs/localhost+3-key.pem"), handler)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}
