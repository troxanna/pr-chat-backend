package integration


import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	tg "github.com/troxanna/pr-chat-backend/pkg/bot"
)

const messageQuestionTemplate = `Сформулируй один открытый вопрос на русском языке для собеседования, чтобы оценить уровень компетенции REST API у сотрудника. Уровень указан как 2 по шкале от 0 до 5:
0 — Нет желания изучать
1 — Нет экспертизы. Не изучал и не применял на практике
2 — Средняя экспертиза. Изучал самостоятельно, практики было мало
3 — Хорошая экспертиза. Регулярно применяет на практике
4 — Эксперт. Знает тонкости, делится лайфхаками
5 — Гуру. Готов выступать на конференциях
Построй вопрос так, чтобы он был релевантен именно для уровня 2 и позволял раскрыть глубину знаний сотрудника. Используй профессиональный стиль.`

type TelegramBotService struct {
	client     *tg.Client
	chatGPTSvc ChatGPTService // ← Интеграция с ChatGPT
}

func NewTelegramBotService(client *tg.Client, chatGPTSvc ChatGPTService) TelegramBotService {
	return TelegramBotService{
		client:     client,
		chatGPTSvc: chatGPTSvc,
	}
}

func (s TelegramBotService) Start(ctx context.Context) error {
	updates := s.client.UpdatesChan()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:

			switch {
			// Обработка команд
			case update.Message != nil:
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я твой Telegram-бот.")
						msg.ReplyMarkup = s.buildStartKeyboard()
						s.client.API().Send(msg)
					default:
						s.client.SendMessage(update.Message.Chat.ID, "Неизвестная команда.")
					}
				} else {
					s.client.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Вы сказали: %s", update.Message.Text))
				}

			// Обработка кнопок
			case update.CallbackQuery != nil:
				data := update.CallbackQuery.Data
				switch data {
				case "fill_form":
					go s.handleFormRequest(ctx, update.CallbackQuery)
				default:
					s.client.SendMessage(update.CallbackQuery.Message.Chat.ID, "Неизвестное действие.")
				}
			}
		}
	}
}

func (s TelegramBotService) buildStartKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("📝 Заполнить анкету", "fill_form")
	row := tgbotapi.NewInlineKeyboardRow(button)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func (s TelegramBotService) handleFormRequest(ctx context.Context, q *tgbotapi.CallbackQuery) {
	chatID := q.Message.Chat.ID

	// Уведомление Telegram'у, что нажата кнопка
	_, err := s.client.API().Request(tgbotapi.NewCallback(q.ID, "Формируем анкету..."))

	if err != nil {
		log.Printf("failed to confirm callback: %v", err)
	}

	// Вызов ChatGPT
	response, err := s.chatGPTSvc.AskUser(ctx, messageQuestionTemplate)
	if err != nil {
		log.Printf("ChatGPT error: %v", err)
		s.client.SendMessage(chatID, "Ошибка при запросе анкеты.")
		return
	}

	s.client.SendMessage(chatID, "Вот ваша анкета:\n\n"+response)
}
