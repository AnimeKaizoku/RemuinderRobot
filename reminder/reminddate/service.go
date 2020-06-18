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
	AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (NextScheduleChatTime, error)
	AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (NextScheduleChatTime, error)
	AddRepeatableReminderOnDateTime(chatID int,
		command string,
		dateTime *reminder.RepeatableDateTime,
		message string,
	) (NextScheduleChatTime, error)
	AddReminderIn(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (NextScheduleChatTime, error)
	AddReminderEvery(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (NextScheduleChatTime, error)
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

func (s *Service) AddReminderOnDateTime(chatID int,
	command string,
	dateTime reminder.DateTime,
	message string,
) (NextScheduleChatTime, error) {
	newReminder := &reminder.Reminder{
		Job: cron.Job{
			ChatID:      chatID,
			Schedule:    buildScheduleForDateTime(&dateTime),
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

func (s *Service) AddReminderOnWordDateTime(chatID int,
	command string,
	dateTime reminder.WordDateTime,
	message string,
) (NextScheduleChatTime, error) {
	chatLocalTime, err := s.convertWordDateTimeToChatLocalDateTime(chatID, dateTime)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	err = s.validateInFuture(chatLocalTime.In(time.UTC))
	if err != nil {
		return NextScheduleChatTime{}, err
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
) (NextScheduleChatTime, error) {
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
) (NextScheduleChatTime, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return NextScheduleChatTime{}, err
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
) (NextScheduleChatTime, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	loc, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return NextScheduleChatTime{}, err
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

type NextScheduleChatTime struct {
	Time     time.Time
	Location *time.Location
}

func (s *Service) ScheduleAndAddReminder(rem *reminder.Reminder) (NextScheduleChatTime, error) {
	cronID, err := s.reminderScheduler.AddReminder(rem)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	rem.CronID = cronID
	_, err = s.reminderStore.CreateReminder(rem)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	nextScheduleTime, err := s.reminderScheduler.GetNextScheduleTime(cronID)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	cp, err := s.chatPreferenceStore.GetChatPreference(rem.ChatID)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	loc, err := time.LoadLocation(cp.TimeZone)
	if err != nil {
		return NextScheduleChatTime{}, err
	}

	return NextScheduleChatTime{Time: nextScheduleTime, Location: loc}, nil
}

func (s *Service) validateInFuture(t time.Time) error {
	minutesInFutureBeforeInvalid := 2 * time.Minute
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
