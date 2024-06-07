package handler

import (
	"context"

	"github.com/EXPORTER-DEV/go-telegram-bot/internal/router"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
)

func RegisterNotFoundHandler(requester api.Requester, r router.Router) {
	r.AddNotFoundHandler(
		func(ctx context.Context, update *responses.Update) error {
			requester.ReplyTo(ctx, update, "I have not found requested resource, forward you to home page")
			return HomeHandler(requester)(ctx, update)
		},
	)
}
