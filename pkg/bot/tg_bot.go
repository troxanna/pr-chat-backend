package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	api *tgbotapi.BotAPI
}

func NewClient(token string) (*Client, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Client{api: api}, nil
}

func (c *Client) UpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return c.api.GetUpdatesChan(u)
}

func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.api.Send(msg)
	return err
}

func (c *Client) API() *tgbotapi.BotAPI {
	return c.api
}