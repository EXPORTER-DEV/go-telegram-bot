package builder

import "github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"

type KeyboardBuilder interface {
	GetReplyMarkup() *requests.ReplyMarkup
}
