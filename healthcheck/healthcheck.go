package healthcheck

import (
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Healthcheck struct {
	cron.Job
	Data HealthcheckData `json:"data"`
}

type HealthcheckData struct {
	RecipientID int    `json:"recipient_id"`
	Message     string `json:"message"`
}

func NewHealthcheckCronFunc(b *bot.Bot, h *Healthcheck) func() {
	return func() {
		b.Send(&tb.User{ID: h.Data.RecipientID}, h.Data.Message)
	}
}
