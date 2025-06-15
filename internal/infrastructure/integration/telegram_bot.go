package integration

import (
	"context"
	"fmt"
	tg "github.com/troxanna/pr-chat-backend/pkg/bot"
)

type TelegramBotService struct {
	client *tg.Client
}

func NewTelegramBotService(client *tg.Client) TelegramBotService {
	return TelegramBotService{client: client}
}

func (s TelegramBotService) Start(ctx context.Context) error {
	updates := s.client.UpdatesChan()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			if update.Message != nil {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						s.client.SendMessage(update.Message.Chat.ID, "Привет! Я твой Telegram-бот.")
					default:
						s.client.SendMessage(update.Message.Chat.ID, "Неизвестная команда.")
					}
				} else {
					reply := fmt.Sprintf("Вы сказали: %s", update.Message.Text)
					s.client.SendMessage(update.Message.Chat.ID, reply)
				}
			}
		}
	}
}
