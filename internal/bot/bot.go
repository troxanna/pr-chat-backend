package bot

import (
  "context"
  "fmt"
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

type Competency struct {
  level     int
  name      string
  Questions map[int][]string
}

func (bw *BotWrapper) CommandHandlers() {
  groupCompetency := map[string]Competency{
    "PostgresSQL": Competency{
      level:     0,
      name:      "",
      Questions: nil,
    },
  }
  startButton := telebot.InlineButton{
    Unique: "Start_PR",
    Text:   "Launch Performance Review",
  }
  bw.RegisterHandler("/start", func(c telebot.Context) error {
    return c.Send("Добро пожаловать, для того чтобы начать Performance Review",
      &telebot.ReplyMarkup{
        InlineKeyboard: [][]telebot.InlineButton{
          {startButton},
        },
      })
  })
  bw.Bot.Handle(&startButton, func(c telebot.Context) error {
    c.Send("Начинаем Performance Review 🎯" + "\nЗаполни карту компетенций")
    for i, _ := range groupCompetency {
      c.Send(fmt.Sprintf("Как ты оцениваешь свои знания в %s?", groupCompetency[i].name),
        &telebot.ReplyMarkup{
          InlineKeyboard: [][]telebot.InlineButton{
            {
              {
                Unique: "level_1",
              },
              {
                Unique: "level_2",
              },
              {
                Unique: "level_3",
              },
              {
                Unique: "level_4",
              },
              {
                Unique: "level_5",
              },
            },
          },
        },
      )

      bw.Bot.Handle("level_1", func(c telebot.Context) error {
        return nil
      })

    }
    return nil
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