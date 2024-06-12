package requests

import (
	"encoding/json"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/definitions"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/errors"
)

type ParseMode string

var MarkdownV2ParseMode ParseMode = "MarkdownV2"

type ReplyParameters struct {
	MessageId int `json:"message_id"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type SendMessageRequest struct {
	ChatId          string           `json:"chat_id"`
	Text            string           `json:"text"`
	ParseMode       ParseMode        `json:"parse_mode,omitempty"`
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup     `json:"reply_markup,omitempty"`
}

func (r *SendMessageRequest) Validate() error {
	if r.ChatId == "" {
		return errors.NewErrValidate("got empty ChatId")
	}

	return nil
}

func (req *SendMessageRequest) Serialize() ([]byte, error) {
	return json.Marshal(req)
}

func NewSendMessageRequest(
	chatId string,
	text string,
) definitions.Requester {
	return &SendMessageRequest{
		ChatId: chatId,
		Text:   text,
	}
}
