package bot

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/troxanna/pr-chat-backend/internal/infrastructure/integration"
	"gopkg.in/telebot.v4"
	"github.com/google/uuid"
)

var messageQuestion = `Сформулируй один открытый вопрос для собеседования, чтобы оценить уровень компетенции PosgreSql у сотрудника. Уровень указан как 2 по следующей шкале:
0 — Нет желания изучать
1 — Нет экспертизы. Не изучал и не применял на практике
2 — Средняя экспертиза. Изучал самостоятельно, практики было мало
3 — Хорошая экспертиза. Регулярно применяет на практике
4 — Эксперт. Знает тонкости, делится лайфхаками
5 — Гуру. Готов выступать на конференциях
Построй вопрос так, чтобы он был релевантен именно для уровня 2 и позволял раскрыть глубину знаний сотрудника. Используй профессиональный стиль.
)`

var messageResult = `Ты выступаешь в роли эксперта, оценивающего уровень профессиональной компетенции по ответу сотрудника.

Задача:
1. Дать оценку по шкале от 0 до 5 (только цифра)
2. Написать один краткий комментарий — **обоснование**, почему именно этот уровень

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
Ответ сотрудника:  `
// var clientAI integration.Client

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

	bw := &BotWrapper{
		Bot:      bot,
		Handlers: make(map[string]HandlerFunc),
		Client: integration.NewClient(
				&http.Client{Transport: http.DefaultTransport},
				"app.cfg.ClientAI.BaseURL",
				"OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
			),
	}
	log.Println(bw.Client)
	log.Println("test3")

	return bw, nil
}

func (bw *BotWrapper) RegisterHandler(command string, handler HandlerFunc) {
	bw.Handlers[command] = handler
	bw.Bot.Handle(command, telebot.HandlerFunc(handler))
}

type Competency struct {
	level     string
	name      string
	Questions map[int][]string
}

func (bw *BotWrapper) CommandHandlers() {
	log.Println(bw.Client)
	log.Println("test1")
	startButton := telebot.InlineButton{
		Unique: "Start_PR",
		Text:   "Launch Performance Review",
	}
	bw.RegisterHandler("/start", func(c telebot.Context) error {
		return c.Send("Добро пожаловать, для того чтобы начать Performance Review",
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{startButton},
					{{
						Text: "Launch Admin Space",
						WebApp: &telebot.WebApp{
							URL: "https://www.appsmith.com/",
						},
					}},
				},
			})
	})
	bw.Bot.Handle(&startButton, func(c telebot.Context) error {
		uid := uuid.NewString()
		defer bw.Client.CleanContextRequest(uid)
		log.Println(bw.Client)
		log.Println("test2")
		bw.Client.SendPromptForQuestion(uid, messageQuestion)
		result := false
		mes := ""
		for !result {
			result, mes = bw.Client.GetResultForQuestionRequest(uid)
		}
		
		return c.Send(mes)
	})

	bw.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		tmp := fmt.Sprintf("%s\n%s",messageResult,c.Update().Message.Text)
		uid := uuid.NewString()
		defer bw.Client.CleanContextRequest(uid)
		log.Println(bw.Client)
		log.Println("test2")
		bw.Client.SendPromptForQuestion(uid, tmp)
		result := false
		mes := ""
		for !result {
			result, mes = bw.Client.GetResultForQuestionRequest(uid)
		}
		
		return c.Send(mes)
	})

	// bw.Bot.Handle(telebot.OnAudio, func(c telebot.Context) error {
	// 	uid := uuid.NewString()
	// 	defer bw.Client.CleanContextRequest(uid)
	// 	log.Println(bw.Client)
	// 	log.Println("test2")
	// 	bw.Client.SendPromptForQuestion(uid)
	// 	result := false
	// 	mes := ""
	// 	for !result {
	// 		result, mes = bw.Client.GetResultForQuestionRequest(uid)
	// 	}
		
	// 	return c.Send(mes)
	// })
}

func (bw *BotWrapper) Start(ctx context.Context) error {
	bw.CommandHandlers()
	errCh := make(chan error, 1)

	go func() {
		// clientAI = integration.NewClient(
		// 	&http.Client{Transport: http.DefaultTransport},
		// 	"app.cfg.ClientAI.BaseURL",
		// 	"OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
		// )
		bw.Bot.Start()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
