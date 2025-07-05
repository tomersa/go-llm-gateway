package provider

import (
	"io"
	"net/http"
)

type AiService interface {
	HandleRequest(apiKey string, bodyReader io.Reader, r *http.Request) (*http.Response, error)
}
