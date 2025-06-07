package main

import (
	"fmt"
	"os"

	"github.com/troxanna/pr-chat-backend/internal/config"
	"github.com/troxanna/pr-chat-backend/internal/application"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = application.New("pr", cfg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
