package remindat

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Who     string `regexpGroup:"who"`
	Hour    int    `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind (?P<who>me|chat) at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}) (?P<message>.*)`

func HandleRemindAt(service reminddate.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime := mapMessageToReminderWordDateTime(message)
		nextSchedule, err := service.AddReminderOnWordDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(reminddate.ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageToReminderWordDateTime(m *Message) reminder.WordDateTime {
	return reminder.WordDateTime{
		When:   reminder.Today,
		Hour:   m.Hour,
		Minute: m.Minute,
	}
}
