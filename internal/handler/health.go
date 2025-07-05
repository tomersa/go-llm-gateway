package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tomersa/llm-gateway/internal/provider"
)

type ProviderStatus struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Online bool   `json:"online"`
}

type HealthResponse struct {
	Status    string           `json:"status"`
	Providers []ProviderStatus `json:"providers"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	statuses := make([]ProviderStatus, 0, len(provider.AiServiceEndpoints))
	for name, url := range provider.AiServiceEndpoints {
		client := http.Client{Timeout: 2 * time.Second}
		resp, err := client.Head(url)
		online := err == nil && resp.StatusCode < 500
		statuses = append(statuses, ProviderStatus{
			Name:   name,
			URL:    url,
			Online: online,
		})
		if resp != nil {
			resp.Body.Close()
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "ok",
		Providers: statuses,
	})
}
