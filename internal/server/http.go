package server

import (
	"log"
	"net/http"

	"github.com/tomersa/llm-gateway/internal/config"
	"github.com/tomersa/llm-gateway/internal/handler"

	"github.com/go-chi/chi/v5"
)

func Run() error {
	// Load config
	if err := config.LoadConfig("keys.json"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Setting up router
	r := chi.NewRouter()

	r.Post("/chat/completions", handler.HandleChat)
	r.Get("/health", handler.HandleHealth)
	r.Get("/metrics", handler.HandleMetrics)

	log.Println("Server listening on :8080")
	return http.ListenAndServe(":8080", r)
}
