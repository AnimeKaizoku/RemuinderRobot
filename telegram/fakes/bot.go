package fakes

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	handler              map[string]func(m *tb.Message)
	OutboundSendMessages []string
}

func NewBot() *Bot {
	return &Bot{
		handler: make(map[string]func(m *tb.Message)),
	}
}

func (t *Bot) Handle(endpoint, h interface{}) {
	if handler, ok := h.(func(*tb.Message)); ok {
		t.handler[endpoint.(string)] = handler
		return
	}
}

func (t *Bot) Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error {
	return nil
}

func (t *Bot) Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	if message, ok := what.(string); ok {
		t.OutboundSendMessages = append(t.OutboundSendMessages, message)
		return nil, nil
	}
	return nil, nil
}

func (t *Bot) Delete(chatID int64, messageID int) error {
	return nil
}

func (t *Bot) Start() {}

func (t *Bot) SimulateIncomingMessageToChat(chatID int, text string) {
	if handler, ok := t.handler[text]; ok {
		handler(&tb.Message{Text: text, Chat: &tb.Chat{ID: int64(chatID)}})
		return
	}

	t.handler[tb.OnText](&tb.Message{Text: text, Chat: &tb.Chat{ID: 1}})
}
