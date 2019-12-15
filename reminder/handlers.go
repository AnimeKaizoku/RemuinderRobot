package reminder

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/parser"
	"gopkg.in/tucnak/telebot.v2"
)

type ReminderHandler struct {
	b               *bot.Bot
	scheduler       *cron.Scheduler
	reminderService *ReminderService
	reminderParser  *parser.Parser
	ReminderButtons map[string]*telebot.InlineButton
}

func NewReminderHandlers(telegramBot *bot.Bot, reminderService *ReminderService, reminderParser *parser.Parser, scheduler *cron.Scheduler) *ReminderHandler {
	return &ReminderHandler{
		b:               telegramBot,
		scheduler:       scheduler,
		reminderParser:  reminderParser,
		reminderService: reminderService,
		ReminderButtons: NewReminderButtons(),
	}
}

func NewReminderButtons() map[string]*telebot.InlineButton {
	reminderDetailDeleteBtn := telebot.InlineButton{
		Unique: "reminderDetailDeleteBtn",
		Text:   "🗑 Delete Reminder",
	}
	reminderDetailShowReminderCommandBtn := telebot.InlineButton{
		Unique: "reminderDetailShowReminderCommandBtn",
		Text:   "📄 Show Reminder Command",
	}

	return map[string]*telebot.InlineButton{
		"ReminderDetailDeleteBtn":              &reminderDetailDeleteBtn,
		"ReminderDetailShowReminderCommandBtn": &reminderDetailShowReminderCommandBtn,
	}
}

func (r *ReminderHandler) HandleRemind(m *telebot.Message) {
	log.Println("HandleRemind")
	timezone := "Europe/London"

	key, vars, err := r.reminderParser.Parse(m.Text)
	if err != nil {
		log.Println(err)
		return
	}

	ownerID := int(m.Chat.ID)
	switch key {
	case "regexOnDayMonthYear":
		fmt.Println(vars)
		fmt.Println(timezone)

		// TODO validate date not in past
		r.reminderService.AddReminderOnDayMonthYearHourMin(ownerID, m.Text, toInt(vars["day"]),
			toNumericMonth(vars["month"]), toInt(vars["year"]), 22, 16, vars["message"], timezone)
		r.b.Send(m.Chat, "reminder added regexOnDayMonthYear")
	case "regexOnDayMonthYearHourMin":
		fmt.Println(vars)
		fmt.Println(timezone)
		// TODO validate date not in past
		r.reminderService.AddReminderOnDayMonthYearHourMin(ownerID, m.Text, toInt(vars["day"]),
			toNumericMonth(vars["month"]), toInt(vars["year"]), toInt(vars["hour"]), toInt(vars["minute"]), vars["message"], timezone)
		r.b.Send(m.Chat, "reminder added regexOnDayMonthYearHourMin")
	}
}

func (r *ReminderHandler) HandleRemindList(m *telebot.Message) {
	log.Println("HandleRemindList")
	timezone := "Europe/London"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println(err)
		return
	}

	if m.Chat == nil {
		log.Println("chat is nil")
		return
	}

	ownerID := int(m.Chat.ID)
	reminders, err := r.reminderService.ListAllRemindersByOwnerID(ownerID, timezone)
	if err != nil {
		log.Println(err)
		return
	}

	remindersByStatus := map[cron.JobStatus][]ReminderListEntry{}

	for j := range reminders {
		rLE := ReminderListEntry{
			Reminder: reminders[j],
		}

		if reminders[j].Status == cron.Active {
			cronEntry := r.scheduler.GetEntryByID(reminders[j].CronID)
			rLE.NextSchedule = &cronEntry.Next
		}

		remindersByStatus[reminders[j].Status] = append(remindersByStatus[reminders[j].Status], rLE)
	}

	activeStatusList := buildReminderListText(remindersByStatus[cron.Active], loc)
	completedStatusText := buildReminderListText(remindersByStatus[cron.Completed], loc)
	inactiveStatusText := buildReminderListText(remindersByStatus[cron.Inactive], loc)

	text := ""
	if len(remindersByStatus[cron.Active]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Active)
		text += activeStatusList
		text += "\n"
	}

	if len(remindersByStatus[cron.Completed]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Completed)
		text += completedStatusText
		text += "\n"
	}

	if len(remindersByStatus[cron.Inactive]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Inactive)
		text += inactiveStatusText
		text += "\n"
	}

	if text == "" {
		text = "You have no reminders."
	}

	r.b.Send(m.Chat, text)
}

func (r *ReminderHandler) HandleRemindDelete(m *telebot.Message) {
	log.Println("HandleRemindDelete")
	timezone := "Europe/London"

	key, vars, err := r.reminderParser.Parse(m.Text)
	if err != nil {
		log.Println(err)
		return
	}

	ownerID := int(m.Chat.ID)
	switch key {
	case "regexDeleteReminder":
		// TODO implement
		fmt.Println(timezone)
		fmt.Println(vars)

		r.reminderService.DeleteReminder(ownerID, toInt(vars["reminderID"]))

		// reminderService.AddReminderOnDayMonthYearHourMinSec(m.Sender.ID, m.Text, toInt(vars["day"]), toNumericMonth(vars["month"]),
		// 	toInt(vars["year"]), 22, 16, 0, vars["message"], timezone)
		r.b.Send(m.Chat, fmt.Sprintf("Reminder %s has been deleted", vars["reminderID"]))
		return
	}

}

func (r *ReminderHandler) HandleRemindDetail(m *telebot.Message) {
	log.Println("HandleRemindDelete")
	timezone := "Europe/London"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println(err)
		return
	}

	key, vars, err := r.reminderParser.Parse(m.Text)
	if err != nil {
		log.Println(err)
		return
	}

	ownerID := int(m.Chat.ID)
	switch key {
	case "regexDetailReminder":
		// TODO implement
		fmt.Println(timezone)
		fmt.Println(vars)

		reminder, err := r.reminderService.GetReminderByOwnerID(ownerID, toInt(vars["reminderID"]), timezone)
		if err != nil {
			log.Println(err)
			return
		}

		reminderDetail := ReminderDetail{Reminder: *reminder}

		if reminder.Status == cron.Active {
			cronEntry := r.scheduler.GetEntryByID(reminder.CronID)
			reminderDetail.NextSchedule = &cronEntry.Next
		}

		remindDetailInlineKeys := [][]telebot.InlineButton{
			{*r.ReminderButtons["ReminderDetailShowReminderCommandBtn"]},
			{*r.ReminderButtons["ReminderDetailDeleteBtn"]},
		}

		//if inlineBtns != nil {
		remindDetailInlineKeys[0][0].Data = strconv.Itoa(reminder.ID)
		remindDetailInlineKeys[1][0].Data = strconv.Itoa(reminder.ID)
		r.b.Send(m.Chat, buildReminderDetailText(reminderDetail, loc), &telebot.ReplyMarkup{
			InlineKeyboard: remindDetailInlineKeys,
		})
		//} else {
		//r.b.Send(m.Chat, buildReminderDetailText(reminderDetail, loc))
		//}

		return
	}
}

func (r *ReminderHandler) HandleReminderDetailDeleteBtn(c *telebot.Callback) {
	ownerID := int(c.Message.Chat.ID)
	reminderID, err := strconv.Atoi(c.Data)
	if err != nil {
		log.Println(err)
		return
	}
	r.b.Respond(c, &telebot.CallbackResponse{
		ShowAlert: false,
	})
	r.reminderService.DeleteReminder(ownerID, reminderID)
	r.b.Send(c.Message.Chat, fmt.Sprintf("Reminder %s has been deleted", strconv.Itoa(reminderID)))
}

func (r *ReminderHandler) HandleReminderShowReminderCommandBtn(c *telebot.Callback) {
	timezone := "Europe/London"
	ownerID := int(c.Message.Chat.ID)
	reminderID, err := strconv.Atoi(c.Data)
	if err != nil {
		log.Println(err)
		return
	}

	reminder, err := r.reminderService.GetReminderByOwnerID(ownerID, reminderID, timezone)
	if err != nil {
		log.Println(err)
		return
	}
	r.b.Respond(c, &telebot.CallbackResponse{
		ShowAlert: false,
	})
	r.b.Send(c.Message.Chat, reminder.Data.Command)
}

//telegramBot.Handle("/remind", HandleRemind(telegramBot, reminderService, reminderParser))
//telegramBot.Handle("/remindlist", HandleRemindList(telegramBot, reminderService, scheduler))
//telegramBot.Handle("/reminddelete", HandleRemindDelete(telegramBot, reminderService, reminderParser))
//telegramBot.Handle("/reminddetail", HandleRemindDetail(telegramBot, reminderService, reminderParser, scheduler, remindDetailInlineKeys))

func RegisterHandlers(telegramBot *bot.Bot, handler *ReminderHandler) {
	telegramBot.Handle("/remind", handler.HandleRemind)
	telegramBot.Handle("/remindlist", handler.HandleRemindList)
	telegramBot.Handle("/reminddelete", handler.HandleRemindDelete)
	telegramBot.Handle("/reminddetail", handler.HandleRemindDetail)

	// buttons
	telegramBot.Handle(handler.ReminderButtons["ReminderDetailDeleteBtn"], handler.HandleReminderDetailDeleteBtn)
	telegramBot.Handle(handler.ReminderButtons["ReminderDetailShowReminderCommandBtn"], handler.HandleReminderShowReminderCommandBtn)

}

//func RegisterHandlers(telegramBot *bot.Bot, reminderService *ReminderService, reminderParser *parser.Parser, scheduler *cron.Scheduler) {
//	//confirmRemindDeleteBtn := telebot.InlineButton{
//	//	Unique:          "confirmRemindDeleteBtn",
//	//	Text:            "Yes",
//	//	Data:            "1",
//	//}
//	//cancelRemindDeleteBtn := telebot.InlineButton{
//	//	Unique:          "cancelRemindDeleteBtn",
//	//	Text:            "No",
//	//	Data:            "1",
//	//}
//	//remindDeleteInlineKeys := [][]telebot.InlineButton{
//	//	{confirmRemindDeleteBtn},
//	//	{cancelRemindDeleteBtn},
//	//}
//	//fmt.Println(remindDeleteInlineKeys)
//	//
//	//telegramBot.Handle(&confirmRemindDeleteBtn, func(c *telebot.Callback) {
//	//	// TODO keep global map of last item to delete
//	//	fmt.Println("c.Data")
//	//	fmt.Println(c.Data)
//	//	telegramBot.Respond(c, &telebot.CallbackResponse{
//	//		ShowAlert: false,
//	//	})
//	//	telegramBot.Send(c.Sender, "reminder deleted")
//	//})
//
//	reminderDetailDeleteBtn := telebot.InlineButton{
//		Unique: "reminderDetailDeleteBtn",
//		Text:   "🗑 Delete Reminder",
//	}
//	reminderDetailPrintCommandBtn := telebot.InlineButton{
//		Unique: "reminderDetailPrintCommandBtn",
//		Text:   "🗑 Print Reminder Command",
//	}
//	remindDetailInlineKeys := [][]telebot.InlineButton{
//		{reminderDetailDeleteBtn},
//		{reminderDetailPrintCommandBtn},
//	}
//
//	telegramBot.Handle(&reminderDetailDeleteBtn, func(c *telebot.Callback) {
//		// TODO keep global map of last item to delete
//		fmt.Println("c.Data found for reminderDetailDeleteBtn")
//		fmt.Println(c.Data)
//		telegramBot.Respond(c, &telebot.CallbackResponse{
//			ShowAlert: false,
//		})
//		telegramBot.Send(c.Sender, "clicked on reminderDetailDeleteBtn")
//	})
//
//	telegramBot.Handle(&reminderDetailPrintCommandBtn, func(c *telebot.Callback) {
//		// TODO keep global map of last item to delete
//		fmt.Println("c.Data found for reminderDetailPrintCommandBtn")
//		fmt.Println(c.Data)
//		telegramBot.Respond(c, &telebot.CallbackResponse{
//			ShowAlert: false,
//		})
//		telegramBot.Send(c.Sender, "clicked on reminderDetailPrintCommandBtn")
//	})
//
//	telegramBot.Handle("/remind", HandleRemind(telegramBot, reminderService, reminderParser))
//	telegramBot.Handle("/remindlist", HandleRemindList(telegramBot, reminderService, scheduler))
//	telegramBot.Handle("/reminddelete", HandleRemindDelete(telegramBot, reminderService, reminderParser))
//	telegramBot.Handle("/reminddetail", HandleRemindDetail(telegramBot, reminderService, reminderParser, scheduler, remindDetailInlineKeys))
//}

func HandleRemind(b *bot.Bot, reminderService *ReminderService, reminderParser *parser.Parser) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		log.Println("HandleRemind")
		timezone := "Europe/London"

		key, vars, err := reminderParser.Parse(m.Text)
		if err != nil {
			log.Println(err)
			return
		}

		ownerID := int(m.Chat.ID)
		switch key {
		case "regexOnDayMonthYear":
			fmt.Println(vars)
			fmt.Println(timezone)

			// TODO validate date not in past
			reminderService.AddReminderOnDayMonthYearHourMin(ownerID, m.Text, toInt(vars["day"]),
				toNumericMonth(vars["month"]), toInt(vars["year"]), 22, 16, vars["message"], timezone)
			b.Send(m.Chat, "reminder added regexOnDayMonthYear")
		case "regexOnDayMonthYearHourMin":
			fmt.Println(vars)
			fmt.Println(timezone)
			// TODO validate date not in past
			reminderService.AddReminderOnDayMonthYearHourMin(ownerID, m.Text, toInt(vars["day"]),
				toNumericMonth(vars["month"]), toInt(vars["year"]), toInt(vars["hour"]), toInt(vars["minute"]), vars["message"], timezone)
			b.Send(m.Chat, "reminder added regexOnDayMonthYearHourMin")
		}

	}
}

func toInt(s string) int {
	if s == "" {
		panic("empty string")
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func toNumericMonth(month string) int {
	switch month {
	case "january":
		return 1
	case "february":
		return 2
	case "march":
		return 3
	case "april":
		return 4
	case "may":
		return 5
	case "june":
		return 6
	case "july":
		return 7
	case "august":
		return 8
	case "september":
		return 9
	case "october":
		return 10
	case "november":
		return 11
	case "december":
		return 12
	default:
		return 0
	}
}

type ReminderListEntry struct {
	Reminder
	NextSchedule *time.Time
}

func HandleRemindList(b *bot.Bot, reminderService *ReminderService, scheduler *cron.Scheduler) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		log.Println("HandleRemindList")
		timezone := "Europe/London"

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			log.Println(err)
			return
		}

		if m.Chat == nil {
			log.Println("chat is nil")
			return
		}

		ownerID := int(m.Chat.ID)
		reminders, err := reminderService.ListAllRemindersByOwnerID(ownerID, timezone)
		if err != nil {
			log.Println(err)
			return
		}

		remindersByStatus := map[cron.JobStatus][]ReminderListEntry{}

		for j := range reminders {
			r := ReminderListEntry{
				Reminder: reminders[j],
			}

			if reminders[j].Status == cron.Active {
				cronEntry := scheduler.GetEntryByID(reminders[j].CronID)
				r.NextSchedule = &cronEntry.Next
			}

			remindersByStatus[reminders[j].Status] = append(remindersByStatus[reminders[j].Status], r)
		}

		activeStatusList := buildReminderListText(remindersByStatus[cron.Active], loc)
		completedStatusText := buildReminderListText(remindersByStatus[cron.Completed], loc)
		inactiveStatusText := buildReminderListText(remindersByStatus[cron.Inactive], loc)

		text := ""
		if len(remindersByStatus[cron.Active]) > 0 {
			text += fmt.Sprintf("*%s*\n", cron.Active)
			text += activeStatusList
			text += "\n"
		}

		if len(remindersByStatus[cron.Completed]) > 0 {
			text += fmt.Sprintf("*%s*\n", cron.Completed)
			text += completedStatusText
			text += "\n"
		}

		if len(remindersByStatus[cron.Inactive]) > 0 {
			text += fmt.Sprintf("*%s*\n", cron.Inactive)
			text += inactiveStatusText
			text += "\n"
		}

		if text == "" {
			text = "You have no reminders."
		}

		b.Send(m.Chat, text)
	}
}

type ReminderDetail struct {
	Reminder
	NextSchedule *time.Time
}

func HandleRemindDetail(b *bot.Bot, reminderService *ReminderService, reminderParser *parser.Parser, scheduler *cron.Scheduler, inlineBtns [][]telebot.InlineButton) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		log.Println("HandleRemindDelete")
		timezone := "Europe/London"

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			log.Println(err)
			return
		}

		key, vars, err := reminderParser.Parse(m.Text)
		if err != nil {
			log.Println(err)
			return
		}

		ownerID := int(m.Chat.ID)
		switch key {
		case "regexDetailReminder":
			// TODO implement
			fmt.Println(timezone)
			fmt.Println(vars)

			reminder, err := reminderService.GetReminderByOwnerID(ownerID, toInt(vars["reminderID"]), timezone)
			if err != nil {
				log.Println(err)
			}

			cronEntry := scheduler.GetEntryByID(reminder.CronID)
			reminderDetail := ReminderDetail{
				Reminder:     *reminder,
				NextSchedule: &cronEntry.Next,
			}

			if inlineBtns != nil {
				inlineBtns[0][0].Data = strconv.Itoa(reminder.ID)
				b.Send(m.Chat, buildReminderDetailText(reminderDetail, loc), &telebot.ReplyMarkup{
					InlineKeyboard: inlineBtns,
				})
			} else {
				b.Send(m.Chat, buildReminderDetailText(reminderDetail, loc))
			}

			return
		}
	}
}

func HandleRemindDelete(b *bot.Bot, reminderService *ReminderService, reminderParser *parser.Parser) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		log.Println("HandleRemindDelete")
		timezone := "Europe/London"

		key, vars, err := reminderParser.Parse(m.Text)
		if err != nil {
			log.Println(err)
			return
		}

		ownerID := int(m.Chat.ID)
		switch key {
		case "regexDeleteReminder":
			// TODO implement
			fmt.Println(timezone)
			fmt.Println(vars)

			reminderService.DeleteReminder(ownerID, toInt(vars["reminderID"]))

			// reminderService.AddReminderOnDayMonthYearHourMinSec(m.Sender.ID, m.Text, toInt(vars["day"]), toNumericMonth(vars["month"]),
			// 	toInt(vars["year"]), 22, 16, 0, vars["message"], timezone)
			b.Send(m.Chat, fmt.Sprintf("Reminder %s has been deleted", vars["reminderID"]))
			return
		}
	}
}

// https://unicode.org/emoji/charts/full-emoji-list.html
func buildReminderListText(list []ReminderListEntry, loc *time.Location) string {
	text := ""
	for j := range list {
		if list[j].NextSchedule != nil {
			text += fmt.Sprintf(
				"%d) %s \n[📋[/reminddetail_%d]]\n",
				list[j].ID,
				list[j].Data.Message,
				list[j].ID,
			)
		} else {
			text += fmt.Sprintf(
				"%d) %s \n[📋[/reminddetail_%d]]\n",
				list[j].ID,
				list[j].Data.Message,
				list[j].ID,
			)
		}

	}
	return text
}

// TODO convert this to a text template?
func buildReminderDetailText(reminder ReminderDetail, loc *time.Location) string {
	text := fmt.Sprintf("*Id*: %d\n", reminder.ID)
	text += fmt.Sprintf("*Status*: %s\n", reminder.Status)
	text += fmt.Sprintf("*Message*: %s\n", reminder.Data.Message)
	text += fmt.Sprintf("*Command*: %s\n", reminder.Data.Command)
	if reminder.NextSchedule != nil {
		text += fmt.Sprintf("*Next Schedule*: %s\n", reminder.NextSchedule.In(loc).Format(time.RFC1123))
	}
	return text
}

func (r *ReminderHandler) HandleNewRemindDayMonthYear(c *bot.Context) error {
	// TODO validate date not in past
	timezone := "Europe/London"
	day := toInt(c.Param("day"))
	month := toNumericMonth(c.Param("month"))
	year := toInt(c.Param("year"))
	message := c.Param("message")

	r.reminderService.AddReminderOnDayMonthYearHourMin(c.OwnerID, c.Text, day,
		month, year, 22, 16, message, timezone)
	r.b.Send(c.Chat, "reminder added regexOnDayMonthYear")
	return nil
}

func (r *ReminderHandler) HandleNewRemindDayMonthYearHourMin(c *bot.Context) error {
	// TODO validate date not in past
	timezone := "Europe/London"
	day := toInt(c.Param("day"))
	month := toNumericMonth(c.Param("month"))
	year := toInt(c.Param("year"))
	hour := toInt(c.Param("hour"))
	minute := toInt(c.Param("minute"))
	message := c.Param("message")

	r.reminderService.AddReminderOnDayMonthYearHourMin(c.OwnerID, c.Text, day,
		month, year, hour, minute, message, timezone)
	r.b.Send(c.Chat, "reminder added regexOnDayMonthYearHourMin")
	return nil
}

func (r *ReminderHandler) HandleNewRemindList(c *bot.Context) error {
	log.Println("HandleRemindList")
	timezone := "Europe/London"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println(err)
		return err
	}

	reminders, err := r.reminderService.ListAllRemindersByOwnerID(c.OwnerID, timezone)
	if err != nil {
		log.Println(err)
		return err
	}

	remindersByStatus := map[cron.JobStatus][]ReminderListEntry{}

	for j := range reminders {
		rLE := ReminderListEntry{
			Reminder: reminders[j],
		}

		if reminders[j].Status == cron.Active {
			cronEntry := r.scheduler.GetEntryByID(reminders[j].CronID)
			rLE.NextSchedule = &cronEntry.Next
		}

		remindersByStatus[reminders[j].Status] = append(remindersByStatus[reminders[j].Status], rLE)
	}

	activeStatusList := buildReminderListText(remindersByStatus[cron.Active], loc)
	completedStatusText := buildReminderListText(remindersByStatus[cron.Completed], loc)
	inactiveStatusText := buildReminderListText(remindersByStatus[cron.Inactive], loc)

	text := ""
	if len(remindersByStatus[cron.Active]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Active)
		text += activeStatusList
		text += "\n"
	}

	if len(remindersByStatus[cron.Completed]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Completed)
		text += completedStatusText
		text += "\n"
	}

	if len(remindersByStatus[cron.Inactive]) > 0 {
		text += fmt.Sprintf("*%s*\n", cron.Inactive)
		text += inactiveStatusText
		text += "\n"
	}

	if text == "" {
		text = "You have no reminders."
	}

	r.b.Send(c.Chat, text)
	return nil
}

func (r *ReminderHandler) HandleNewReminderDetailDeleteBtn(c *bot.Context) error {
	reminderID, err := strconv.Atoi(c.Callback.Data)
	if err != nil {
		log.Println(err)
		return err
	}
	r.b.Respond(c.Callback, &telebot.CallbackResponse{
		ShowAlert: false,
	})
	r.reminderService.DeleteReminder(c.OwnerID, reminderID)
	r.b.Send(c.Chat, fmt.Sprintf("Reminder %s has been deleted", strconv.Itoa(reminderID)))
	return nil
}

func (r *ReminderHandler) HandleNewReminderShowReminderCommandBtn(c *bot.Context) error {
	timezone := "Europe/London"
	reminderID, err := strconv.Atoi(c.Callback.Data)
	if err != nil {
		log.Println(err)
		return err
	}

	reminder, err := r.reminderService.GetReminderByOwnerID(c.OwnerID, reminderID, timezone)
	if err != nil {
		log.Println(err)
		return err
	}
	r.b.Respond(c.Callback, &telebot.CallbackResponse{
		ShowAlert: false,
	})
	r.b.Send(c.Chat, reminder.Data.Command)
	return nil
}

func (r *ReminderHandler) HandleNewRemindDelete(c *bot.Context) error {
	log.Println("HandleRemindDelete")
	reminderID := toInt(c.Param("reminderID"))

	r.reminderService.DeleteReminder(c.OwnerID, reminderID)

	// reminderService.AddReminderOnDayMonthYearHourMinSec(m.Sender.ID, m.Text, toInt(vars["day"]), toNumericMonth(vars["month"]),
	// 	toInt(vars["year"]), 22, 16, 0, vars["message"], timezone)
	r.b.Send(c.Chat, fmt.Sprintf("Reminder %d has been deleted", reminderID))
	return nil
}

func (r *ReminderHandler) HandleNewRemindDetail(c *bot.Context) error {
	log.Println("HandleRemindDelete")
	timezone := "Europe/London"
	reminderID := toInt(c.Param("reminderID"))

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println(err)
		return err
	}

	reminder, err := r.reminderService.GetReminderByOwnerID(c.OwnerID, reminderID, timezone)
	if err != nil {
		log.Println(err)
		return err
	}

	reminderDetail := ReminderDetail{Reminder: *reminder}

	if reminder.Status == cron.Active {
		cronEntry := r.scheduler.GetEntryByID(reminder.CronID)
		reminderDetail.NextSchedule = &cronEntry.Next
	}

	remindDetailInlineKeys := [][]telebot.InlineButton{
		{*r.ReminderButtons["ReminderDetailShowReminderCommandBtn"]},
		{*r.ReminderButtons["ReminderDetailDeleteBtn"]},
	}

	//if inlineBtns != nil {
	remindDetailInlineKeys[0][0].Data = strconv.Itoa(reminder.ID)
	remindDetailInlineKeys[1][0].Data = strconv.Itoa(reminder.ID)
	r.b.Send(c.Chat, buildReminderDetailText(reminderDetail, loc), &telebot.ReplyMarkup{
		InlineKeyboard: remindDetailInlineKeys,
	})
	//} else {
	//r.b.Send(m.Chat, buildReminderDetailText(reminderDetail, loc))
	//}

	return nil

}
