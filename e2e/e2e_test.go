package e2e_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

const chatID = 123456

func TestE2E(t *testing.T) {
	checkSkip(t)

	dbFile := MustGetEnv("TEST_E2E_DB_FILE")
	allowedChats := []int{chatID}

	telebot, database, err := setup(dbFile, allowedChats)
	defer database.Close()
	require.NoError(t, err)

	// RemindAt
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me tomorrow at 20:45 MSG1_")
	require.Contains(t, telebot.OutboundSendMessages[0], `Reminder "MSG1_" has been added`)

	// RemindDayMonth
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on the 1st of december at 8:23 MSG2_")
	require.Contains(t, telebot.OutboundSendMessages[1], `Reminder "MSG2_" has been added`)

	// RemindEvery
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 2 minutes MSG3_")
	require.Contains(t, telebot.OutboundSendMessages[2], `Reminder "MSG3_" has been added`)

	// RemindEveryDayNumber
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of the month at 8:23 MSG4_")
	require.Contains(t, telebot.OutboundSendMessages[3], `Reminder "MSG4_" has been added`)

	// RemindEveryDayNumberMonth
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of december at 8:23 MSG5_")
	require.Contains(t, telebot.OutboundSendMessages[4], `Reminder "MSG5_" has been added`)

	// RemindEveryDayOfWeek
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every tuesday at 8:23 MSG6_")
	require.Contains(t, telebot.OutboundSendMessages[5], `Reminder "MSG6_" has been added`)

	// RemindIn
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me in 5 minutes MSG7_")
	require.Contains(t, telebot.OutboundSendMessages[6], `Reminder "MSG7_" has been added`)

	// RemindWhen
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me tomorrow MSG8_")
	require.Contains(t, telebot.OutboundSendMessages[7], `Reminder "MSG8_" has been added`)

	// RemindDayOfWeek
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on tuesday MSG9_")
	require.Contains(t, telebot.OutboundSendMessages[8], `Reminder "MSG9_" has been added`)

	// RemindEveryDay
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every day MSG10_")
	require.Contains(t, telebot.OutboundSendMessages[9], `Reminder "MSG10_" has been added`)

	telebot.SimulateIncomingMessageToChat(chatID, "/remindlist")
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG1_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG2_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG3_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG4_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG5_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG6_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG7_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG8_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG9_`)
	require.Contains(t, telebot.OutboundSendMessages[10], `MSG10_`)
}

func setup(dbFile string, allowedChats []int) (*fakes.TeleBot, *bolt.DB, error) {
	database, err := db.SetupDB(dbFile, allowedChats)
	if err != nil {
		return nil, nil, err
	}

	teleBot := fakes.NewTeleBot()
	botConfig := tbwrap.Config{
		AllowedChats: allowedChats,
		TBot:         teleBot,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		return nil, nil, err
	}

	appBot := bot.New(allowedChats, database, telegramBot)
	appBot.Start()

	return teleBot, database, nil
}

func MustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}

func checkSkip(t *testing.T) {
	testDBFile := os.Getenv("TEST_E2E_DB_FILE")
	if testDBFile == "" {
		t.Skip()
	}
}
