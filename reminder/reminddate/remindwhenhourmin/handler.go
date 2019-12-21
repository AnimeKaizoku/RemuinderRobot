package remindwhenhourmin

import (
	"errors"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Who     string `regexpGroup:"who"`
	When    string `regexpGroup:"when"`
	Hour    int    `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind (?P<who>me|chat) (?P<when>today|this afternoon|this evening|tonight|tomorrow morning|tomorrow afternoon|tomorrow evening|tomorrow) at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}) (?P<message>.*)`

func HandleRemindWhenHourMin(service reminddate.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime, err := mapMessageToReminderWordDateTime(message)
		if err != nil {
			return err
		}

		nextSchedule, err := service.AddReminderOnWordDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(reminddate.ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageToReminderWordDateTime(m *Message) (reminder.WordDateTime, error) {
	wordDateTime := reminder.WordDateTime{}

	switch m.When {
	case "today", "this afternoon", "this evening", "tonight":
		wordDateTime = reminder.WordDateTime{
			When: reminder.Today,
		}
	case "tomorrow", "tomorrow morning", "tomorrow afternoon", "tomorrow evening":
		wordDateTime = reminder.WordDateTime{
			When: reminder.Tomorrow,
		}
	default:
		return reminder.WordDateTime{}, errors.New("no match")
	}

	wordDateTime.Hour = m.Hour
	wordDateTime.Minute = m.Minute

	return wordDateTime, nil
}
