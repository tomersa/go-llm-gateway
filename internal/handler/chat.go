package handler

import (
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

	var aiService provider.AiService
	switch info.Provider {
	case "openai":
		aiService = provider.OpenAI{}
	case "anthropic":
		aiService = provider.Anthropic{}
	default:
		http.Error(w, fmt.Sprintf("Unsupported provider: %s", info.Provider), http.StatusBadRequest)
		return
	}

	resp, err := aiService.HandleRequest(info.APIKey, request.Body, request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
}

func extractVirtualKey(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1], nil
	}
	return "", fmt.Errorf("invalid authorization header")
}
