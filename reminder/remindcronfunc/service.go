package remindcronfunc

// nolint: lll
//go:generate mockgen -destination=./mocks/mock_Servicer.go -package=mocks github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc Servicer

import (
	"fmt"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/telegram"
)

type Servicer interface {
	Complete(r *reminder.Reminder) error
	UpdateReminderWithRepeatSchedule(rem *reminder.Reminder) error
}

type Service struct {
	b                   telegram.TBWrapBot
	scheduler           cron.Scheduler
	reminderStore       reminder.Storer
	chatPreferenceStore chatpreference.Storer
}

func NewService(
	b telegram.TBWrapBot,
	scheduler cron.Scheduler,
	reminderStore reminder.Storer,
	chatPreferenceStore chatpreference.Storer,
) *Service {
	return &Service{
		b:                   b,
		scheduler:           scheduler,
		reminderStore:       reminderStore,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *Service) Complete(r *reminder.Reminder) error {
	r.Status = cron.Completed
	timeNow := time.Now().In(time.UTC)
	r.CompletedAt = &timeNow

	err := s.reminderStore.UpdateReminder(r)
	if err != nil {
		return err
	}

	s.scheduler.Remove(r.CronID)
	return nil
}

func (s *Service) UpdateReminderWithRepeatSchedule(rem *reminder.Reminder) error {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(rem.Job.ChatID)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return err
	}

	addedTime := time.Now().In(loc).Add(
		time.Duration(rem.RepeatSchedule.Days)*24*time.Hour +
			time.Duration(rem.RepeatSchedule.Hours)*time.Hour +
			time.Duration(rem.RepeatSchedule.Minutes)*time.Minute,
	)

	schedule := fmt.Sprintf("%d %d %d %d *",
		addedTime.Minute(),
		addedTime.Hour(),
		addedTime.Day(),
		addedTime.Month(),
	)
	rem.Job.Schedule = schedule

	// remove previous cron job before scheduling new one
	s.scheduler.Remove(rem.CronID)

	scheduleWithTZ := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, schedule)
	reminderCronID, err := s.scheduler.Add(scheduleWithTZ, New(s, s.b, rem))
	if err != nil {
		return err
	}

	rem.CronID = reminderCronID
	err = s.reminderStore.UpdateReminder(rem)
	if err != nil {
		return err
	}

	return nil
}
