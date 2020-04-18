package remindcronfunc

import (
	"fmt"
	"log"

	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
)

func New(s Servicer, b telegram.TBWrapBot, r *reminder.Reminder) func() {
	return func() {
		messageWithIcon := fmt.Sprintf("🗓 %s", r.Data.Message)
		_, err := b.Send(&tb.Chat{ID: int64(r.Data.RecipientID)}, messageWithIcon)
		if err != nil {
			log.Printf("NewReminderCronFunc err: %q", err)
			return
		}

		if !r.Job.RunOnlyOnce {
			return
		}

		err = s.Complete(r)
		if err != nil {
			log.Printf("NewReminderCronFunc complete err: %q", err)
			return
		}

		if r.Job.RepeatSchedule != nil {
			err := s.AddReminderRepeatSchedule(r)
			if err != nil {
				log.Printf("NewReminderCronFunc AddReminderRepeatSchedule err: %q", err)
				return
			}
		}
	}
}
