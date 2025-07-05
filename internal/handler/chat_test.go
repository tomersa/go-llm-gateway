package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomersa/llm-gateway/internal/config"
	"github.com/tomersa/llm-gateway/internal/provider"
)

type testConfigEntry struct {
	APIKey   string
	Provider string
}

func setup() {
	// Setup a fake config
	config.Config = map[string]config.ProviderInfo{
		"testkey":     {APIKey: "real-api-key", Provider: "openai"},
		"badprovider": {APIKey: "real-api-key", Provider: "doesnotexist"},
	}
	provider.AiServiceEndpoints = map[string]string{
		"openai":    "http://mock-openai-endpoint",
		"anthropic": "http://mock-anthropic-endpoint",
	}
}

func TestHandleChat_ValidProvider(t *testing.T) {
	setup()
	// Mock http.DefaultClient.Do to avoid real HTTP calls to third part AI providers.
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"result": "ok"}`)),
				Header:     make(http.Header),
			}
		}),
	}
	defer func() { http.DefaultClient = oldClient }()

	body := bytes.NewBufferString(`{"message": "hi"}`)
	req := httptest.NewRequest("POST", "/chat", body)
	req.Header.Set("Authorization", "Bearer testkey")
	recorder := httptest.NewRecorder()

	HandleChat(recorder, req)

	if recorder.Code != 200 {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	var resp map[string]string
	json.NewDecoder(recorder.Body).Decode(&resp)
	if resp["result"] != "ok" {
		t.Errorf("unexpected response: %v", resp)
	}
}

func TestHandleChat_InvalidProvider(t *testing.T) {
	setup()
	body := bytes.NewBufferString(`{"message": "hi"}`)
	req := httptest.NewRequest("POST", "/chat", body)
	req.Header.Set("Authorization", "Bearer badprovider")
	recorder := httptest.NewRecorder()

	HandleChat(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandleChat_InvalidAuthHeader(t *testing.T) {
	setup()
	body := bytes.NewBufferString(`{"message": "hi"}`)
	req := httptest.NewRequest("POST", "/chat", body)
	// No Authorization header
	recorder := httptest.NewRecorder()

	HandleChat(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
