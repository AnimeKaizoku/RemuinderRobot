package chatpreference_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/stretchr/testify/assert"
)

func TestChatPreferenceStore_CreateChatPreference(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()
	chatPreferenceStore := chatpreference.NewStore(database)

	t.Run("success", func(t *testing.T) {
		cp := &chatpreference.ChatPreference{
			ChatID: chatID,
		}
		err := chatPreferenceStore.CreateChatPreference(cp)
		assert.NoError(t, err)

		checkChatPreference, err := chatPreferenceStore.GetChatPreference(chatID)
		assert.NoError(t, err)
		assert.Equal(t, &chatpreference.ChatPreference{
			ChatID: chatID,
		}, checkChatPreference)
	})
}

func TestChatPreferenceStore_GetChatPreference(t *testing.T) {
	checkSkip(t)

	chatID := rand.Intn(100000)
	database, err := db.SetupDB(testDBFile(), []int{chatID})
	assert.NoError(t, err)
	defer database.Close()

	chatPreferenceStore := chatpreference.NewStore(database)
	cp := &chatpreference.ChatPreference{
		ChatID: chatID,
	}
	err = chatPreferenceStore.CreateChatPreference(cp)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		foundChatPreference, err := chatPreferenceStore.GetChatPreference(chatID)
		assert.NoError(t, err)
		assert.Equal(t, cp, foundChatPreference)
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
