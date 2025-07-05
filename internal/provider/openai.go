package provider

import (
	"io"
	"net/http"
)

type OpenAI struct{}

func (o OpenAI) HandleRequest(apiKey string, bodyReader io.Reader, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header.Clone()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return http.DefaultClient.Do(req)
}
