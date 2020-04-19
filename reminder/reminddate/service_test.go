package reminddate_test

import (
	"testing"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	chatpreferenceMocks "github.com/enrico5b1b4/telegram-bot/chatpreference/mocks"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	reminderMocks "github.com/enrico5b1b4/telegram-bot/reminder/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	schedulerMocks "github.com/enrico5b1b4/telegram-bot/reminder/scheduler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type Mocks struct {
	ReminderStore       *reminderMocks.MockStorer
	Scheduler           *schedulerMocks.MockScheduler
	ChatPreferenceStore *chatpreferenceMocks.MockStorer
}

const (
	message = "message"
	command = "command"
)

func TestService_AddReminderOnDateTime(t *testing.T) {
	t.Run("success with day of month", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnDateTime(chatID, command, reminder.DateTime{
			DayOfMonth: 1,
			Month:      date.ToNumericMonth(time.April.String()),
			Hour:       13,
			Minute:     52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
	t.Run("success with day of week", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 * * 1",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 * * 1",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnDateTime(chatID, command, reminder.DateTime{
			DayOfWeek: "1",
			Hour:      13,
			Minute:    52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
}

func TestService_AddReminderOnWordDateTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: "Europe/London",
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 1 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderOnWordDateTime(chatID, command, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   13,
			Minute: 52,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
}

func TestService_AddRepeatableReminderOnDateTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "52 13 31 April *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: false,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "52 13 31 April *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: false,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddRepeatableReminderOnDateTime(chatID, command, &reminder.RepeatableDateTime{
			DayOfMonth: "31",
			Month:      time.April.String(),
			Hour:       "13",
			Minute:     "52",
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
}

func TestService_AddReminderIn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: "Europe/London",
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderIn(chatID, command, reminder.AmountDateTime{
			Minutes: 1,
			Hours:   2,
			Days:    3,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
}

func TestService_AddReminderEvery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		chatID := 1
		cronID := 2
		reminderID := 3
		stubNextScheduleTime := timeNow()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mocks := createMocks(mockCtrl)
		mocks.ChatPreferenceStore.EXPECT().GetChatPreference(chatID).Return(&chatpreference.ChatPreference{
			ChatID:   chatID,
			TimeZone: "Europe/London",
		}, nil)
		mocks.Scheduler.EXPECT().AddReminder(&reminder.Reminder{
			Job: cron.Job{
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				RepeatSchedule: &cron.JobRepeatSchedule{
					Minutes: 1,
					Hours:   2,
					Days:    3,
				},
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(cronID, nil)
		mocks.ReminderStore.EXPECT().CreateReminder(&reminder.Reminder{
			Job: cron.Job{
				CronID:      cronID,
				ChatID:      chatID,
				Schedule:    "46 15 4 4 *",
				Type:        cron.Reminder,
				Status:      cron.Active,
				RunOnlyOnce: true,
				RepeatSchedule: &cron.JobRepeatSchedule{
					Minutes: 1,
					Hours:   2,
					Days:    3,
				},
			},
			Data: reminder.Data{
				RecipientID: chatID,
				Message:     message,
				Command:     command,
			},
		}).Return(reminderID, nil)
		mocks.Scheduler.EXPECT().GetNextScheduleTime(cronID).Return(stubNextScheduleTime, nil)

		service := reminddate.NewService(mocks.Scheduler, mocks.ReminderStore, mocks.ChatPreferenceStore, timeNow)
		nextScheduleTime, err := service.AddReminderEvery(chatID, command, reminder.AmountDateTime{
			Minutes: 1,
			Hours:   2,
			Days:    3,
		}, message)
		assert.NoError(t, err)
		assert.Equal(t, stubNextScheduleTime, nextScheduleTime)
	})
}

func createMocks(mockCtrl *gomock.Controller) Mocks {
	return Mocks{
		ReminderStore:       reminderMocks.NewMockStorer(mockCtrl),
		Scheduler:           schedulerMocks.NewMockScheduler(mockCtrl),
		ChatPreferenceStore: chatpreferenceMocks.NewMockStorer(mockCtrl),
	}
}

func timeNow() time.Time {
	timeLoc, _ := time.LoadLocation("Europe/London")
	return time.Date(2020, time.April, 1, 13, 45, 0, 0, timeLoc)
}
