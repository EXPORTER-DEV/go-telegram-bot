package handler

import (
	"context"

	"github.com/EXPORTER-DEV/go-telegram-bot/internal/navigation"
	"github.com/EXPORTER-DEV/go-telegram-bot/internal/router"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain/builder"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
)

func AboutMeHandler(requester api.Requester) router.HandlerFunc {
	return func(ctx context.Context, update *responses.Update) error {
		mb := builder.NewMessageBuilderReplyTo(update, "I'm just a Telegram bot written on: https://github.com/EXPORTER-DEV/go-telegram-bot")

		mb.WithKeyboard(navigation.GoToHomeKeyboard)

		return requester.SendMessage(ctx, mb)
	}
}

func RegisterAboutMeHandler(requester api.Requester, r router.Router) {
	r.AddByCallbackData(
		string(navigation.AboutMePath),
		AboutMeHandler(requester),
	)
}
