package settimezone

import (
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/reminder/loader"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/chatpreference/timezone/settimezone/mocks/${GOFILE} -package=mocks

type Servicer interface {
	SetTimeZone(chatID int, timezone string) error
}

type Service struct {
	reminderLoader      loader.Servicer
	chatPreferenceStore chatpreference.Storer
}

func NewService(chatPreferenceStore chatpreference.Storer, reminderLoader loader.Servicer) *Service {
	return &Service{
		reminderLoader:      reminderLoader,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *Service) SetTimeZone(chatID int, timezone string) error {
	if err := validateTimeZone(timezone); err != nil {
		return err
	}

	if err := s.chatPreferenceStore.UpsertChatPreference(&chatpreference.ChatPreference{
		ChatID:   chatID,
		TimeZone: timezone,
	}); err != nil {
		return err
	}

	_, err := s.reminderLoader.ReloadSchedulesForChat(chatID)
	if err != nil {
		return err
	}

	return nil
}

// validateTimeZone validates input timezone
func validateTimeZone(tz string) error {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}

	return nil
}
