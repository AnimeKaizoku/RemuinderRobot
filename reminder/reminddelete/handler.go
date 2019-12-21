package reminddelete

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type Message struct {
	ReminderID int `regexpGroup:"reminderID"`
}

var HandlePattern = []string{
	`\/reminddelete (?P<reminderID>\d{1,5})`,
	`\/reminddelete_(?P<reminderID>\d{1,5})`,
}

func HandleRemindDelete(reminderDeleteService Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		err := reminderDeleteService.DeleteReminder(int(c.ChatID()), message.ReminderID)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Reminder %d has been deleted", message.ReminderID))

		return err
	}
}
