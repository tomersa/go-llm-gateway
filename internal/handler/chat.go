package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tomersa/llm-gateway/internal/config"
	"github.com/tomersa/llm-gateway/internal/provider"
)

func HandleChat(w http.ResponseWriter, request *http.Request) {
	start := time.Now()

	vk, err := extractVirtualKey(request.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Failed to extract virtual key", http.StatusBadRequest)
		return
	}

	info, ok := config.Config[vk]
	if !ok {
		http.Error(w, fmt.Sprintf("Unauthorized: invalid virtual key %s", vk), http.StatusUnauthorized)
		return
	}

	aiServiceEndpoint, ok := provider.AiServiceEndpoints[info.Provider]
	if !ok {
		http.Error(w, fmt.Sprintf("Unsupported provider: %s", info.Provider), http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}
	request.Body.Close() //  must close
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	req, err := http.NewRequest("POST", aiServiceEndpoint, request.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header = request.Header.Clone()
	req.Header.Set("Authorization", "Bearer "+info.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}

	io.Copy(w, resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}

	logInteraction(vk, info.Provider, request.Method, resp.StatusCode, time.Since(start), bodyBytes, responseBody)

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	resp.Body.Close()
}

func extractVirtualKey(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1], nil
	}
	return "", fmt.Errorf("invalid authorization header")
}

func logInteraction(virtualKey, provider, method string, status int, duration time.Duration, req, resp []byte) {
	logData := map[string]interface{}{
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"virtual_key": virtualKey,
		"provider":    provider,
		"method":      method,
		"status":      status,
		"duration_ms": duration.Milliseconds(),
		"request":     json.RawMessage(req),
		"response":    json.RawMessage(resp),
	}
	logJSON, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal log data: %v", err)
		return
	}

	fmt.Println(string(logJSON))
}
