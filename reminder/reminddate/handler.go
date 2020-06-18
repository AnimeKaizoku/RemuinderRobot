package reminddate

import (
	"fmt"
)

func ReminderAddedSuccessMessage(message string, nextSchedule NextScheduleChatTime) string {
	return fmt.Sprintf("Reminder \"%s\" has been added for %s",
		message,
		nextSchedule.Time.In(nextSchedule.Location).Format("Mon, 02 Jan 2006 15:04 MST"),
	)
}
