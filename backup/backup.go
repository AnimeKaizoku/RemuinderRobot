package backup

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/enrico5b1b4/telegram-bot/cron"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Backup struct {
	cron.Job
	Data BackupData `json:"data"`
}

type BackupData struct {
	Path        string `json:"path"`
	RecipientID int    `json:"recipient_id"`
}

func NewBackupCronFunc(b *tb.Bot, h *Backup, dbxFilesClient files.Client) func() {
	return func() {

		//b.Send(&tb.User{ID: h.Data}, h.Data.Message)
	}
}
