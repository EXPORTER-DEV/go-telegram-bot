package builder

import (
	"strconv"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
)

type MessageBuilder interface {
	WithParseMode(parseMode requests.ParseMode) MessageBuilder
	WithReplyToMessageId(messageId int) MessageBuilder
	WithKeyboard(keyboard KeyboardBuilder) MessageBuilder
	GetRequest() *requests.SendMessageRequest
}

type messageBuilder struct {
	req *requests.SendMessageRequest
}

func (m *messageBuilder) WithParseMode(parseMode requests.ParseMode) MessageBuilder {
	m.req.ParseMode = parseMode

	return m
}

func (m *messageBuilder) WithReplyToMessageId(messageId int) MessageBuilder {
	m.req.ReplyParameters = &requests.ReplyParameters{
		MessageId: messageId,
	}

	return m
}

func (m *messageBuilder) WithKeyboard(keybord KeyboardBuilder) MessageBuilder {
	m.req.ReplyMarkup = *keybord.GetReplyMarkup()

	return m
}

func (m *messageBuilder) GetRequest() *requests.SendMessageRequest {
	return m.req
}

func NewMessageBuilder(chatId string, text string) MessageBuilder {
	return &messageBuilder{
		&requests.SendMessageRequest{
			ChatId: chatId,
			Text:   text,
		},
	}
}

func NewMessageBuilderReplyTo(update *responses.Update, text string) MessageBuilder {
	var chatId string

	if update.Message != nil {
		chatId = strconv.Itoa(update.Message.Chat.Id)
	}

	if update.CallbackQuery != nil {
		chatId = strconv.Itoa(update.CallbackQuery.From.Id)
	}

	return &messageBuilder{
		&requests.SendMessageRequest{
			ChatId: chatId,
			Text:   text,
		},
	}
}
