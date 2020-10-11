package remindcronfunc

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	Snooze10MinuteBtn          = "Snooze10MinuteBtn"
	Snooze20MinuteBtn          = "Snooze20MinuteBtn"
	Snooze30MinuteBtn          = "Snooze30MinuteBtn"
	Snooze1HourBtn             = "Snooze1HourBtn"
	SnoozeThisAfternoonBtn     = "SnoozeThisAfternoonBtn"
	SnoozeThisEveningBtn       = "SnoozeThisEveningBtn"
	SnoozeTomorrowMorningBtn   = "SnoozeTomorrowMorningBtn"
	SnoozeTomorrowAfternoonBtn = "SnoozeTomorrowAfternoonBtn"
	SnoozeTomorrowEveningBtn   = "SnoozeTomorrowEveningBtn"
	SnoozeBtn                  = "SnoozeBtn"
	SnoozeCloseBtn             = "SnoozeCloseBtn"
	CompleteBtn                = "CompleteBtn"
)

func NewButtons() map[string]*telebot.InlineButton {
	snooze10MinuteBtn := telebot.InlineButton{
		Unique: Snooze10MinuteBtn,
		Text:   "⏰ 10m",
	}
	snooze20MinuteBtn := telebot.InlineButton{
		Unique: Snooze20MinuteBtn,
		Text:   "⏰ 20m",
	}
	snooze30MinuteBtn := telebot.InlineButton{
		Unique: Snooze30MinuteBtn,
		Text:   "⏰ 30m",
	}
	snooze1HourBtn := telebot.InlineButton{
		Unique: Snooze1HourBtn,
		Text:   "⏰ 1h",
	}
	snoozeThisAfternoonBtn := telebot.InlineButton{
		Unique: SnoozeThisAfternoonBtn,
		Text:   "⏰ This Afternoon",
	}
	snoozeThisEveningBtn := telebot.InlineButton{
		Unique: SnoozeThisEveningBtn,
		Text:   "⏰ This Evening",
	}
	snoozeTomorrowMorningBtn := telebot.InlineButton{
		Unique: SnoozeTomorrowMorningBtn,
		Text:   "⏰ Tomorrow Morning",
	}
	snoozeTomorrowAfternoonBtn := telebot.InlineButton{
		Unique: SnoozeTomorrowAfternoonBtn,
		Text:   "⏰ Tomorrow Afternoon",
	}
	snoozeTomorrowEveningBtn := telebot.InlineButton{
		Unique: SnoozeTomorrowEveningBtn,
		Text:   "⏰ Tomorrow Evening",
	}
	snoozeBtn := telebot.InlineButton{
		Unique: SnoozeBtn,
		Text:   "⏰ Snooze",
	}
	snoozeCloseBtn := telebot.InlineButton{
		Unique: SnoozeCloseBtn,
		Text:   "❌ Close",
	}
	completeBtn := telebot.InlineButton{
		Unique: CompleteBtn,
		Text:   "✅ Done",
	}

	return map[string]*telebot.InlineButton{
		Snooze10MinuteBtn:          &snooze10MinuteBtn,
		Snooze20MinuteBtn:          &snooze20MinuteBtn,
		Snooze30MinuteBtn:          &snooze30MinuteBtn,
		Snooze1HourBtn:             &snooze1HourBtn,
		SnoozeThisAfternoonBtn:     &snoozeThisAfternoonBtn,
		SnoozeThisEveningBtn:       &snoozeThisEveningBtn,
		SnoozeTomorrowMorningBtn:   &snoozeTomorrowMorningBtn,
		SnoozeTomorrowAfternoonBtn: &snoozeTomorrowAfternoonBtn,
		SnoozeTomorrowEveningBtn:   &snoozeTomorrowEveningBtn,
		CompleteBtn:                &completeBtn,
		SnoozeBtn:                  &snoozeBtn,
		SnoozeCloseBtn:             &snoozeCloseBtn,
	}
}

type RemindDateServicer interface {
	AddReminderIn(
		chatID int,
		command string,
		amountDateTime reminder.AmountDateTime,
		message string,
	) (reminder.NextScheduleChatTime, error)
	AddReminderOnWordDateTime(
		chatID int,
		command string,
		dateTime reminder.WordDateTime,
		message string,
	) (reminder.NextScheduleChatTime, error)
}

// nolint:dupl
func HandleReminderSnoozeAmountDateTimeBtn(
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
		if err != nil {
			return err
		}

		return c.Delete(c.ChatID(), c.Message().ID)
	}
}

// nolint:dupl
func HandleReminderSnoozeWordDateTimeBtn(
	service RemindDateServicer,
	store reminder.Storer,
	wordDateTime reminder.WordDateTime,
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

		nextSchedule, err := service.AddReminderOnWordDateTime(chatID, rem.Data.Command, wordDateTime, rem.Data.Message)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Reminder \"%s\" has been rescheduled for %s",
			rem.Data.Message,
			nextSchedule.Time.In(nextSchedule.Location).Format("Mon, 02 Jan 2006 15:04 MST"),
		))
		if err != nil {
			return err
		}

		return c.Delete(c.ChatID(), c.Message().ID)
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

func HandleReminderSnoozeBtn(store reminder.Storer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		err := c.Respond(c.Callback())
		if err != nil {
			return err
		}

		reminderID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		rem, err := store.GetReminder(int(c.ChatID()), reminderID)
		if err != nil {
			return err
		}

		buttons := NewButtons()

		snooze10MinuteBtn := *buttons[Snooze10MinuteBtn]
		snooze10MinuteBtn.Data = strconv.Itoa(reminderID)
		snooze20MinuteBtn := *buttons[Snooze20MinuteBtn]
		snooze20MinuteBtn.Data = strconv.Itoa(reminderID)
		snooze30MinuteBtn := *buttons[Snooze30MinuteBtn]
		snooze30MinuteBtn.Data = strconv.Itoa(reminderID)
		snooze1HourBtn := *buttons[Snooze1HourBtn]
		snooze1HourBtn.Data = strconv.Itoa(reminderID)
		snoozeThisAfternoonBtn := *buttons[SnoozeThisAfternoonBtn]
		snoozeThisAfternoonBtn.Data = strconv.Itoa(reminderID)
		snoozeThisEveningBtn := *buttons[SnoozeThisEveningBtn]
		snoozeThisEveningBtn.Data = strconv.Itoa(reminderID)
		snoozeTomorrowMorningBtn := *buttons[SnoozeTomorrowMorningBtn]
		snoozeTomorrowMorningBtn.Data = strconv.Itoa(reminderID)
		snoozeTomorrowAfternoonBtn := *buttons[SnoozeTomorrowAfternoonBtn]
		snoozeTomorrowAfternoonBtn.Data = strconv.Itoa(reminderID)
		snoozeTomorrowEveningBtn := *buttons[SnoozeTomorrowEveningBtn]
		snoozeTomorrowEveningBtn.Data = strconv.Itoa(reminderID)
		snoozeBtn := *buttons[SnoozeBtn]
		snoozeBtn.Data = strconv.Itoa(reminderID)
		snoozeCloseBtn := *buttons[SnoozeCloseBtn]

		inlineKeys := [][]telebot.InlineButton{
			{
				snooze10MinuteBtn,
				snooze30MinuteBtn,
				snooze1HourBtn,
			},
			{
				snoozeThisAfternoonBtn,
			},
			{
				snoozeThisEveningBtn,
			},
			{
				snoozeTomorrowMorningBtn,
			},
			{
				snoozeTomorrowAfternoonBtn,
			},
			{
				snoozeTomorrowEveningBtn,
			},
			{
				snoozeCloseBtn,
			},
		}

		messageWithIcon := fmt.Sprintf("🗓 %s", rem.Data.Message)
		_, err = c.Send(messageWithIcon, &telebot.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
		if err != nil {
			return err
		}

		return err
	}
}

func HandleReminderSnoozeCloseBtn() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		return c.Delete(c.ChatID(), c.Message().ID)
	}
}
