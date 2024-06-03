package domain

import "github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"

type MessageBuilder struct {
	req *requests.SendMessageRequest
}

func (m *MessageBuilder) WithParseMode(parseMode requests.ParseMode) {
	m.req.ParseMode = parseMode
}

func (m *MessageBuilder) WithReplyToMessageId(messageId int) {
	m.req.ReplyParameters = &requests.ReplyParameters{
		MessageId: messageId,
	}
}

func (m *MessageBuilder) GetRequest() *requests.SendMessageRequest {
	return m.req
}

func NewMessageBuilder(chatId string, text string) *MessageBuilder {
	return &MessageBuilder{
		&requests.SendMessageRequest{
			ChatId: chatId,
			Text:   text,
		},
	}
}
