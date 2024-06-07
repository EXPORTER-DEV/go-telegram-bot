package router

import (
	"context"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
)

type HandlerFunc func(ctx context.Context, update *responses.Update) error

type Router interface {
	AddByCallbackData(data string, handler HandlerFunc) Router
	AddByMessageText(text string, handler HandlerFunc) Router
	AddNotFoundHandler(handler HandlerFunc) Router
	Handle(ctx context.Context, update *responses.Update) (found bool, err error)
}

type router struct {
	byCallbackData  map[string]HandlerFunc
	byMessageText   map[string]HandlerFunc
	notFoundHandler HandlerFunc
}

func (r *router) AddByCallbackData(data string, handler HandlerFunc) Router {
	r.byCallbackData[data] = handler

	return r
}

func (r *router) AddByMessageText(text string, handler HandlerFunc) Router {
	r.byMessageText[text] = handler

	return r
}

func (r *router) AddNotFoundHandler(handler HandlerFunc) Router {
	r.notFoundHandler = handler

	return r
}

func (r *router) Handle(ctx context.Context, update *responses.Update) (found bool, err error) {
	if update.CallbackQuery != nil {
		v, ok := r.byCallbackData[update.CallbackQuery.Data]

		if ok {
			return true, v(ctx, update)
		}
	}

	if update.Message != nil {
		v, ok := r.byMessageText[update.Message.Text]

		if ok {
			return true, v(ctx, update)
		}
	}

	if r.notFoundHandler != nil {
		return true, r.notFoundHandler(ctx, update)
	}

	return false, nil
}

func New() Router {
	r := router{}

	r.byCallbackData = make(map[string]HandlerFunc)
	r.byMessageText = make(map[string]HandlerFunc)

	return &r
}
