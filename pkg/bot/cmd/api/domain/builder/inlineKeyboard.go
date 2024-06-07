package builder

import "github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"

type InlineKeyboardBuilder interface {
	GetReplyMarkup() *requests.ReplyMarkup
	AddRow(buttons ...InlineKeyboardButtonBuilder) InlineKeyboardBuilder
}

type InlineKeyboardButtonBuilder interface {
	GetInlineKeyboardButton() *requests.InlineKeyboardButton
	WithURL(url string) InlineKeyboardButtonBuilder
	WithCallbackData(data string) InlineKeyboardButtonBuilder
}

type button struct {
	Text         string
	URL          string
	CallbackData string
}

func (button *button) GetInlineKeyboardButton() *requests.InlineKeyboardButton {
	return &requests.InlineKeyboardButton{
		Text:         button.Text,
		URL:          button.URL,
		CallbackData: button.CallbackData,
	}
}

func (button *button) WithCallbackData(data string) InlineKeyboardButtonBuilder {
	button.CallbackData = data

	return button
}

func (button *button) WithURL(url string) InlineKeyboardButtonBuilder {
	button.URL = url

	return button
}

type inline struct {
	rows [][]*requests.InlineKeyboardButton
}

func (inline *inline) GetReplyMarkup() *requests.ReplyMarkup {
	return &requests.ReplyMarkup{
		InlineKeyboard: inline.rows,
	}
}

func (inline *inline) AddRow(buttons ...InlineKeyboardButtonBuilder) InlineKeyboardBuilder {
	b := []*requests.InlineKeyboardButton{}

	for i := range buttons {
		b = append(b, buttons[i].GetInlineKeyboardButton())
	}

	inline.rows = append(inline.rows, b)

	return inline
}

func NewInlineKeyboard() InlineKeyboardBuilder {
	return &inline{}
}

func NewInlineKeyboardButton(text string) InlineKeyboardButtonBuilder {
	return &button{
		text,
		"",
		"",
	}
}
