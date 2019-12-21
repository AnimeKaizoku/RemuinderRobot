package reminddelete

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/reminddelete/mocks/${GOFILE} -package=mocks

import (
	"errors"

	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
)

type Servicer interface {
	DeleteReminder(chatID, ID int) error
}

type Service struct {
	reminderStore reminder.Storer
	scheduler     cron.Scheduler
}

func NewService(reminderStore reminder.Storer, scheduler cron.Scheduler) *Service {
	return &Service{
		reminderStore: reminderStore,
		scheduler:     scheduler,
	}
}

func (s *Service) DeleteReminder(chatID, id int) error {
	r, err := s.reminderStore.GetReminder(chatID, id)
	if err != nil {
		return err
	}

	if chatID != r.ChatID {
		return errors.New("unauthorized to delete reminder")
	}

	s.scheduler.Remove(r.CronID)

	return s.reminderStore.DeleteReminder(r.ChatID, id)
}
