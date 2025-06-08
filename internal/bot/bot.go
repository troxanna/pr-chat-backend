package bot

import (
	"context"
	"gopkg.in/telebot.v4"
	"time"
)

type HandlerFunc func(c telebot.Context) error

type BotWrapper struct {
	Bot      *telebot.Bot
	Handlers map[string]HandlerFunc
}

func NewBot(token string) (*BotWrapper, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	bw := &BotWrapper{
		Bot:      bot,
		Handlers: make(map[string]HandlerFunc),
	}

	return bw, nil
}

func (bw *BotWrapper) RegisterHandler(command string, handler HandlerFunc) {
	bw.Handlers[command] = handler
	bw.Bot.Handle(command, telebot.HandlerFunc(handler))
}

func (bw *BotWrapper) CommandHandlers() {
	bw.RegisterHandler("/start", func(c telebot.Context) error {
		return c.Send("Добро пожаловать, для того чтобы начать Performance Review",
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{{
						Text:   "Launch Performance Review",
						WebApp: &telebot.WebApp{URL: "https://www.appsmith.com/"},
					}},
				},
			})
	})
}

func (bw *BotWrapper) Start(ctx context.Context) error {
	bw.CommandHandlers()
	errCh := make(chan error, 1)

	go func() {
		bw.Bot.Start()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
