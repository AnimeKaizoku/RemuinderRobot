package loader

import (
	"fmt"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/telegram"
)

type Service struct {
	b                   telegram.Bot
	scheduler           cron.Scheduler
	reminderStore       reminder.Storer
	reminderJobService  remindcronfunc.Servicer
	chatPreferenceStore chatpreference.Storer
}

func NewService(
	b telegram.Bot,
	scheduler cron.Scheduler,
	reminderStore reminder.Storer,
	chatPreferenceStore chatpreference.Storer,
	reminderJobService remindcronfunc.Servicer,
) *Service {
	return &Service{
		b:                   b,
		scheduler:           scheduler,
		reminderStore:       reminderStore,
		chatPreferenceStore: chatPreferenceStore,
		reminderJobService:  reminderJobService,
	}
}

func (s *Service) LoadExistingSchedules() (int, error) {
	remindersAdded := 0
	rmdrListByChat, err := s.reminderStore.GetAllRemindersByChat()
	if err != nil {
		return 0, err
	}

	for chatID := range rmdrListByChat {
		for i := range rmdrListByChat[chatID] {
			if rmdrListByChat[chatID][i].Status != cron.Active {
				continue
			}

			chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
			if err != nil {
				return 0, err
			}

			schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rmdrListByChat[chatID][i].Job.Schedule)
			reminderID, err := s.scheduler.Add(
				schedule,
				remindcronfunc.New(s.reminderJobService, s.b, &rmdrListByChat[chatID][i]),
			)
			if err != nil {
				return 0, err
			}

			rmdrListByChat[chatID][i].CronID = reminderID
			err = s.reminderStore.UpdateReminder(&rmdrListByChat[chatID][i])
			if err != nil {
				return 0, err
			}

			remindersAdded++
		}
	}

	return len(rmdrListByChat), nil
}
