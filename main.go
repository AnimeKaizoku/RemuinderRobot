package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/db"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/loader"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindat"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonthyear"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonthyearhourmin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindevery"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumber"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumberhourmin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumbermonth"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydaynumbermonthhourmin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindwhen"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindwhenhourmin"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddelete"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddetail"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindhelp"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindlist"
)

// nolint:funlen
func main() {
	dbFile := MustGetEnv("DB_FILE")
	telegramBotToken := MustGetEnv("TELEGRAM_BOT_TOKEN")
	allowedChats := parseAllowedChats(MustGetEnv("ALLOWED_CHATS"))

	database, err := db.SetupDB(dbFile, allowedChats)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	botConfig := tbwrap.Config{
		Token:        telegramBotToken,
		AllowedChats: allowedChats,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Println(err)
		return
	}

	scheduler := cron.NewScheduler()
	reminderStore := reminder.NewStore(database)
	chatPreferenceStore := chatpreference.NewStore(database)
	chatPreferenceService := chatpreference.NewService(chatPreferenceStore)
	remindCronFuncService := remindcronfunc.NewService(telegramBot, scheduler, reminderStore, chatPreferenceStore)
	remindListService := remindlist.NewService(reminderStore, scheduler, chatPreferenceStore)
	remindDeleteService := reminddelete.NewService(reminderStore, scheduler)
	remindDateService := reminddate.NewService(telegramBot, remindCronFuncService, reminderStore, scheduler, chatPreferenceStore)
	remindDetailService := reminddetail.NewService(reminderStore, scheduler, chatPreferenceStore)
	reminderLoader := loader.NewService(telegramBot, scheduler, reminderStore, chatPreferenceStore, remindCronFuncService)
	remindDetailButtons := reminddetail.NewButtons()
	remindListButtons := remindlist.NewButtons()

	chatPreferenceService.CreateDefaultChatPreferences(allowedChats)

	// check if DB exists and load schedules
	remindersLoaded, err := reminderLoader.LoadExistingSchedules()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d reminders", remindersLoaded)

	telegramBot.Handle(remindlist.HandlePattern, remindlist.HandleRemindList(remindListService, remindListButtons))
	telegramBot.Handle(remindhelp.HandlePattern, remindhelp.HandleRemindHelp())
	telegramBot.HandleMultiRegExp(reminddetail.HandlePattern, reminddetail.HandleRemindDetail(remindDetailService, reminddetail.NewButtons()))
	telegramBot.HandleMultiRegExp(reminddelete.HandlePattern, reminddelete.HandleRemindDelete(remindDeleteService))
	telegramBot.HandleRegExp(
		reminddaymonthyearhourmin.HandlePattern,
		reminddaymonthyearhourmin.HandleRemindDayMonthYearHourMinute(remindDateService),
	)
	telegramBot.HandleRegExp(
		reminddaymonthyear.HandlePattern,
		reminddaymonthyear.HandleRemindDayMonthYear(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumberhourmin.HandlePattern,
		remindeverydaynumberhourmin.HandleRemindEveryDayNumberHourMin(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumber.HandlePattern,
		remindeverydaynumber.HandleRemindEveryDayNumber(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindeverydaynumbermonthhourmin.HandlePattern,
		remindeverydaynumbermonthhourmin.HandleRemindEveryDayNumberMonthHourMin(remindDateService),
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
		remindwhenhourmin.HandlePattern,
		remindwhenhourmin.HandleRemindWhenHourMin(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindwhen.HandlePattern,
		remindwhen.HandleRemindWhen(remindDateService),
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

	scheduler.Start()
	telegramBot.Start()
}

func parseAllowedChats(list string) []int {
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
