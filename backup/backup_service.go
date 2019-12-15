package backup

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/cron"
)

type BackupService struct {
	b              *bot.Bot
	scheduler      *cron.Scheduler
	backupStore    *BackupStore
	dbxFilesClient files.Client
}

func NewBackupService(b *bot.Bot, scheduler *cron.Scheduler, backupStore *BackupStore, dbxFilesClient files.Client) *BackupService {
	return &BackupService{
		b:              b,
		scheduler:      scheduler,
		backupStore:    backupStore,
		dbxFilesClient: dbxFilesClient,
	}
}
