package e2e_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindat"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonth"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddayofweek"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindevery"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumber"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumbermonth"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydayofweek"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindwhen"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddelete"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddetail"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindhelp"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindlist"
	"github.com/enrico5b1b4/telegram-bot/reminder/scheduler"
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

	telebot.SimulateIncomingMessageToChat(chatID, "/remind me at 20:45 RemindAt")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on the 1st of december at 8:23 RemindDayMonth")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 2 minutes RemindEvery")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of the month at 8:23 RemindEveryDayNumber")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every 1st of december at 8:23 RemindEveryDayNumberMonth")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me every tuesday at 8:23 RemindEveryDayOfWeek")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me in 5 minutes RemindIn")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me tonight RemindWhen")
	telebot.SimulateIncomingMessageToChat(chatID, "/remind me on tuesday RemindDayOfWeek")

	require.Contains(t, telebot.OutboundSendMessages[0], `Reminder "RemindAt" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[1], `Reminder "RemindDayMonth" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[2], `Reminder "RemindEvery" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[3], `Reminder "RemindEveryDayNumber" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[4], `Reminder "RemindEveryDayNumberMonth" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[5], `Reminder "RemindEveryDayOfWeek" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[6], `Reminder "RemindIn" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[7], `Reminder "RemindWhen" has been added`)
	require.Contains(t, telebot.OutboundSendMessages[8], `Reminder "RemindDayOfWeek" has been added`)
}

// TODO refactor main to make this setup easier
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

	cronScheduler := cron.NewScheduler()
	reminderStore := reminder.NewStore(database)
	chatPreferenceStore := chatpreference.NewStore(database)
	chatPreferenceService := chatpreference.NewService(chatPreferenceStore)
	remindCronFuncService := remindcronfunc.NewService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore)
	remindListService := remindlist.NewService(reminderStore, cronScheduler, chatPreferenceStore)
	remindDeleteService := reminddelete.NewService(reminderStore, cronScheduler)
	reminderScheduler := scheduler.NewReminderScheduler(telegramBot, remindCronFuncService, reminderStore, cronScheduler, chatPreferenceStore)
	remindDateService := reminddate.NewService(reminderScheduler, reminderStore, chatPreferenceStore, reminder.RealTimeNow)
	remindDetailService := reminddetail.NewService(reminderStore, cronScheduler, chatPreferenceStore)
	remindDetailButtons := reminddetail.NewButtons()
	remindListButtons := remindlist.NewButtons()

	chatPreferenceService.CreateDefaultChatPreferences(allowedChats)

	telegramBot.Handle(remindlist.HandlePattern, remindlist.HandleRemindList(remindListService, remindListButtons))
	telegramBot.Handle(remindhelp.HandlePattern, remindhelp.HandleRemindHelp())
	telegramBot.HandleMultiRegExp(reminddetail.HandlePattern, reminddetail.HandleRemindDetail(remindDetailService, reminddetail.NewButtons()))
	telegramBot.HandleMultiRegExp(reminddelete.HandlePattern, reminddelete.HandleRemindDelete(remindDeleteService))
	telegramBot.HandleRegExp(
		reminddaymonth.HandlePattern,
		reminddaymonth.HandleRemindDayMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		reminddayofweek.HandlePattern,
		reminddayofweek.HandleRemindDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumber.HandlePattern,
		remindeverydaynumber.HandleRemindEveryDayNumber(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumbermonth.HandlePattern,
		remindeverydaynumbermonth.HandleRemindEveryDayNumberMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindin.HandlePattern3,
		remindin.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindin.HandlePattern2,
		remindin.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindin.HandlePattern1,
		remindin.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindevery.HandlePattern3,
		remindevery.HandleRemindEvery(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindevery.HandlePattern2,
		remindevery.HandleRemindEvery(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindevery.HandlePattern1,
		remindevery.HandleRemindEvery(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindwhen.HandlePattern,
		remindwhen.HandleRemindWhen(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydayofweek.HandlePattern,
		remindeverydayofweek.HandleRemindEveryDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindat.HandlePattern,
		remindat.HandleRemindAt(remindDateService),
	)

	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailCloseCommandBtn],
		reminddetail.HandleCloseBtn(remindDetailButtons),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailDeleteBtn],
		reminddetail.HandleReminderDetailDeleteBtn(remindDetailService, remindDetailButtons),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailShowReminderCommandBtn],
		reminddetail.HandleReminderShowReminderCommandBtn(remindDetailService, remindDetailButtons),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListRemoveCompletedRemindersBtn],
		remindlist.HandleReminderListRemoveCompletedRemindersBtn(remindListService, remindListButtons),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListCloseCommandBtn],
		remindlist.HandleCloseBtn(remindListButtons),
	)

	cronScheduler.Start()
	telegramBot.Start()

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
