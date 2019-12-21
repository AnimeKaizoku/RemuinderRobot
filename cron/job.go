package cron

import "time"

type JobType int

const (
	HealthCheck JobType = iota + 1
	Backup
	Reminder
)

func (j JobType) String() string {
	return [...]string{"", "HealthCheck", "Backup", "Reminder"}[j]
}

type JobStatus int

const (
	Active JobStatus = iota + 1
	Inactive
	Completed
)

func (j JobStatus) String() string {
	return [...]string{"", "Active", "Inactive", "Completed"}[j]
}

type Job struct {
	ID             int                `json:"id"`
	CronID         int                `json:"cron_id"`
	ChatID         int                `json:"owner_id"`
	Schedule       string             `json:"schedule"`
	Type           JobType            `json:"type"`
	Status         JobStatus          `json:"status"`
	RunOnlyOnce    bool               `json:"run_only_once"`
	RepeatSchedule *JobRepeatSchedule `json:"repeat_schedule"`
	CompletedAt    *time.Time         `json:"completed_at"`
}

type JobRepeatSchedule struct {
	Minutes int `json:"minutes"`
	Hours   int `json:"hours"`
	Days    int `json:"days"`
	Months  int `json:"months"`
}
