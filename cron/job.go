package cron

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
	ID          int       `json:"id"`
	CronID      int       `json:"cron_id"`
	OwnerID     int       `json:"owner_id"`
	Schedule    string    `json:"schedule"`
	Type        JobType   `json:"type"`
	Status      JobStatus `json:"status"`
	RunOnlyOnce bool      `json:"run_only_once"`
}
