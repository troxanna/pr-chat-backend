package bot

import (
	"context"
	"fmt"
	"github.com/troxanna/pr-chat-backend/internal/config"
	"log"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v4"
)

func StartBot(ctx context.Context, cfg *config.Config, router *gin.Engine) error {
	webhook := &tele.Webhook{
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: cfg.Telegram.WebhookUrl,
			Cert:      cfg.Telegram.Cert,
		},
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:     cfg.Telegram.BotToken,
		Poller:    webhook,
		ParseMode: tele.ModeHTML,
		Verbose:   true,
	})
	if err != nil {
		return err
	}

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет, для того чтобы начать Performance Review Нажми на кнопку Launch PR", tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{
					Text: "Launch Bot",
					URL:  fmt.Sprintf("%s+%s", cfg.HTTP.ListenAddress, "/webhook"),
				}},
			},
		})
	})

	router.POST("/webhook", func(c *gin.Context) {
		webhook.ServeHTTP(c.Writer, c.Request)
	})
	go bot.Start()

	<-ctx.Done()

	log.Println("Shutting down bot...")

	bot.Stop()
	return nil
}
