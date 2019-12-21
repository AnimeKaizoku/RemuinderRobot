package remindlist

//go:generate mockgen -source=$GOFILE -destination=$PWD/reminder/remindlist/mocks/${GOFILE} -package=mocks

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/cron"
	"github.com/enrico5b1b4/telegram-bot/reminder"
)

var jobStatusIndexMap = map[cron.JobStatus]int{
	cron.Active:    0,
	cron.Inactive:  1,
	cron.Completed: 2,
}

type ByJobStatusList struct {
	Status  cron.JobStatus
	Entries []ListEntryGroup
}

type ListEntryGroup struct {
	Time    *time.Time
	Entries []ListEntry
}

type ListEntry struct {
	reminder.Reminder
	NextSchedule *time.Time
}

const maxLengthMessageEntry = 20

type Servicer interface {
	GetRemindersByChatID(chatID int) ([]ByJobStatusList, error)
	RemoveCompletedReminders(chatID int) error
}

type Service struct {
	reminderStore       reminder.Storer
	scheduler           cron.Scheduler
	chatPreferenceStore chatpreference.Storer
}

func NewService(reminderStore reminder.Storer, scheduler cron.Scheduler, chatPreferenceStore chatpreference.Storer) *Service {
	return &Service{
		reminderStore:       reminderStore,
		scheduler:           scheduler,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *Service) GetRemindersByChatID(chatID int) ([]ByJobStatusList, error) {
	reminders, err := s.reminderStore.GetAllRemindersByChatID(chatID)
	if err != nil {
		return nil, err
	}

	chatPreference, err := s.chatPreferenceStore.GetChatPreference(chatID)
	if err != nil {
		return nil, err
	}

	chatLocalTimezone, err := time.LoadLocation(chatPreference.TimeZone)
	if err != nil {
		return nil, err
	}

	// index 0 = cron.Active, index 1 = cron.Inactive, index 2 = cron.Completed
	remindersByStatusAndTime := []ByJobStatusList{
		{Status: cron.Active, Entries: []ListEntryGroup{}},
		{Status: cron.Inactive, Entries: []ListEntryGroup{}},
		{Status: cron.Completed, Entries: []ListEntryGroup{}},
	}
	for _, rem := range reminders {
		rLE := ListEntry{
			Reminder: rem,
		}
		rLE.Reminder.Data.Message = truncateString(rLE.Reminder.Data.Message, maxLengthMessageEntry)
		jobStatusIndex := jobStatusIndexMap[rem.Status]
		var timeKey *time.Time

		// if reminder is still active then fetch next schedule date to display sorted entries to user
		if rem.Status == cron.Active {
			cronEntry := s.scheduler.GetEntryByID(rem.CronID)
			nextSchedule := cronEntry.Next.In(chatLocalTimezone)

			rLE.NextSchedule = &nextSchedule
			var err error
			timeKey, err = createTimeKey(nextSchedule)
			if err != nil {
				return nil, err
			}
		}

		// insert reminder in correct group based on timeKey
		inserted := false
		for i, group := range remindersByStatusAndTime[jobStatusIndex].Entries {
			if !inserted && equalTime(group.Time, timeKey) {
				if remindersByStatusAndTime[jobStatusIndex].Entries[i].Entries == nil {
					remindersByStatusAndTime[jobStatusIndex].Entries[i].Entries = []ListEntry{}
				}

				remindersByStatusAndTime[jobStatusIndex].Entries[i].Entries = insertEntrySorted(
					remindersByStatusAndTime[jobStatusIndex].Entries[i].Entries,
					rLE,
				)
				inserted = true
				break
			}
		}

		// if group doesn't exist then create it and then add the reminder to its entries
		if !inserted {
			var index int
			remindersByStatusAndTime[jobStatusIndex].Entries, index = insertListEntrySorted(
				remindersByStatusAndTime[jobStatusIndex].Entries,
				ListEntryGroup{
					Time:    timeKey,
					Entries: []ListEntry{},
				})

			remindersByStatusAndTime[jobStatusIndex].Entries[index].Entries = insertEntrySorted(
				remindersByStatusAndTime[jobStatusIndex].Entries[index].Entries,
				rLE,
			)
		}
	}

	return remindersByStatusAndTime, nil
}

func (s *Service) RemoveCompletedReminders(chatID int) error {
	reminders, err := s.reminderStore.GetAllRemindersByChatID(chatID)
	if err != nil {
		return err
	}

	for i := range reminders {
		if reminders[i].Status == cron.Completed {
			err := s.reminderStore.DeleteReminder(chatID, reminders[i].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func equalTime(val1, val2 *time.Time) bool {
	if val1 == nil && val2 == nil {
		return true
	}

	if val1 != nil && val2 != nil {
		return *val1 == *val2
	}

	return false
}

func truncateString(message string, length int) string {
	runes := bytes.Runes([]byte(message))
	if len(runes) > length {
		return fmt.Sprintf("%s...", string(runes[:length]))
	}

	return string(runes)
}

//nolint: gocritic
func insertListEntrySorted(list []ListEntryGroup, entry ListEntryGroup) ([]ListEntryGroup, int) {
	index := sort.Search(len(list), func(i int) bool { return list[i].Time.Unix() >= entry.Time.Unix() })
	list = append(list, ListEntryGroup{})
	copy(list[index+1:], list[index:])
	list[index] = entry

	return list, index
}

// nolint:gocritic
func insertEntrySorted(list []ListEntry, entry ListEntry) []ListEntry {
	index := sort.Search(len(list), func(i int) bool {
		if list[i].NextSchedule != nil {
			return list[i].NextSchedule.Unix() >= entry.NextSchedule.Unix()
		}
		return true
	})
	list = append(list, ListEntry{})
	copy(list[index+1:], list[index:])
	list[index] = entry

	return list
}

func createTimeKey(tm time.Time) (*time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, tm.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	return &t, nil
}
