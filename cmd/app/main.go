package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain"
)

func bootstap() api.TelegramAPIInterface {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("Got no TOKEN in env")
	}

	b, err := bot.New(os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return b
}

func polling(ctx context.Context, b api.TelegramAPIInterface) {
	ch := b.Poll(ctx)

	for update := range ch {
		m := domain.NewMessageBuilder(strconv.Itoa(update.Message.Chat.Id), "text")
		m.WithReplyToMessageId(update.Message.Id)

		err := b.Reply(ctx, update, "test")

		if err != nil {
			fmt.Printf("Failed to send message back: %v", err)
		}
	}
}

func main() {
	b := bootstap()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go polling(ctx, b)

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt)
	signal.Notify(exit, syscall.SIGTERM)

	<-exit

	cancel()
}
