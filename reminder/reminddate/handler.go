package reminddate

import (
	"fmt"
	"time"
)

func ReminderAddedSuccessMessage(message string, nextSchedule time.Time) string {
	return fmt.Sprintf("Reminder \"%s\" has been added for %s",
		message,
		nextSchedule.Format("Mon, 02 Jan 2006 15:04 MST"),
	)
}
