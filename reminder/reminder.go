package reminder

import (
	"time"

	"github.com/enrico5b1b4/telegram-bot/cron"
)

type Reminder struct {
	cron.Job
	Data Data `json:"data"`
}

type Data struct {
	RecipientID int    `json:"recipient_id"`
	Command     string `json:"command"`
	Message     string `json:"message"`
}

type DateTime struct {
	DayOfMonth int
	DayOfWeek  string
	Month      int
	Hour       int
	Minute     int
}

type RepeatableDateTime struct {
	DayOfMonth string
	DayOfWeek  string
	Month      string
	Hour       string
	Minute     string
}

type AmountDateTime struct {
	Minutes int
	Hours   int
	Days    int
}

// TODO: Find a better name for these variables...
type WordTimes int

const (
	Today WordTimes = iota + 1
	Tomorrow
)

func (w WordTimes) String() string {
	return [...]string{"", "Today", "Tomorrow"}[w]
}

type WordDateTime struct {
	When   WordTimes
	Hour   int
	Minute int
}

type NextScheduleChatTime struct {
	Time     time.Time
	Location *time.Location
}
