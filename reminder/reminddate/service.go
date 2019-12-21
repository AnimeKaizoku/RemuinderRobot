package reminddate

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/reminddate/mocks/${GOFILE} -package=mocks

import (
	"errors"
	"fmt"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/telegram"
)

type Servicer interface {
	AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (time.Time, error)
	AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (time.Time, error)
	AddRepeatableReminderOnDateTime(chatID int, command string, dateTime reminder.RepeatableDateTime, message string) (time.Time, error)
	AddReminderIn(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error)
	AddReminderEvery(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error)
}

type Service struct {
	reminderStore           reminder.Storer
	reminderCronFuncService remindcronfunc.Servicer
	scheduler               cron.Scheduler
	chatPreferenceStore     chatpreference.Storer
	b                       telegram.Bot
}

func NewService(
	b telegram.Bot,
	reminderCronFuncService remindcronfunc.Servicer,
	reminderStore reminder.Storer,
	scheduler cron.Scheduler,
	chatPreferenceStore chatpreference.Storer,
) *Service {
	return &Service{
		b:                       b,
		reminderStore:           reminderStore,
		reminderCronFuncService: reminderCronFuncService,
		scheduler:               scheduler,
		chatPreferenceStore:     chatPreferenceStore,
	}
}

func (s *Service) AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (time.Time, error) {
	chatLocalTime, err := s.getChatLocalDateTime(chatID, dateTime.Year, dateTime.Month, dateTime.Day, dateTime.Hour, dateTime.Minute)
	if err != nil {
		return time.Now(), err
	}

	err = validateInFuture(chatLocalTime.In(time.UTC))
	if err != nil {
		return time.Now(), err
	}

	schedule := fmt.Sprintf("%d %d %d %d *", chatLocalTime.Minute(), chatLocalTime.Hour(), chatLocalTime.Day(), chatLocalTime.Month())
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    schedule,
			Type:        cron.Reminder,
			Status:      cron.Active,
			RunOnlyOnce: true,
		},
		Data: reminder.Data{
			RecipientID: chatID,
			Message:     message,
			Command:     command,
		},
	}

	id, err := s.addReminder(newReminder)
	if err != nil {
		return time.Now(), err
	}

	nextScheduleTime, err := s.getNextScheduleTime(chatID, id)
	if err != nil {
		return time.Now(), err
	}

	return nextScheduleTime, nil
}

func (s *Service) AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (time.Time, error) {
	chatLocalTime, err := s.convertWordDateTimeToChatLocalDateTime(chatID, dateTime)
	if err != nil {
		return time.Now(), err
	}

	err = validateInFuture(chatLocalTime.In(time.UTC))
	if err != nil {
		return time.Now(), err
	}

	schedule := fmt.Sprintf("%d %d %d %d *", chatLocalTime.Minute(), chatLocalTime.Hour(), chatLocalTime.Day(), chatLocalTime.Month())
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    schedule,
			Type:        cron.Reminder,
			Status:      cron.Active,
			RunOnlyOnce: true,
		},
		Data: reminder.Data{
			RecipientID: chatID,
			Message:     message,
			Command:     command,
		},
	}

	id, err := s.addReminder(newReminder)
	if err != nil {
		return time.Now(), err
	}

	nextScheduleTime, err := s.getNextScheduleTime(chatID, id)
	if err != nil {
		return time.Now(), err
	}

	return nextScheduleTime, nil
}

func (s *Service) convertWordDateTimeToChatLocalDateTime(chatID int, dateTime reminder.WordDateTime) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return time.Now(), err
	}

	// default to today
	timeNowChatLocalTime := time.Now().In(loc)
	hours := 24
	if dateTime.When == reminder.Tomorrow {
		timeNowChatLocalTime = timeNowChatLocalTime.Add(time.Duration(hours) * time.Hour)
	}

	return time.Date(
		timeNowChatLocalTime.Year(),
		timeNowChatLocalTime.Month(),
		timeNowChatLocalTime.Day(),
		dateTime.Hour,
		dateTime.Minute,
		0,
		0,
		loc,
	), nil
}

func (s *Service) AddRepeatableReminderOnDateTime(
	chatID int, command string, repeatDateTime reminder.RepeatableDateTime, message string,
) (time.Time, error) {
	schedule := fmt.Sprintf("%s %s %s %s *",
		repeatDateTime.Minute,
		repeatDateTime.Hour,
		repeatDateTime.Day,
		repeatDateTime.Month,
	)
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    schedule,
			Type:        cron.Reminder,
			Status:      cron.Active,
			RunOnlyOnce: false,
		},
		Data: reminder.Data{
			RecipientID: chatID,
			Message:     message,
			Command:     command,
		},
	}

	id, err := s.addReminder(newReminder)
	if err != nil {
		return time.Now(), err
	}

	nextScheduleTime, err := s.getNextScheduleTime(chatID, id)
	if err != nil {
		return time.Now(), err
	}

	return nextScheduleTime, nil
}

func (s *Service) AddReminderIn(
	chatID int, command string, amountDateTime reminder.AmountDateTime, message string,
) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return time.Now(), err
	}

	addedTime := time.Now().In(loc).Add(
		time.Duration(amountDateTime.Days)*24*time.Hour +
			time.Duration(amountDateTime.Hours)*time.Hour +
			time.Duration(amountDateTime.Minutes)*time.Minute,
	)

	schedule := fmt.Sprintf("%d %d %d %d *", addedTime.Minute(), addedTime.Hour(), addedTime.Day(), addedTime.Month())
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    schedule,
			Type:        cron.Reminder,
			Status:      cron.Active,
			RunOnlyOnce: true,
		},
		Data: reminder.Data{
			RecipientID: chatID,
			Message:     message,
			Command:     command,
		},
	}

	id, err := s.addReminder(newReminder)
	if err != nil {
		return time.Now(), err
	}

	nextScheduleTime, err := s.getNextScheduleTime(chatID, id)
	if err != nil {
		return time.Now(), err
	}

	return nextScheduleTime, nil
}

func (s *Service) AddReminderEvery(
	chatID int, command string, amountDateTime reminder.AmountDateTime, message string,
) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return time.Now(), err
	}

	addedTime := time.Now().In(loc).Add(
		time.Duration(amountDateTime.Days)*24*time.Hour +
			time.Duration(amountDateTime.Hours)*time.Hour +
			time.Duration(amountDateTime.Minutes)*time.Minute,
	)

	schedule := fmt.Sprintf("%d %d %d %d *", addedTime.Minute(), addedTime.Hour(), addedTime.Day(), addedTime.Month())
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    schedule,
			Type:        cron.Reminder,
			Status:      cron.Active,
			RunOnlyOnce: true,
			RepeatSchedule: &cron.JobRepeatSchedule{
				Hours:   amountDateTime.Hours,
				Days:    amountDateTime.Days,
				Minutes: amountDateTime.Minutes,
			},
		},
		Data: reminder.Data{
			RecipientID: chatID,
			Message:     message,
			Command:     command,
		},
	}

	id, err := s.addReminder(newReminder)
	if err != nil {
		return time.Now(), err
	}

	nextScheduleTime, err := s.getNextScheduleTime(chatID, id)
	if err != nil {
		return time.Now(), err
	}

	return nextScheduleTime, nil
}

func (s *Service) getChatLocalDateTime(chatID, year, month, day, hour, minute int) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc), nil
}

func (s *Service) addReminder(rem *reminder.Reminder) (int, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(rem.Job.ChatID)
	if err != nil {
		return 0, err
	}

	schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rem.Job.Schedule)
	reminderCronID, err := s.scheduler.Add(schedule, remindcronfunc.New(s.reminderCronFuncService, s.b, rem))
	if err != nil {
		return 0, err
	}

	rem.CronID = reminderCronID
	ID, err := s.reminderStore.CreateReminder(rem)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *Service) getNextScheduleTime(chatID, reminderID int) (time.Time, error) {
	r, err := s.reminderStore.GetReminder(chatID, reminderID)
	if err != nil {
		return time.Now(), err
	}

	cronEntry := s.scheduler.GetEntryByID(r.CronID)

	return cronEntry.Next, nil
}

var minutesInFutureBeforeInvalid = 2 * time.Minute

func validateInFuture(t time.Time) error {
	currentTimeUTC := time.Now().Add(minutesInFutureBeforeInvalid).In(time.UTC)
	if t.Before(currentTimeUTC) {
		return errors.New("error: time must be at least 3 minutes in the future")
	}

	return nil
}
