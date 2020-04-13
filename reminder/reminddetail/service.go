package reminddetail

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/reminddetail/mocks/${GOFILE} -package=mocks

import (
	"errors"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
)

type Servicer interface {
	GetReminder(chatID, reminderID int) (*ReminderDetail, error)
	DeleteReminder(chatID, ID int) error
}

type Service struct {
	reminderStore       reminder.Storer
	scheduler           cron.Scheduler
	chatPreferenceStore chatpreference.Storer
}

func NewService(reminderStore reminder.Storer, scheduler cron.Scheduler, chatPreferenceStore chatpreference.Storer) *Service {
	return &Service{
		reminderStore:       reminderStore,
		scheduler:           scheduler,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *Service) GetReminder(chatID, reminderID int) (*ReminderDetail, error) {
	rem, err := s.reminderStore.GetReminder(chatID, reminderID)
	if err != nil {
		return nil, err
	}

	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return nil, err
	}

	reminderDetail := &ReminderDetail{Reminder: *rem}
	if rem.Status == cron.Active {
		cronEntry := s.scheduler.GetEntryByID(rem.CronID)
		nextScheduleInChatTimezone := cronEntry.Next.In(loc)
		reminderDetail.NextSchedule = &nextScheduleInChatTimezone
	}
	if rem.Status == cron.Completed && rem.CompletedAt != nil {
		completedAtChatTimezone := rem.CompletedAt.In(loc)
		reminderDetail.CompletedAt = &completedAtChatTimezone
	}

	return reminderDetail, nil
}

func (s *Service) DeleteReminder(chatID, id int) error {
	rem, err := s.reminderStore.GetReminder(chatID, id)
	if err != nil {
		return err
	}

	if chatID != rem.ChatID {
		return errors.New("unauthorised to delete reminder")
	}

	s.scheduler.Remove(rem.CronID)

	return s.reminderStore.DeleteReminder(rem.ChatID, id)
}
