package healthcheck

import (
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
)

type HealthcheckService struct {
	b                *bot.Bot
	scheduler        *cron.Scheduler
	healthcheckStore *HealthcheckStore
}

func NewHealtcheckService(b *bot.Bot, scheduler *cron.Scheduler, healthcheckStore *HealthcheckStore) *HealthcheckService {
	return &HealthcheckService{
		b:                b,
		scheduler:        scheduler,
		healthcheckStore: healthcheckStore,
	}
}

func (s *HealthcheckService) CreateHealthcheck(schedule, message string, recipientID int) (int, error) {
	aliveHealthcheck := &Healthcheck{
		Job: cron.Job{
			Schedule: schedule,
			Type:     cron.HealthCheck,
			Status:   cron.Active,
		},
		Data: HealthcheckData{
			RecipientID: recipientID,
			Message:     message,
		},
	}
	aliveHealthcheckCronID, err := s.scheduler.Add(aliveHealthcheck.Job.Schedule, NewHealthcheckCronFunc(s.b, aliveHealthcheck))
	if err != nil {
		return 0, err
	}

	aliveHealthcheck.CronID = aliveHealthcheckCronID
	ID, err := s.healthcheckStore.CreateHealthcheck(aliveHealthcheck)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *HealthcheckService) LoadExistingSchedules() (int, error) {
	hcList, err := s.healthcheckStore.GetAllHealtchecks()
	if err != nil {
		return 0, err
	}

	for i := range hcList {
		aliveHealthcheckCronID, err := s.scheduler.Add(hcList[i].Job.Schedule, NewHealthcheckCronFunc(s.b, &hcList[i]))
		if err != nil {
			return 0, err
		}

		hcList[i].CronID = aliveHealthcheckCronID
		err = s.healthcheckStore.UpdateHealthcheck(&hcList[i])
		if err != nil {
			return 0, err
		}
	}

	return len(hcList), nil
}
