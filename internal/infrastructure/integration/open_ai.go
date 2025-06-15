package integration

import (
	"context"
	"fmt"

	"github.com/troxanna/pr-chat-backend/pkg/openai"
)

type ChatGPTService struct {
	client openai.Client
	model  string
}

func NewChatGPTService(
	client openai.Client,
	model string,
) ChatGPTService {
	return ChatGPTService{
		client: client,
		model:  model,
	}
}

// AskUser отправляет вопрос от пользователя в ChatGPT и возвращает ответ.
func (s ChatGPTService) AskUser(ctx context.Context, prompt string) (string, error) {
	messages := []openai.ChatMessage{
		{Role: "user", Content: prompt},
	}

	response, err := s.client.SendMessage(ctx, s.model, messages)
	if err != nil {
		return "", fmt.Errorf("client.SendMessage: %w", err)
	}

	return response, nil
}