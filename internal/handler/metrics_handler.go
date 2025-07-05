package handler

import (
	"encoding/json"
	"net/http"
)

type MetricsResponse struct {
	TotalRequests         int                `json:"total_requests"`
	RequestsPerProvider   map[string]int     `json:"requests_per_provider"`
	AverageResponseTimeMs float64            `json:"average_response_time_ms"`
	AveragePerProviderMs  map[string]float64 `json:"average_response_time_per_provider_ms"`
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	total, perProvider, avg, avgPerProvider := metrics.Snapshot()
	resp := MetricsResponse{
		TotalRequests:         total,
		RequestsPerProvider:   perProvider,
		AverageResponseTimeMs: avg,
		AveragePerProviderMs:  avgPerProvider,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
