package remindwhen

import (
	"errors"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Who     string `regexpGroup:"who"`
	When    string `regexpGroup:"when"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind (?P<who>me|chat) (?P<when>this afternoon|this evening|tonight|tomorrow morning|tomorrow afternoon|tomorrow evening|tomorrow) (?P<message>.*)`

func HandleRemindWhen(service reminddate.Servicer) func(c tbwrap.Context) error {
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
	switch m.When {
	case "this afternoon":
		return reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   15,
			Minute: 0,
		}, nil
	case "this evening", "tonight":
		return reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   20,
			Minute: 0,
		}, nil
	case "tomorrow", "tomorrow morning":
		return reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   9,
			Minute: 0,
		}, nil
	case "tomorrow afternoon":
		return reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   15,
			Minute: 0,
		}, nil
	case "tomorrow evening":
		return reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   20,
			Minute: 0,
		}, nil
	}

	return reminder.WordDateTime{}, errors.New("no match")
}
