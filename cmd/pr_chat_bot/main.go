package main

import (
	"fmt"
	"os"

	"github.com/troxanna/pr-chat-backend/internal/config"
)

func main() {
	_, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// if err = application.New(appName, appVersion, cfg).Run(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
}
