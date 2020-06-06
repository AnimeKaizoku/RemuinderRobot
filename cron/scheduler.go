package cron

//go:generate mockgen -destination=./mocks/mock_Scheduler.go -package=mocks github.com/enrico5b1b4/telegram-bot/cron Scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	Add(spec string, cmd func()) (int, error)
	Remove(ID int)
	GetEntryByID(ID int) Entry
	Start()
}

type JobScheduler struct {
	c *cron.Cron
}

type Entry struct {
	ID   int
	Next time.Time
	Prev time.Time
}

func NewScheduler() *JobScheduler {
	return &JobScheduler{
		c: cron.New(),
	}
}

func (s *JobScheduler) Start() {
	s.c.Start()
}

func (s *JobScheduler) Add(spec string, cmd func()) (int, error) {
	entryID, err := s.c.AddFunc(spec, cmd)

	return int(entryID), err
}

func (s *JobScheduler) Remove(id int) {
	s.c.Remove(cron.EntryID(id))
}

func (s *JobScheduler) GetAllEntries(id int) []Entry {
	var entries []Entry

	for _, e := range s.c.Entries() {
		entries = append(entries, convertToEntry(e))
	}

	return entries
}

func (s *JobScheduler) GetEntryByID(id int) Entry {
	cronEntry := s.c.Entry(cron.EntryID(id))

	return convertToEntry(cronEntry)
}

// nolint:gocritic
func convertToEntry(entry cron.Entry) Entry {
	return Entry{
		ID:   int(entry.ID),
		Next: entry.Next,
		Prev: entry.Prev,
	}
}
