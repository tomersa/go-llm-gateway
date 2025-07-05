package server

import (
	"log"
	"net/http"

	"github.com/tomersa/llm-gateway/internal/handler"

	"github.com/go-chi/chi/v5"
)

func Run() error {
	r := chi.NewRouter()

	r.Post("/chat/completions", handler.HandleChat)

	log.Println("Server listening on :8080")
	return http.ListenAndServe(":8080", r)
}
