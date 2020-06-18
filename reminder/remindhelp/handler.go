package remindhelp

import (
	"github.com/enrico5b1b4/tbwrap"
)

const HandlePattern = "/remindhelp"

func HandleRemindHelp() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		_, err := c.Send(text)

		return err
	}
}

const text = `
*Available commands*

_list reminders_
/remindlist

_get details of a reminder_
[/r_ID]

_delete a reminder_
[/reminddelete_ID]

_set a reminder_
/remind me at 21:00 update weekly report
/remind me on the 1st of december at 8:23 update weekly report
/remind me on the 1st of december update weekly report
/remind me in 4 minutes update weekly report
/remind me in 5 hours update weekly report
/remind me in 3 hours, 4 minutes update weekly report
/remind me in 5 days, 3 hours, 4 minutes update weekly report
/remind me tonight/this evening/tomorrow/tomorrow morning update weekly report
/remind me today/tomorrow at 21:00 update weekly report
/remind me on tuesday update weekly report

_set a recurring reminder_
/remind me every 1st of the month at 8:23 update weekly report
/remind me every 1st of the month update weekly report
/remind me every 1st of december at 8:23 update weekly report
/remind me every 1st of december update weekly report
/remind me every 2 minutes update weekly report
/remind me every 3 hours, 4 minutes update weekly report
/remind me every 5 days, 3 hours, 4 minutes update weekly report
/remind me every day at 8pm update weekly report

_set timezone for chat reminders_
/gettimezone
/settimezone Europe/London
`
