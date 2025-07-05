package provider

import (
	"fmt"
	"io"
	"net/http"
)

type Anthropic struct{}

func (a Anthropic) HandleRequest(apiKey string, bodyReader io.Reader, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = r.Header.Clone()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return http.DefaultClient.Do(req)
}
