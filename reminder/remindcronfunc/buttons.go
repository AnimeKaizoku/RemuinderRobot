package remindcronfunc

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	Snooze15MinuteBtn = "Snooze15MinuteBtn"
	Snooze30MinuteBtn = "Snooze30MinuteBtn"
	Snooze1HourBtn    = "Snooze1HourBtn"
)

func NewButtons() map[string]*telebot.InlineButton {
	snooze15MinuteBtn := telebot.InlineButton{
		Unique: Snooze15MinuteBtn,
		Text:   "⏰ 15m",
	}
	snooze30MinuteBtn := telebot.InlineButton{
		Unique: Snooze30MinuteBtn,
		Text:   "⏰ 30m",
	}
	snooze1HourBtn := telebot.InlineButton{
		Unique: Snooze1HourBtn,
		Text:   "⏰ 1h",
	}

	return map[string]*telebot.InlineButton{
		Snooze15MinuteBtn: &snooze15MinuteBtn,
		Snooze30MinuteBtn: &snooze30MinuteBtn,
		Snooze1HourBtn:    &snooze1HourBtn,
	}
}

type RemindDateServicer interface {
	AddReminderIn(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (reminder.NextScheduleChatTime, error)
}

func HandleReminderSnoozeBtn(
	service RemindDateServicer,
	store reminder.Storer,
	amountDateTime reminder.AmountDateTime,
) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		err := c.Respond(c.Callback())
		if err != nil {
			return err
		}

		chatID := int(c.ChatID())
		reminderID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		rem, err := store.GetReminder(chatID, reminderID)
		if err != nil {
			return err
		}

		nextSchedule, err := service.AddReminderIn(chatID, rem.Data.Command, amountDateTime, rem.Data.Message)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Reminder \"%s\" has been rescheduled for %s",
			rem.Data.Message,
			nextSchedule.Time.In(nextSchedule.Location).Format("Mon, 02 Jan 2006 15:04 MST"),
		))

		return err
	}
}
