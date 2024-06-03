package responses

import "github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/definitions"

type InlineQuery struct {
	Id       string            `json:"id"`
	From     *definitions.User `json:"user"`
	Query    string            `json:"query"`
	ChatType string            `json:"chat_type"`
}

type Update struct {
	UpdateId    int                  `json:"update_id"`
	Message     *definitions.Message `json:"message"`
	InlineQuery *InlineQuery         `json:"inline_query"`
}

type GetUpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}
