package reminddetail

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	ReminderDetailDeleteBtn              = "ReminderDetailDeleteBtn"
	ReminderDetailShowReminderCommandBtn = "ReminderDetailShowReminderCommandBtn"
	ReminderDetailCloseCommandBtn        = "ReminderDetailCloseCommandBtn"
)

func NewButtons() map[string]*telebot.InlineButton {
	reminderDetailDeleteBtn := telebot.InlineButton{
		Unique: ReminderDetailDeleteBtn,
		Text:   "🗑 Delete Reminder",
	}
	reminderDetailShowReminderCommandBtn := telebot.InlineButton{
		Unique: ReminderDetailShowReminderCommandBtn,
		Text:   "📄 Show Reminder Command",
	}
	closeCommandBtn := telebot.InlineButton{
		Unique: ReminderDetailCloseCommandBtn,
		Text:   "❌ Close Details",
	}

	return map[string]*telebot.InlineButton{
		ReminderDetailDeleteBtn:              &reminderDetailDeleteBtn,
		ReminderDetailShowReminderCommandBtn: &reminderDetailShowReminderCommandBtn,
		ReminderDetailCloseCommandBtn:        &closeCommandBtn,
	}
}

func HandleReminderDetailDeleteBtn(reminderDetailService Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		reminderID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		err = c.Respond(c.Callback())
		if err != nil {
			return err
		}

		err = reminderDetailService.DeleteReminder(int(c.ChatID()), reminderID)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Reminder %s has been deleted", strconv.Itoa(reminderID)))

		return err
	}
}

func HandleReminderShowReminderCommandBtn(reminderDetailService Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		reminderID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		reminder, err := reminderDetailService.GetReminder(int(c.ChatID()), reminderID)
		if err != nil {
			return err
		}

		err = c.Respond(c.Callback())
		if err != nil {
			return err
		}

		_, err = c.Send(reminder.Data.Command)

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
