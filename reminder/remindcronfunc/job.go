package remindcronfunc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/telegram"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

// New creates a function which is called when a reminder is due
// Note: repeatable jobs can be of two kinds:
// - Reminders set as "remind me every 31 april at 13:52" will have a cron job like "52 13 31 April *"
//   These are occurrences which occur once a year
//   In this case RunOnlyOnce = false as we want to keep the schedule
// - Reminders set as "remind me every 3 minutes" will have a cron job set on a very specific date like "46 15 4 4 *".
//   These reminders are set with RunOnlyOnce = true as they should only run once.
//   They will have a RepeatSchedule which will reschedule the job for the following occurrence (e.g. in 3 minutes from now)
func New(s Servicer, b telegram.TBWrapBot, r *reminder.Reminder) func() {
	return func() {
		buttons := NewButtons()
		var inlineKeys [][]telebot.InlineButton
		var inlineButtons []telebot.InlineButton

		snoozeBtn := *buttons[SnoozeBtn]
		snoozeBtn.Data = strconv.Itoa(r.ID)
		inlineButtons = append(
			inlineButtons,
			snoozeBtn,
		)

		// if repeatable job add button to complete it
		if !r.Job.RunOnlyOnce || (r.Job.RunOnlyOnce && r.Job.RepeatSchedule != nil) {
			completeBtn := *buttons[CompleteBtn]
			completeBtn.Data = strconv.Itoa(r.ID)
			inlineButtons = append(inlineButtons, completeBtn)
		}
		inlineKeys = append(inlineKeys, inlineButtons)

		messageWithIcon := fmt.Sprintf("🗓 %s", r.Data.Message)
		_, err := b.Send(&tb.Chat{ID: int64(r.Data.RecipientID)}, messageWithIcon, &telebot.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
		if err != nil {
			log.Printf("NewReminderCronFunc err: %q", err)
			return
		}

		if !r.Job.RunOnlyOnce {
			return
		}

		if r.Job.RepeatSchedule != nil {
			updateErr := s.UpdateReminderWithRepeatSchedule(r)
			if updateErr != nil {
				log.Printf("NewReminderCronFunc UpdateReminderWithRepeatSchedule err: %q", updateErr)
				return
			}
			return
		}

		err = s.Complete(r)
		if err != nil {
			log.Printf("NewReminderCronFunc complete err: %q", err)
			return
		}
	}
}
