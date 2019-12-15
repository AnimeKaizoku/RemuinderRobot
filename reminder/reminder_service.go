package reminder

import (
	"errors"
	"fmt"
	"time"

	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
)

type ReminderService struct {
	b             *bot.Bot
	scheduler     *cron.Scheduler
	ReminderStore *ReminderStore
}

func NewReminderService(b *bot.Bot, scheduler *cron.Scheduler, reminderStore *ReminderStore) *ReminderService {
	return &ReminderService{
		b:             b,
		scheduler:     scheduler,
		ReminderStore: reminderStore,
	}
}

// // timezones??
// func (s *ReminderService) AddReminderEveryDayWeekForNMonths(recipientID int, command, timezone, dayWeek string, nMonths int, message string) error {
// 	reminder := Reminder{
// 		RecipientID: recipientID,
// 		Command:     command,
// 		Message:     message,
// 		Dates:       []ReminderDate{},
// 	}

// 	fmt.Printf("Reminder: %#v", reminder)
// 	// fmt.Println(time.Now().Format(time.RFC3339))

// 	// dayToFind := "Tuesday"
// 	loc, err := time.LoadLocation(timezone)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	calcTime := time.Now().In(loc)
// 	untilTime := time.Now().In(loc).AddDate(0, nMonths, 0) // 1 month
// 	for i := 0; i <= 7; i++ {
// 		if strings.ToLower(calcTime.Weekday().String()) == strings.ToLower(dayWeek) {
// 			fmt.Println(calcTime.Format(time.RFC1123Z))
// 			break
// 		}
// 		calcTime = calcTime.AddDate(0, 0, 1)
// 	}

// 	for i := calcTime; i.Before(untilTime); i = i.AddDate(0, 0, 7) {
// 		reminder.Dates = append(reminder.Dates, ReminderDate{DateTime: calcTime})
// 		calcTime = calcTime.AddDate(0, 0, 7)
// 	}

// 	for i := range reminder.Dates {
// 		fmt.Printf("Reminder date: %#v\n", reminder.Dates[i].DateTime.Format(time.UnixDate))
// 	}

// 	return nil
// }

func (s *ReminderService) AddReminderOnDayMonthYearHourMin(ownerID int, command string, day, month, year, hour, min int, message, timezone string) (int, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, err
	}

	t := time.Date(year, time.Month(month), day, hour, min, 0, 0, loc)
	tInUTC := MustTimeIn(t, "UTC")

	timeNow := time.Now().In(time.UTC).Add(time.Minute * 2)
	fmt.Printf("Scheduling next reminder for %s\n", timeNow.Format(time.RFC1123))

	// TODO replace minute for testing
	schedule := fmt.Sprintf("%d %d %d %d *", timeNow.Minute(), timeNow.Hour(), tInUTC.Day(), tInUTC.Month())

	reminder := &Reminder{
		Job: cron.Job{
			OwnerID:     ownerID,
			Schedule:    schedule,
			Type:        cron.HealthCheck,
			Status:      cron.Active,
			RunOnlyOnce: true,
		},
		Data: ReminderData{
			RecipientID: ownerID,
			Message:     message,
			Command:     command,
		},
	}

	reminderCronID, err := s.scheduler.Add(reminder.Job.Schedule, NewReminderCronFunc(s, s.b, reminder))
	if err != nil {
		return 0, err
	}

	reminder.CronID = reminderCronID
	ID, err := s.ReminderStore.CreateReminder(reminder)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *ReminderService) ListAllRemindersByUserAndGroup(recipientID int, timezone string) (map[int][]Reminder, error) {
	//loc, err := time.LoadLocation(timezone)
	//if err != nil {
	//	return nil, err
	//}

	remindersByUserAndGroup, err := s.ReminderStore.GetAllRemindersByUserAndGroup()
	if err != nil {
		return nil, err
	}

	return remindersByUserAndGroup, nil
}

func (s *ReminderService) ListAllRemindersByOwnerID(ownerID int, timezone string) ([]Reminder, error) {
	//loc, err := time.LoadLocation(timezone)
	//if err != nil {
	//	return nil, err
	//}

	reminders, err := s.ReminderStore.GetAllRemindersByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (s *ReminderService) GetReminderByOwnerID(ownerID, ID int, timezone string) (*Reminder, error) {
	//loc, err := time.LoadLocation(timezone)
	//if err != nil {
	//	return nil, err
	//}

	reminder, err := s.ReminderStore.GetReminder(ownerID, ID)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (s *ReminderService) DeleteReminder(ownerID, ID int) error {
	reminder, err := s.ReminderStore.GetReminder(ownerID, ID)
	if err != nil {
		return err
	}

	if ownerID != reminder.OwnerID {
		return errors.New("unauthorised to delete reminder")
	}

	s.scheduler.Remove(reminder.CronID)

	return s.ReminderStore.DeleteReminder(reminder.OwnerID, ID)
}

func (s *ReminderService) LoadExistingSchedules() (int, error) {
	count := 0
	rmdrListByUserAndGroup, err := s.ReminderStore.GetAllRemindersByUserAndGroup()
	if err != nil {
		return 0, err
	}

	for userGroupID := range rmdrListByUserAndGroup {
		for i := range rmdrListByUserAndGroup[userGroupID] {
			fmt.Println(rmdrListByUserAndGroup[userGroupID][i].Schedule)
			fmt.Println(rmdrListByUserAndGroup[userGroupID][i].RunOnlyOnce)
			if rmdrListByUserAndGroup[userGroupID][i].Status == cron.Active {
				reminderID, err := s.scheduler.Add(rmdrListByUserAndGroup[userGroupID][i].Job.Schedule, NewReminderCronFunc(s, s.b, &rmdrListByUserAndGroup[userGroupID][i]))
				if err != nil {
					return 0, err
				}

				rmdrListByUserAndGroup[userGroupID][i].CronID = reminderID
				err = s.ReminderStore.UpdateReminder(&rmdrListByUserAndGroup[userGroupID][i])
				if err != nil {
					return 0, err
				}

				count++
			}
		}
	}

	return len(rmdrListByUserAndGroup), nil
}

func (s *ReminderService) Complete(r *Reminder) error {
	r.Status = cron.Completed

	err := s.ReminderStore.UpdateReminder(r)
	if err != nil {
		return err
	}

	s.scheduler.Remove(r.CronID)
	return nil
}

// func (s *ReminderService) DeleteReminder(recipientID, reminderID int) error {
// 	// reminder, err := s.ReminderStore.GetReminder(reminderID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// for i := range reminder.Dates {
// 	// 	s.ReminderByDateStore.
// 	// }
// 	return nil
// }

// func TimeIn(t time.Time, name string) (time.Time, error) {
// 	loc, err := time.LoadLocation(name)
// 	if err == nil {
// 		t = t.In(loc)
// 	}
// 	return t, err
// }

func MustTimeIn(t time.Time, name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}

// func ParseInLocation(layout, value, locName string) (time.Time, error) {
// 	loc, err := time.LoadLocation(locName)
// 	if err != nil {
// 		return time.Now(), err
// 	}
// 	locTime, err := time.ParseInLocation(layout, value, loc)
// 	if err != nil {
// 		return time.Now(), err
// 	}

// 	return locTime, nil
// }

// func mainbabababab() {
// 	fmt.Println("Hello, playground")

// 	layout := time.RFC3339
// 	str := "2014-11-12T11:45:26.371Z"
// 	t, err := ParseInLocation(layout, str, "Europe/London")

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(t)
// 	fmt.Println(TimeIn(t, "Europe/Rome"))
// }
