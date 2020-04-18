package scheduler

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/scheduler/mocks/${GOFILE} -package=mocks

import (
	"fmt"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/telegram"
)

type Scheduler interface {
	AddReminder(r *reminder.Reminder) (int, error)
	GetNextScheduleTime(cronID int) (time.Time, error)
}

type ReminderScheduler struct {
	reminderStore           reminder.Storer
	reminderCronFuncService remindcronfunc.Servicer
	scheduler               cron.Scheduler
	bot                     telegram.TBWrapBot
	chatPreferenceStore     chatpreference.Storer
}

func NewReminderScheduler(
	bot telegram.TBWrapBot,
	reminderCronFuncService remindcronfunc.Servicer,
	reminderStore reminder.Storer,
	scheduler cron.Scheduler,
	chatPreferenceStore chatpreference.Storer,
) *ReminderScheduler {
	return &ReminderScheduler{
		bot:                     bot,
		reminderStore:           reminderStore,
		reminderCronFuncService: reminderCronFuncService,
		scheduler:               scheduler,
		chatPreferenceStore:     chatPreferenceStore,
	}
}

func (s *ReminderScheduler) AddReminder(rem *reminder.Reminder) (int, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(rem.Job.ChatID)
	if err != nil {
		return 0, err
	}

	schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rem.Job.Schedule)
	reminderCronID, err := s.scheduler.Add(schedule, remindcronfunc.New(s.reminderCronFuncService, s.bot, rem))
	if err != nil {
		return 0, err
	}

	return reminderCronID, nil
}

func (s *ReminderScheduler) GetNextScheduleTime(cronID int) (time.Time, error) {
	cronEntry := s.scheduler.GetEntryByID(cronID)

	return cronEntry.Next, nil
}
