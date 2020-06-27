package remindat

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Hour    int    `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind me at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)? (?P<message>.*)`

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
	hour, minute := date.ConvertTo24H(m.Hour, m.Minute, m.AMPM)

	return reminder.WordDateTime{
		When:   reminder.Today,
		Hour:   hour,
		Minute: minute,
	}
}
