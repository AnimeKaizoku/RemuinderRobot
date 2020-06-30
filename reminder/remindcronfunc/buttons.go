package remindcronfunc

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	Snooze20MinuteBtn = "Snooze20MinuteBtn"
	Snooze1HourBtn    = "Snooze1HourBtn"
	CompleteBtn       = "CompleteBtn"
)

func NewButtons() map[string]*telebot.InlineButton {
	snooze20MinuteBtn := telebot.InlineButton{
		Unique: Snooze20MinuteBtn,
		Text:   "⏰ 20m",
	}
	snooze1HourBtn := telebot.InlineButton{
		Unique: Snooze1HourBtn,
		Text:   "⏰ 1h",
	}
	completeBtn := telebot.InlineButton{
		Unique: CompleteBtn,
		Text:   "✅ Done",
	}

	return map[string]*telebot.InlineButton{
		Snooze20MinuteBtn: &snooze20MinuteBtn,
		Snooze1HourBtn:    &snooze1HourBtn,
		CompleteBtn:       &completeBtn,
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

func HandleReminderCompleteBtn(
	service Servicer,
	store reminder.Storer,
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

		err = service.Complete(rem)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Reminder \"%s\" has been completed",
			rem.Data.Message,
		))

		return err
	}
}
