package handler

import (
	"context"

	"github.com/EXPORTER-DEV/go-telegram-bot/internal/navigation"
	"github.com/EXPORTER-DEV/go-telegram-bot/internal/router"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain/builder"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
)

func HomeHandler(requester api.Requester) router.HandlerFunc {
	return func(ctx context.Context, update *responses.Update) error {
		mb := builder.NewMessageBuilderReplyTo(update, "You are at home page!")

		mb.WithKeyboard(navigation.HomeKeyboard)

		return requester.SendMessage(ctx, mb)
	}
}

func RegisterHomeHandler(requester api.Requester, r router.Router) {
	r.AddByCallbackData(
		string(navigation.HomePath),
		HomeHandler(requester),
	)
}
