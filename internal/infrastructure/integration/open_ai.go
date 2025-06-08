package integration

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(
	httpClient *http.Client,
	baseURL string,
	apiKey string,
) Client {
	return Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (c Client) sendPromptForQuestion() {
	
}

