package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/EXPORTER-DEV/go-telegram-bot/internal/handler"
	"github.com/EXPORTER-DEV/go-telegram-bot/internal/router"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/joho/godotenv"
)

var TokenEnvKey = "TOKEN"

func loadConfigToEnvironment() {
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		log.Fatalf("Failed to get working directory")
	}

	log.Print(filename)

	p := path.Join(filename, "../../../", ".env")

	log.Print(p)

	if err := godotenv.Load(p); err != nil {
		log.Fatalf("Got failed load .env file: %v\n", err)
	}
}

func bootstap() api.Requester {
	token := os.Getenv(TokenEnvKey)
	if token == "" {
		log.Fatalf("Got not set: %s key in environment", TokenEnvKey)
	}

	b, err := bot.New(token)

	if err != nil {
		log.Fatalf("Error while init bot: %v", err)
	}

	return b
}

func polling(ctx context.Context, requestor api.Requester) {
	ch := requestor.Poll(ctx)

	router := router.New()

	handler.RegisterAboutMeHandler(requestor, router)
	handler.RegisterHomeHandler(requestor, router)
	handler.RegisterNotFoundHandler(requestor, router)

	for update := range ch {
		found, err := router.Handle(ctx, update)

		if !found {
			log.Printf("Got not found handler for update: %+v, ignore it\n", update)
		}

		if err != nil {
			log.Printf("Error: %+v\n", fmt.Errorf("got error while handling router: %w", err))
		}
	}
}

func main() {
	loadConfigToEnvironment()

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
