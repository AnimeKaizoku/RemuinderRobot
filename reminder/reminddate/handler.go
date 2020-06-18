package reminddate

import (
	"fmt"

	"github.com/enrico5b1b4/telegram-bot/reminder"
)

func ReminderAddedSuccessMessage(message string, nextSchedule reminder.NextScheduleChatTime) string {
	return fmt.Sprintf("Reminder \"%s\" has been added for %s",
		message,
		nextSchedule.Time.In(nextSchedule.Location).Format("Mon, 02 Jan 2006 15:04 MST"),
	)
}
