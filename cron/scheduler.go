package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

type Runner interface {
	Run()
}

type Entry struct {
	ID   int
	Next time.Time
	Prev time.Time
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		c: cron.New(cron.WithLocation(time.UTC)),
	}
}

func (s *Scheduler) Start() {
	s.c.Start()

	return
}

func (s *Scheduler) Add(spec string, cmd func()) (int, error) {
	entryId, err := s.c.AddFunc(spec, cmd)

	for i := range s.c.Entries() {
		fmt.Printf("%#v\n", s.c.Entries()[i])
	}

	return int(entryId), err
}

func (s *Scheduler) Remove(ID int) {
	s.c.Remove(cron.EntryID(ID))
}

func (s *Scheduler) GetAllEntries(ID int) []Entry {
	var entries []Entry

	for _, e := range s.c.Entries() {
		entries = append(entries, convertToEntry(e))
	}

	return entries
}

func (s *Scheduler) GetEntryByID(ID int) Entry {
	cronEntry := s.c.Entry(cron.EntryID(ID))

	return convertToEntry(cronEntry)
}

func convertToEntry(entry cron.Entry) Entry {
	return Entry{
		ID:   int(entry.ID),
		Next: entry.Next,
		Prev: entry.Prev,
	}
}
