package telegram

import (
	"github.com/enrico5b1b4/tbwrap"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot interface {
	Handle(path string, handler tbwrap.HandlerFunc)
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Start()
}
