package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tomersa/llm-gateway/internal/config"
	"github.com/tomersa/llm-gateway/internal/provider"
)

func HandleChat(w http.ResponseWriter, request *http.Request) {
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

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func extractVirtualKey(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1], nil
	}
	return "", fmt.Errorf("invalid authorization header")
}
