package remindcronfunc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/telegram"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

// New creates a function which is called when a reminder is due
func New(s Servicer, b telegram.TBWrapBot, r *reminder.Reminder) func() {
	return func() {
		buttons := NewButtons()
		var inlineKeys [][]telebot.InlineButton

		snooze15MinuteBtn := *buttons[Snooze15MinuteBtn]
		snooze15MinuteBtn.Data = strconv.Itoa(r.ID)
		snooze30MinuteBtn := *buttons[Snooze30MinuteBtn]
		snooze30MinuteBtn.Data = strconv.Itoa(r.ID)
		snooze1HourBtn := *buttons[Snooze1HourBtn]
		snooze1HourBtn.Data = strconv.Itoa(r.ID)
		inlineKeys = append(inlineKeys, []telebot.InlineButton{snooze15MinuteBtn, snooze30MinuteBtn, snooze1HourBtn})

		messageWithIcon := fmt.Sprintf("🗓 %s", r.Data.Message)
		_, err := b.Send(&tb.Chat{ID: int64(r.Data.RecipientID)}, messageWithIcon, &telebot.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
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
