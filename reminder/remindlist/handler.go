package remindlist

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/enrico5b1b4/tbwrap"
	"gopkg.in/tucnak/telebot.v2"
)

const HandlePattern = "/remindlist"

func HandleRemindList(reminderListService Servicer, buttons map[string]*telebot.InlineButton) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		remindersByStatus, err := reminderListService.GetRemindersByChatID(int(c.ChatID()))
		if err != nil {
			return err
		}

		// index 0 = cron.Active, index 1 = cron.Inactive, index 2 = cron.Completed
		if len(remindersByStatus) == 0 ||
			(len(remindersByStatus[0].Entries) == 0 &&
				len(remindersByStatus[1].Entries) == 0 &&
				len(remindersByStatus[2].Entries) == 0) {
			_, err = c.Send("You have no reminders.")

			return err
		}

		var remindListInlineKeys [][]telebot.InlineButton
		if len(buttons) > 0 {
			closeCommandBtn := *buttons[ReminderListCloseCommandBtn]
			closeCommandBtn.Data = strconv.Itoa(c.Message().ID)
			remindListInlineKeys = append(remindListInlineKeys, []telebot.InlineButton{closeCommandBtn})

			// only show button if there are completed reminders
			if len(remindersByStatus[2].Entries) > 0 {
				reminderListRemoveCompletedRemindersBtn := *buttons[ReminderListRemoveCompletedRemindersBtn]
				remindListInlineKeys = append(
					remindListInlineKeys,
					[]telebot.InlineButton{reminderListRemoveCompletedRemindersBtn},
				)
			}
		}

		t := template.Must(template.New("text").Parse(text))
		var buf bytes.Buffer
		if execErr := t.Execute(&buf, remindersByStatus); execErr != nil {
			return execErr
		}

		_, err = c.Send(buf.String(), &telebot.ReplyMarkup{
			InlineKeyboard: remindListInlineKeys,
		})

		return err
	}
}

// nolint:lll
const text = `
{{ range . }}{{if .Entries}}*{{.Status}}*{{$previousTimeKey:=""}}
{{ range .Entries }}{{ if .Time }}{{$currentTimeKey:=.Time.Format "2 Jan 2006"}}{{if ne $currentTimeKey $previousTimeKey}}*{{$currentTimeKey}}*{{printf "\n"}}{{$previousTimeKey = $currentTimeKey}}{{end}}{{ end }}{{ range .Entries }}{{$nextSchedule:=""}}{{if .NextSchedule}}{{$nextSchedule = .NextSchedule.Format "15:04"}}{{end}}{{if ( and (.RunOnlyOnce) (not .RepeatSchedule))}}{{printf "- _%s_ %s [[/r_%d]]" $nextSchedule .Data.Message .ID}}{{ else }}{{printf "- 🔁 _%s_ %s [[/r_%d]]" $nextSchedule .Data.Message .ID}}{{ end }}{{printf "\n"}}{{ end }}
{{ end }}{{ end }}
{{ end }}
`
