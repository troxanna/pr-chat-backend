package main

import (
	"context"
	"fmt"
	"github.com/troxanna/pr-chat-backend/internal/bot"
	"os/signal"
	"syscall"
)

func main() {
	//cfg, err := config.Load()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	myBot, err := bot.NewBot("7661896995:AAE5P63Q-bBNvA1dZDmFrIeg31xtsOHkqvY")
	if err != nil {
		fmt.Println("Ошибка создания бота:", err)
		return
	}

	go func() {
		if err := myBot.Start(ctx); err != nil {
			fmt.Println("Ошибка работы бота:", err)
		}
		fmt.Println("Бот остановлен")
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Получен сигнал завершения, завершаем main")
			return
		default:
			//if err = application.New("pr", cfg).Run(); err != nil {
			//	fmt.Fprintln(os.Stderr, err)
			//	os.Exit(1)
			//}
			fmt.Println("ПУ ПУ ПУ")
		}
	}

}
