package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/enrico5b1b4/telegram-bot/backup"
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/enrico5b1b4/telegram-bot/healthcheck"
	"github.com/enrico5b1b4/telegram-bot/parser"
	"github.com/enrico5b1b4/telegram-bot/reminder"
)

// /remind me on the 30 of november 2019 buy bread, milk and eggs without hour minutes
func main() {
	dbFile := MustGetEnv("DB_FILE") // TODO read from env
	telegramBotToken := MustGetEnv("TELEGRAM_BOT_TOKEN")
	allowedUsers := parseUsersAndGroups(MustGetEnv("ALLOWED_USERS"))
	allowedGroups := parseUsersAndGroups(MustGetEnv("ALLOWED_GROUPS"))

	database, err := db.SetupDB(dbFile, append(allowedUsers, allowedGroups...))
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	botConfig := bot.Config{
		Token:         telegramBotToken, // TODO read from env
		AllowedUsers:  allowedUsers,     // TODO read from env
		AllowedGroups: allowedGroups,    // TODO read from env
	}
	telegramBot, err := bot.NewBot(botConfig)

	//telegramBot.Send(&tb.User{ID: 1234567890}, "starting Bot")

	dbxConfig := dropbox.Config{
		Token:    "replaceME",
		LogLevel: dropbox.LogOff,
	}
	dbxFilesClient := files.New(dbxConfig)
	scheduler := cron.NewScheduler()
	reminderStore := reminder.NewReminderStore(database)
	reminderService := reminder.NewReminderService(telegramBot, scheduler, reminderStore)
	reminderParser := parser.NewParser(reminder.ReminderRegexes)
	reminderHandlers := reminder.NewReminderHandlers(telegramBot, reminderService, reminderParser, scheduler)
	healthcheckStore := healthcheck.NewHealthcheckStore(database)
	healthcheckService := healthcheck.NewHealtcheckService(telegramBot, scheduler, healthcheckStore)
	backupStore := backup.NewBackupStore(database)
	_ = backup.NewBackupService(telegramBot, scheduler, backupStore, dbxFilesClient)

	reminderStore.Debug()

	// check if DB exists and load schedules
	healthchecksLoaded, err := healthcheckService.LoadExistingSchedules()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d healthchecks", healthchecksLoaded)

	remindersLoaded, err := reminderService.LoadExistingSchedules()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d reminders", remindersLoaded)

	// Setup default healthchecks and backups
	//_, err = healthcheckService.CreateHealthcheck("@every 30s", "hello I'm alive", 1234567890)
	//if err != nil {
	//	panic(err)
	//}

	//reminder.RegisterHandlers(telegramBot, reminderService, reminderParser, scheduler)
	//reminder.RegisterHandlers(telegramBot, reminderHandlers)

	telegramBot.AddMultiRegExp([]string{
		`\/reminddetail (?P<reminderID>\d{1,5})`,
		`\/reminddetail_(?P<reminderID>\d{1,5})`,
	}, reminderHandlers.HandleNewRemindDetail)
	telegramBot.AddMultiRegExp([]string{
		`\/reminddelete (?P<reminderID>\d{1,5})`,
		`\/reminddelete_(?P<reminderID>\d{1,5})`,
	}, reminderHandlers.HandleNewRemindDelete)
	telegramBot.Add("/remindlist", reminderHandlers.HandleNewRemindList)
	telegramBot.AddRegExp(
		`\/remind (?P<who>me|group) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>october|november|december) (?P<year>\d{4}) (?P<message>.*)`,
		reminderHandlers.HandleNewRemindDayMonthYear,
	)
	telegramBot.AddRegExp(
		`\/remind (?P<who>me|group) on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? of (?P<month>october|november|december) (?P<year>\d{4}) at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}) (?P<message>.*)`,
		reminderHandlers.HandleNewRemindDayMonthYearHourMin,
	)
	telegramBot.AddButton(reminderHandlers.ReminderButtons["ReminderDetailDeleteBtn"], reminderHandlers.HandleNewReminderDetailDeleteBtn)
	telegramBot.AddButton(reminderHandlers.ReminderButtons["ReminderDetailShowReminderCommandBtn"], reminderHandlers.HandleNewReminderShowReminderCommandBtn)

	scheduler.Start()
	telegramBot.Start()
}

func parseUsersAndGroups(list string) []int {
	sepList := strings.Split(list, ",")
	intList := make([]int, len(sepList))
	var err error

	for i := range sepList {
		intList[i], err = strconv.Atoi(strings.TrimSpace(sepList[i]))
		if err != nil {
			log.Fatalln(err)
		}
	}

	return intList
}

func MustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}
