package reminddetail

import (
	"bytes"
	"html/template"
	"strconv"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"gopkg.in/tucnak/telebot.v2"
)

type Message struct {
	ReminderID int `regexpGroup:"reminderID"`
}

type ReminderDetail struct {
	reminder.Reminder
	NextSchedule *time.Time
}

var HandlePattern = []string{
	`\/reminddetail (?P<reminderID>\d{1,5})`,
	`\/reminddetail_(?P<reminderID>\d{1,5})`,
	`\/r (?P<reminderID>\d{1,5})`,
	`\/r_(?P<reminderID>\d{1,5})`,
}

func HandleRemindDetail(reminderDetailService Servicer, buttons map[string]*telebot.InlineButton) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		reminderDetail, err := reminderDetailService.GetReminder(int(c.ChatID()), message.ReminderID)
		if err != nil {
			return err
		}

		var remindDetailInlineKeys [][]telebot.InlineButton
		if len(buttons) > 0 {
			closeCommandBtn := *buttons[ReminderDetailCloseCommandBtn]
			closeCommandBtn.Data = strconv.Itoa(c.Message().ID)
			remindDetailInlineKeys = append(remindDetailInlineKeys, []telebot.InlineButton{closeCommandBtn})

			reminderDetailShowReminderCommandBtn := *buttons[ReminderDetailShowReminderCommandBtn]
			reminderDetailShowReminderCommandBtn.Data = strconv.Itoa(message.ReminderID)
			remindDetailInlineKeys = append(remindDetailInlineKeys, []telebot.InlineButton{reminderDetailShowReminderCommandBtn})

			reminderDetailDeleteBtn := *buttons[ReminderDetailDeleteBtn]
			reminderDetailDeleteBtn.Data = strconv.Itoa(message.ReminderID)
			remindDetailInlineKeys = append(remindDetailInlineKeys, []telebot.InlineButton{reminderDetailDeleteBtn})
		}

		t := template.Must(template.New("text").Parse(text))
		var buf bytes.Buffer
		if execErr := t.Execute(&buf, reminderDetail); execErr != nil {
			return execErr
		}

		_, err = c.Send(buf.String(), &telebot.ReplyMarkup{
			InlineKeyboard: remindDetailInlineKeys,
		})

		return err
	}
}

// nolint:lll
const text = `
*Id*: {{.ID}}
*Status*: {{.Status}}
*Message*: {{.Data.Message}}
*Command*: {{.Data.Command}}
{{if .NextSchedule}}*Next Schedule*: {{.NextSchedule.Format "Mon, 02 Jan 2006 15:04 MST"}}{{end}}{{if .CompletedAt}}*Completed At*: {{.CompletedAt.Format "Mon, 02 Jan 2006 15:04 MST"}}{{end}}
`
