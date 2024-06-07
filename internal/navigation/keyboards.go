package navigation

import "github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain/builder"

type Path string

var HomePath Path = "HOME"
var AboutMePath Path = "ABOUT_ME"

var HomeKeyboard = builder.NewInlineKeyboard().AddRow(
	builder.NewInlineKeyboardButton("About me").WithCallbackData(string(AboutMePath)),
)

var GoToHomeKeyboard = builder.NewInlineKeyboard().AddRow(
	builder.NewInlineKeyboardButton("Home").WithCallbackData(string(HomePath)),
)
