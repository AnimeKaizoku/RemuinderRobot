package reminddate

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/reminddate/mocks/${GOFILE} -package=mocks

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/scheduler"
)

type Servicer interface {
	AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (time.Time, error)
	AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (time.Time, error)
	AddRepeatableReminderOnDateTime(chatID int, command string, dateTime *reminder.RepeatableDateTime, message string) (time.Time, error)
	AddReminderIn(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error)
	AddReminderEvery(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error)
}

type Service struct {
	reminderStore       reminder.Storer
	reminderScheduler   scheduler.Scheduler
	chatPreferenceStore chatpreference.Storer
	timeNow             func() time.Time
}

func NewService(
	reminderScheduler scheduler.Scheduler,
	reminderStore reminder.Storer,
	chatPreferenceStore chatpreference.Storer,
	timeNow func() time.Time,
) *Service {
	return &Service{
		reminderScheduler:   reminderScheduler,
		reminderStore:       reminderStore,
		chatPreferenceStore: chatPreferenceStore,
		timeNow:             timeNow,
	}
}

func (s *Service) AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (time.Time, error) {
	var schedule string
	if dateTime.DayOfWeek != "" {
		schedule = buildScheduleForDateTime(&dateTime)
	} else {
		chatLocalTime, err := s.getChatLocalDateTime(chatID, dateTime.Month, dateTime.DayOfMonth, dateTime.Hour, dateTime.Minute)
		if err != nil {
			return s.timeNow(), err
		}

		err = s.validateInFuture(chatLocalTime.In(time.UTC))
		if err != nil {
			return s.timeNow(), err
		}

		schedule = buildScheduleForDateTime(&dateTime)
	}

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

	return s.ScheduleAndAddReminder(newReminder)
}

func (s *Service) AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (time.Time, error) {
	chatLocalTime, err := s.convertWordDateTimeToChatLocalDateTime(chatID, dateTime)
	if err != nil {
		return s.timeNow(), err
	}

	err = s.validateInFuture(chatLocalTime.In(time.UTC))
	if err != nil {
		return s.timeNow(), err
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

	return s.ScheduleAndAddReminder(newReminder)
}

func (s *Service) convertWordDateTimeToChatLocalDateTime(chatID int, dateTime reminder.WordDateTime) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return s.timeNow(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return s.timeNow(), err
	}

	// default to today
	timeNowChatLocalTime := s.timeNow().In(loc)
	if dateTime.When == reminder.Tomorrow {
		hours := 24
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
	chatID int, command string, repeatDateTime *reminder.RepeatableDateTime, message string,
) (time.Time, error) {
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    buildScheduleForRepeatableDateTime(repeatDateTime),
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

	return s.ScheduleAndAddReminder(newReminder)
}

func (s *Service) AddReminderIn(
	chatID int, command string, amountDateTime reminder.AmountDateTime, message string,
) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return s.timeNow(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return s.timeNow(), err
	}

	addedTime := s.timeNow().In(loc).Add(
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

	return s.ScheduleAndAddReminder(newReminder)
}

func (s *Service) AddReminderEvery(
	chatID int, command string, amountDateTime reminder.AmountDateTime, message string,
) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return s.timeNow(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return s.timeNow(), err
	}

	addedTime := s.timeNow().In(loc).Add(
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

	return s.ScheduleAndAddReminder(newReminder)
}

func (s *Service) ScheduleAndAddReminder(rem *reminder.Reminder) (time.Time, error) {
	cronID, err := s.reminderScheduler.AddReminder(rem)
	if err != nil {
		return s.timeNow(), err
	}

	rem.CronID = cronID
	_, err = s.reminderStore.CreateReminder(rem)
	if err != nil {
		return s.timeNow(), err
	}

	return s.reminderScheduler.GetNextScheduleTime(cronID)
}

func (s *Service) getChatLocalDateTime(chatID, month, day, hour, minute int) (time.Time, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return s.timeNow(), err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return s.timeNow(), err
	}

	return time.Date(s.timeNow().Year(), time.Month(month), day, hour, minute, 0, 0, loc), nil
}

var minutesInFutureBeforeInvalid = 2 * time.Minute

func (s *Service) validateInFuture(t time.Time) error {
	currentTimeUTC := s.timeNow().Add(minutesInFutureBeforeInvalid).In(time.UTC)
	if t.Before(currentTimeUTC) {
		return errors.New("error: time must be at least 3 minutes in the future")
	}

	return nil
}

func buildScheduleForRepeatableDateTime(repeatDateTime *reminder.RepeatableDateTime) string {
	return fmt.Sprintf("%s %s %s %s %s",
		asteriskIfEmpty(repeatDateTime.Minute),
		asteriskIfEmpty(repeatDateTime.Hour),
		asteriskIfEmpty(repeatDateTime.DayOfMonth),
		asteriskIfEmpty(repeatDateTime.Month),
		asteriskIfEmpty(repeatDateTime.DayOfWeek),
	)
}

func buildScheduleForDateTime(repeatDateTime *reminder.DateTime) string {
	return fmt.Sprintf("%s %s %s %s %s",
		asteriskIfZero(repeatDateTime.Minute),
		asteriskIfZero(repeatDateTime.Hour),
		asteriskIfZero(repeatDateTime.DayOfMonth),
		asteriskIfZero(repeatDateTime.Month),
		asteriskIfEmpty(repeatDateTime.DayOfWeek),
	)
}

func asteriskIfZero(val int) string {
	if val == 0 {
		return "*"
	}

	return strconv.Itoa(val)
}

func asteriskIfEmpty(val string) string {
	if val == "" {
		return "*"
	}

	return val
}
