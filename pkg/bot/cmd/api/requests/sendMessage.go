package requests

import "encoding/json"

type ParseMode string

var MarkdownV2ParseMode ParseMode = "MarkdownV2"

type ReplyParameters struct {
	MessageId int `json:"message_id"`
}

type SendMessageRequest struct {
	ChatId          string           `json:"chat_id"`
	Text            string           `json:"text"`
	ParseMode       ParseMode        `json:"parse_mode,omitempty"`
	ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty"`
}

func (req *SendMessageRequest) Serialize() ([]byte, error) {
	return json.Marshal(req)
}

func NewSendMessageRequest(chatId string, text string) *SendMessageRequest {
	return &SendMessageRequest{
		ChatId: chatId,
		Text:   text,
	}
}
