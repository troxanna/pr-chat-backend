package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/troxanna/pr-chat-backend/internal/infrastructure/integration"
	"gopkg.in/telebot.v4"
)

const messageQuestionTemplate = `Сформулируй один открытый вопрос на русском языке для собеседования, чтобы оценить уровень компетенции {skill} у сотрудника. Уровень указан как {level} по шкале от 0 до 5:
0 — Нет желания изучать
1 — Нет экспертизы. Не изучал и не применял на практике
2 — Средняя экспертиза. Изучал самостоятельно, практики было мало
3 — Хорошая экспертиза. Регулярно применяет на практике
4 — Эксперт. Знает тонкости, делится лайфхаками
5 — Гуру. Готов выступать на конференциях
Построй вопрос так, чтобы он был релевантен именно для уровня {level} и позволял раскрыть глубину знаний сотрудника. Используй профессиональный стиль.`

const messageResultTemplate = `Ты выступаешь в роли эксперта, оценивающего уровень профессиональной компетенции по ответу сотрудника.

Задача:
1. Дать оценку по шкале от 0 до 5 (только цифра)
2. Написать один краткий комментарий на русском языке — **обоснование**, почему именно этот уровень

Оценивай по шкале:
0 — Нет желания изучать  
1 — Нет экспертизы. Не изучал и не применял  
2 — Знает теоретически, но практики почти нет  
3 — Применяет регулярно, описывает кейсы  
4 — Глубокое понимание, делится практиками и тонкостями  
5 — Гуру, обучает других, решает сложные и редкие задачи

Будь строгим:
- Если нет конкретных действий, инструментов, решений — **не выше 2**
- Если ответ общий, без примеров — **макс. 2–3**
- Чтобы получить 4–5, должны быть ясные признаки экспертности и уникальности

Формат вывода:
Оценка: {0–5}  
Комментарий: {одно предложение, строго по содержанию ответа}
Ввод:  
Ответ сотрудника: {answer}`

var (
	skills = []string{"PostgreSQL", "MySQL/MariaDB", "ClickHouse", "MS SQL", "Redis", "MongoDB"}
	count  = 0
)

type HandlerFunc func(c telebot.Context) error

type BotWrapper struct {
	Bot      *telebot.Bot
	Handlers map[string]HandlerFunc
	Client   integration.Client
}

func NewBot(token string) (*BotWrapper, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	client := integration.NewClient(
		&http.Client{Transport: http.DefaultTransport},
		"app.cfg.ClientAI.BaseURL",
		"OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
	)

	return &BotWrapper{
		Bot:      bot,
		Handlers: make(map[string]HandlerFunc),
		Client:   client,
	}, nil
}

func (bw *BotWrapper) RegisterHandler(command string, handler HandlerFunc) {
	bw.Handlers[command] = handler
	bw.Bot.Handle(command, telebot.HandlerFunc(handler))
}

func (bw *BotWrapper) CommandHandlers() {
	startButton := telebot.InlineButton{Unique: "Start_PR", Text: "Launch Performance Review"}
	
	bw.RegisterHandler("/start", func(c telebot.Context) error {
		markup := &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{startButton},
				{{Text: "Launch Admin Space", URL: "http://10.10.169.1:8000/employee-competencies"}},
			},
		}
		count = 0
        bw.Client.CleanContextRequest(fmt.Sprintf("%d", c.Chat().ID))
		return c.Send("Добро пожаловать, для того чтобы начать Performance Review", markup)
	})

	bw.Bot.Handle(&startButton, func(c telebot.Context) error {
		userID := c.Chat().ID
		// defer bw.Client.CleanContextRequest(userID)
		msg, _ := bw.SendQuestion(2, userID)
		return c.Send(msg)
	})

	bw.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := fmt.Sprintf("%d", c.Chat().ID)
		// defer bw.Client.CleanContextRequest(userID)

		answer := c.Text()
		message := strings.ReplaceAll(messageResultTemplate, "{answer}", answer)
		bw.Client.SendPromptForQuestion(userID, message)

		var result string
		for ok := false; !ok; {
			ok, result = bw.Client.GetResultForQuestionRequest(userID)
		}

		if err := c.Send(result); err != nil {
			return err
		}
		msg, _ := bw.SendQuestion(2, c.Chat().ID)
		return c.Send(msg)
	})

	bw.RegisterHandler(telebot.OnVoice, func(c telebot.Context) error {
		voice := c.Message().Voice
		reader, err := bw.Bot.File(voice.MediaFile())
		if err != nil {
			log.Printf("Не удалось получить файл: %v", err)
			return c.Send("Произошла ошибка при получении голосового сообщения.")
		}
		defer reader.Close()

		answ, err := sendVoiceInWhisper(reader, voice.FileID+".ogg")
		if err != nil {
			log.Printf("Ошибка отправки на сервис: %v", err)
			return c.Send("Ошибка при отправке голосового сообщения.")
		}

		userID := fmt.Sprintf("%d", c.Chat().ID)
		// defer bw.Client.CleanContextRequest(userID)

		message := strings.ReplaceAll(messageResultTemplate, "{answer}", answ)
		bw.Client.SendPromptForQuestion(userID, message)

		var result string
		for ok := false; !ok; {
			ok, result = bw.Client.GetResultForQuestionRequest(userID)
		}

		if err := c.Send(answ); err != nil {
			return err
		}
		if err := c.Send(result); err != nil {
			return err
		}
		msg, _ := bw.SendQuestion(2, c.Chat().ID)
		return c.Send(msg)
	})
}

func (bw *BotWrapper) SendQuestion(level int, userId int64) (string, error) {
	if count == len(skills) {
		count = 0
        bw.Client.CleanContextRequest(fmt.Sprintf("%d", userId))
		return "Твой уровень по группе компетенций Базы данных: 3.2", nil
	}
	message := strings.ReplaceAll(messageQuestionTemplate, "{skill}", skills[count])
	message = strings.ReplaceAll(message, "{level}", fmt.Sprintf("%d", level))
	bw.Client.SendPromptForQuestion(fmt.Sprintf("%d", userId), message)

	var result string
	for ok := false; !ok; {
		ok, result = bw.Client.GetResultForQuestionRequest(fmt.Sprintf("%d", userId))
	}
	count++
	return result, nil
}

func (bw *BotWrapper) Start(ctx context.Context) error {
	bw.CommandHandlers()
	errCh := make(chan error, 1)
	go bw.Bot.Start()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
