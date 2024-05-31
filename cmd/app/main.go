package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("Got no TOKEN in env")
	}

	b, err := bot.New(os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		exit := make(chan os.Signal, 1)

		signal.Notify(exit, os.Interrupt)
		signal.Notify(exit, os.Kill)

		<-exit

		cancel()
	}()

	b.Poll(ctx)
}
