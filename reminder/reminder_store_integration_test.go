package reminder_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/stretchr/testify/assert"
)

func TestReminderStore_CreateReminder(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()
	reminderStore := reminder.NewStore(database)

	t.Run("success", func(t *testing.T) {
		r := &reminder.Reminder{
			Job: cron.Job{
				ChatID:   chatID,
				Schedule: "schedule",
			},
		}
		id, err := reminderStore.CreateReminder(r)
		assert.NoError(t, err)

		checkReminder, err := reminderStore.GetReminder(chatID, id)
		assert.NoError(t, err)
		assert.Equal(t, &reminder.Reminder{
			Job: cron.Job{ID: 1, ChatID: chatID, Schedule: "schedule"},
		}, checkReminder)
	})
}

func TestReminderStore_UpdateReminder(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()

	reminderStore := reminder.NewStore(database)
	r := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID,
			Schedule: "schedule",
		},
	}
	id, err := reminderStore.CreateReminder(r)
	assert.NoError(t, err)

	existingReminder, err := reminderStore.GetReminder(chatID, id)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		// nolint: goconst
		existingReminder.Schedule = "updated schedule"
		err = reminderStore.UpdateReminder(existingReminder)
		assert.NoError(t, err)

		updatedReminder, err := reminderStore.GetReminder(chatID, id)
		assert.NoError(t, err)
		assert.Equal(t, "updated schedule", updatedReminder.Schedule)
	})
}

func TestReminderStore_DeleteReminder(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()

	reminderStore := reminder.NewStore(database)
	r := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID,
			Schedule: "schedule",
		},
	}
	id, err := reminderStore.CreateReminder(r)
	assert.NoError(t, err)

	existingReminder, err := reminderStore.GetReminder(chatID, id)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		existingReminder.Schedule = "updated schedule"
		err = reminderStore.DeleteReminder(chatID, id)
		assert.NoError(t, err)

		_, err := reminderStore.GetReminder(chatID, id)
		assert.Error(t, err)
	})
}

func TestReminderStore_GetReminder(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()

	reminderStore := reminder.NewStore(database)
	r := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID,
			Schedule: "schedule",
		},
	}
	id, err := reminderStore.CreateReminder(r)
	assert.NoError(t, err)

	existingReminder, err := reminderStore.GetReminder(chatID, id)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		foundReminder, err := reminderStore.GetReminder(chatID, id)
		assert.NoError(t, err)
		assert.Equal(t, existingReminder, foundReminder)
	})
}

func TestReminderStore_GetAllRemindersByChat(t *testing.T) {
	checkSkip(t)

	chatID1 := rand.Intn(100000)
	chatID2 := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID1, chatID2})
	assert.NoError(t, err)
	defer database.Close()

	reminderStore := reminder.NewStore(database)
	r1 := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID1,
			Schedule: "schedule",
		},
	}
	_, err = reminderStore.CreateReminder(r1)
	assert.NoError(t, err)

	r2 := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID2,
			Schedule: "schedule",
		},
	}
	_, err = reminderStore.CreateReminder(r2)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		reminders, err := reminderStore.GetAllRemindersByChat()
		assert.NoError(t, err)
		assert.Len(t, reminders[chatID1], 1)
		assert.Len(t, reminders[chatID2], 1)
	})
}

func TestReminderStore_GetAllRemindersByChatID(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()

	reminderStore := reminder.NewStore(database)
	r := &reminder.Reminder{
		Job: cron.Job{
			ChatID:   chatID,
			Schedule: "schedule",
		},
	}
	id, err := reminderStore.CreateReminder(r)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		reminders, err := reminderStore.GetAllRemindersByChatID(chatID)
		assert.NoError(t, err)
		assert.Equal(t, id, reminders[0].ID)
	})
}

func checkSkip(t *testing.T) {
	testDBFile := os.Getenv("TEST_DB_FILE")
	if testDBFile == "" {
		t.Skip()
	}
}

func testDBFile() string {
	return fmt.Sprintf("../%s", os.Getenv("TEST_DB_FILE"))
}
