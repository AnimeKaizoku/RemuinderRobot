package remindlist

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	ReminderListRemoveCompletedRemindersBtn = "ReminderListRemoveCompletedRemindersBtn"
	ReminderListCloseCommandBtn             = "ReminderListCloseCommandBtn"
)

func NewButtons() map[string]*telebot.InlineButton {
	reminderListRemoveCompletedRemindersBtn := telebot.InlineButton{
		Unique: ReminderListRemoveCompletedRemindersBtn,
		Text:   "🗑 Remove completed reminders",
	}

	closeCommandBtn := telebot.InlineButton{
		Unique: ReminderListCloseCommandBtn,
		Text:   "❌ Close list",
	}

	return map[string]*telebot.InlineButton{
		ReminderListRemoveCompletedRemindersBtn: &reminderListRemoveCompletedRemindersBtn,
		ReminderListCloseCommandBtn:             &closeCommandBtn,
	}
}

func HandleReminderListRemoveCompletedRemindersBtn(reminderListService Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		err := c.Respond(c.Callback())
		if err != nil {
			return err
		}

		err = reminderListService.RemoveCompletedReminders(int(c.ChatID()))
		if err != nil {
			return err
		}

		_, err = c.Send("Completed reminders have been removed")

		return err
	}
}

func HandleCloseBtn() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		messageID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		err = c.Respond(c.Callback())
		if err != nil {
			return err
		}

		err = c.Delete(c.ChatID(), messageID)
		if err != nil {
			return err
		}

		return c.Delete(c.ChatID(), c.Message().ID)
	}
}
