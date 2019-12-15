package reminder

import (
	"fmt"
	"log"
	"regexp"

	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
	tb "gopkg.in/tucnak/telebot.v2"
)

// "command": "/remind me every monday buy chocolate",
// 			"status": "active",
// 			"message": "buy chocolate",
// 			"dates": [{}]

type Reminder struct {
	cron.Job
	Data ReminderData `json:"data"`
}

type ReminderData struct {
	RecipientID int    `json:"recipient_id"`
	Command     string `json:"command"`
	Message     string `json:"message"`
}

func NewReminderCronFunc(s *ReminderService, b *bot.Bot, r *Reminder) func() {
	return func() {
		fmt.Println("NewReminderCronFunc")
		_, err := b.Send(&tb.Chat{ID: int64(r.Data.RecipientID)}, r.Data.Message)
		if err != nil {
			log.Printf("NewReminderCronFunc err: %q", err)
			return
		}

		if r.Job.RunOnlyOnce {
			err := s.Complete(r)
			if err != nil {
				log.Printf("NewReminderCronFunc complete err: %q", err)
				return
			}
		}
	}
}

var ReminderRegexes = map[string]*regexp.Regexp{
	// "regexEveryDayName":        regexp.MustCompile(`(?P<who>me|group) every (?P<day>monday|tuesday) (?P<message>.*)`),
	// "regexEveryDayNameNMonths": regexp.MustCompile(`(?P<who>me|group) every (?P<day>monday|tuesday) for (?P<nMonths>\d{2}) months (?P<message>.*)`),
	"regexOnDayMonthYearHourMin": regexp.MustCompile(`\/remind (?P<who>me|group) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>october|november|december) (?P<year>\d{4}) at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}) (?P<message>.*)`),
	"regexOnDayMonthYear":        regexp.MustCompile(`\/remind (?P<who>me|group) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>october|november|december) (?P<year>\d{4}) (?P<message>.*)`),
	// "regexOnDayMonth":          regexp.MustCompile(`(?P<who>me|group) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>october|november) (?P<message>.*)`),
	// "regexEveryDayNumber":      regexp.MustCompile(`(?P<who>me|group) every (?P<day>\d{1,2})(?:(st|nd|rd|th))? of the month (?P<message>.*)`),
	"regexDeleteReminder": regexp.MustCompile(`\/reminddelete (?P<reminderID>\d{1,2})`),
	"regexDetailReminder": regexp.MustCompile(`\/reminddetail (?P<reminderID>\d{1,2})`),
}

var ReminderOnTextRegexes = map[string]*regexp.Regexp{
	"regexReminderDeleteByID": regexp.MustCompile(`\/reminddelete_(?P<reminderID>\d{1,5})`),
	"regexReminderDetailByID": regexp.MustCompile(`\/reminddetail_(?P<reminderID>\d{1,5})`),
}
