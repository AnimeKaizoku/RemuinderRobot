package telegram

import (
	"github.com/enrico5b1b4/tbwrap"
	tb "gopkg.in/tucnak/telebot.v2"
)

type TBWrapBot interface {
	Handle(path string, handler tbwrap.HandlerFunc)
	HandleButton(path *tb.InlineButton, handler tbwrap.HandlerFunc)
	HandleRegExp(path string, handler tbwrap.HandlerFunc)
	HandleMultiRegExp(paths []string, handler tbwrap.HandlerFunc)
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Start()
}
