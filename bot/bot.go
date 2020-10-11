package bot

import (
	"log"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/chatpreference/timezone/gettimezone"
	"github.com/enrico5b1b4/telegram-bot/chatpreference/timezone/settimezone"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/loader"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindat"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonth"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddayofweek"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindevery"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeveryday"
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
	"github.com/enrico5b1b4/telegram-bot/telegram"
	"go.etcd.io/bbolt"
)

type Bot struct {
	cronScheduler cron.Scheduler
	telegramBot   telegram.TBWrapBot
}

// nolint:funlen
func New(
	allowedChats []int,
	database *bbolt.DB,
	telegramBot telegram.TBWrapBot,
) *Bot {
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
	reminderLoader := loader.NewService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore, remindCronFuncService)
	setTimeZoneService := settimezone.NewService(chatPreferenceStore, reminderLoader)
	remindDetailButtons := reminddetail.NewButtons()
	remindListButtons := remindlist.NewButtons()
	reminderCompleteButtons := remindcronfunc.NewButtons()

	chatPreferenceService.CreateDefaultChatPreferences(allowedChats)

	// check if DB exists and load schedules
	remindersLoaded, err := reminderLoader.LoadSchedulesFromDB()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d reminders", remindersLoaded)

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
		remindin.HandlePattern,
		remindin.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindevery.HandlePattern,
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
		remindeveryday.HandlePattern,
		remindeveryday.HandleRemindEveryDay(remindDateService),
	)
	telegramBot.HandleRegExp(
		remindat.HandlePattern,
		remindat.HandleRemindAt(remindDateService),
	)
	telegramBot.Handle(gettimezone.HandlePattern, gettimezone.HandleGetTimezone(chatPreferenceStore))
	telegramBot.HandleRegExp(settimezone.HandlePattern, settimezone.HandleSetTimezone(setTimeZoneService))

	// buttons
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailCloseCommandBtn],
		reminddetail.HandleCloseBtn(),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailDeleteBtn],
		reminddetail.HandleReminderDetailDeleteBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindDetailButtons[reminddetail.ReminderDetailShowReminderCommandBtn],
		reminddetail.HandleReminderShowReminderCommandBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListRemoveCompletedRemindersBtn],
		remindlist.HandleReminderListRemoveCompletedRemindersBtn(remindListService),
	)
	telegramBot.HandleButton(
		remindListButtons[remindlist.ReminderListCloseCommandBtn],
		remindlist.HandleCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.Snooze10MinuteBtn],
		remindcronfunc.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 10}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.Snooze20MinuteBtn],
		remindcronfunc.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 20}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.Snooze30MinuteBtn],
		remindcronfunc.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 30}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.Snooze1HourBtn],
		remindcronfunc.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 60}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeThisAfternoonBtn],
		remindcronfunc.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeThisEveningBtn],
		remindcronfunc.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeTomorrowMorningBtn],
		remindcronfunc.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   9,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeTomorrowAfternoonBtn],
		remindcronfunc.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeTomorrowEveningBtn],
		remindcronfunc.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeBtn],
		remindcronfunc.HandleReminderSnoozeBtn(reminderStore),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.SnoozeCloseBtn],
		remindcronfunc.HandleReminderSnoozeCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[remindcronfunc.CompleteBtn],
		remindcronfunc.HandleReminderCompleteBtn(remindCronFuncService, reminderStore),
	)

	return &Bot{
		cronScheduler: cronScheduler,
		telegramBot:   telegramBot,
	}
}

func (b *Bot) Start() {
	b.cronScheduler.Start()
	b.telegramBot.Start()
}
