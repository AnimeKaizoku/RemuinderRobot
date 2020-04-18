package fakes

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type TeleBot struct {
	handler              map[string]func(m *tb.Message)
	OutboundSendMessages []string
}

func NewTeleBot() *TeleBot {
	return &TeleBot{
		handler: make(map[string]func(m *tb.Message)),
	}
}

func (t *TeleBot) Handle(endpoint, h interface{}) {
	if handler, ok := h.(func(*tb.Message)); ok {
		t.handler[endpoint.(string)] = handler
		return
	}
}

func (t *TeleBot) Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error {
	return nil
}

func (t *TeleBot) Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	if message, ok := what.(string); ok {
		t.OutboundSendMessages = append(t.OutboundSendMessages, message)
		return nil, nil
	}
	return nil, nil
}

func (t *TeleBot) Delete(message tb.Editable) error {
	return nil
}

func (t *TeleBot) Start() {}

func (t *TeleBot) SimulateIncomingMessageToChat(chatID int64, text string) {
	if handler, ok := t.handler[text]; ok {
		handler(&tb.Message{Text: text, Chat: &tb.Chat{ID: chatID}})
		return
	}

	t.handler[tb.OnText](&tb.Message{Text: text, Chat: &tb.Chat{ID: chatID}})
}
