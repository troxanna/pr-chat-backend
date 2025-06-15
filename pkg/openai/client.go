package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type ChatMessage struct {
	Role    string `json:"role"`   // "user", "system", "assistant"
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// SendMessage отправляет запрос к ChatGPT API и возвращает ответ.
func (c Client) SendMessage(ctx context.Context, model string, messages []ChatMessage) (string, error) {
	endpoint := c.baseURL + "/chat/completions"

	payload, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	setContentTypeJSON(req)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", requestError(resp)
	}

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from ChatGPT")
	}

	return response.Choices[0].Message.Content, nil
}

func requestError(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	return fmt.Errorf("error from ChatGPT: %s — %s", resp.Status, string(bodyBytes))
}

func setContentTypeJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}
