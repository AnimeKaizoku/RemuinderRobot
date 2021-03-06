package remindeverydayofweek

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Day     string `regexpGroup:"day"`
	When    string `regexpGroup:"when"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind me every (?P<day>(((M|m)(on)|(T|t)(ues)|(W|w)(ednes)|(T|t)(hurs)|(F|f)(ri)|(S|s)(atur)|(S|s)(un))(day))) ?(?P<when>morning|afternoon|evening|night)? ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindEveryDayOfWeek(service reminddate.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		repeatDateTime := mapMessageToReminderDateTime(message)
		nextSchedule, err := service.AddRepeatableReminderOnDateTime(int(c.ChatID()), c.Text(), &repeatDateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(reminddate.ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))

		return err
	}
}

func mapMessageToReminderDateTime(m *Message) reminder.RepeatableDateTime {
	rdt := reminder.RepeatableDateTime{
		DayOfWeek: strconv.Itoa(date.ToNumericDayOfWeek(m.Day)),
		Month:     "*",
		Hour:      "9",
		Minute:    "0",
	}

	switch m.When {
	case "morning":
		rdt.Hour = "9"
		rdt.Minute = "0"

	case "afternoon":
		rdt.Hour = "15"
		rdt.Minute = "0"

	case "evening", "night":
		rdt.Hour = "20"
		rdt.Minute = "0"

	default:
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		rdt.Hour = strconv.Itoa(hour)
		rdt.Minute = strconv.Itoa(minute)
	}

	return rdt
}
