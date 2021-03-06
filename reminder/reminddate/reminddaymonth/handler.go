package reminddaymonth

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Day     int    `regexpGroup:"day"`
	Month   string `regexpGroup:"month"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind me on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? ?(of (?P<month>january|february|march|april|may|june|july|august|september|october|november|december))? ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindDayMonth(service reminddate.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime := mapMessageToReminderDateTime(message)
		nextSchedule, err := service.AddReminderOnDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(reminddate.ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageToReminderDateTime(m *Message) reminder.DateTime {
	dt := reminder.DateTime{
		DayOfMonth: m.Day,
		Month:      date.ToNumericMonth(m.Month),
		Hour:       9,
		Minute:     0,
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		dt.Hour = hour
		dt.Minute = minute
	}

	return dt
}
