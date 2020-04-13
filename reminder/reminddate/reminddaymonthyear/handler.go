package reminddaymonthyear

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Who     string `regexpGroup:"who"`
	Day     int    `regexpGroup:"day"`
	Month   string `regexpGroup:"month"`
	Year    int    `regexpGroup:"year"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  *int   `regexpGroup:"minute"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind (?P<who>me|chat) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>january|february|march|april|may|june|july|august|september|october|november|december) (?P<year>\d{4}) ?(at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}))? (?P<message>.*)`

func HandleRemindDayMonthYear(service reminddate.Servicer) func(c tbwrap.Context) error {
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
		Day:    m.Day,
		Month:  date.ToNumericMonth(m.Month),
		Year:   m.Year,
		Hour:   9,
		Minute: 0,
	}

	if m.Hour != nil {
		dt.Hour = *m.Hour

		if m.Minute != nil {
			dt.Minute = *m.Minute
		}
	}

	return dt
}
