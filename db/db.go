package db

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/telegram-bot/backup"
	"github.com/enrico5b1b4/telegram-bot/healthcheck"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"go.etcd.io/bbolt"
)

func SetupDB(filename string, usersAndGroups []int) (*bbolt.DB, error) {
	db, err := bbolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %#v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		//tx.DeleteBucket(reminder.RemindersBucket)
		rootReminderBucket, err := tx.CreateBucketIfNotExists(reminder.RemindersBucket)
		if err != nil {
			return fmt.Errorf("could not create reminders bucket: %#v", err)
		}

		// create individual buckets for users and groups
		for i := range usersAndGroups {
			_, err := rootReminderBucket.CreateBucketIfNotExists(itob(usersAndGroups[i]))
			if err != nil {
				return fmt.Errorf("could not create reminders bucket for user/group: %d %#v", usersAndGroups[i], err)
			}
		}

		//tx.DeleteBucket(healthcheck.HealthchecksBucket)
		_, err = tx.CreateBucketIfNotExists(healthcheck.HealthchecksBucket)
		if err != nil {
			return fmt.Errorf("could not create healthchecks bucket: %#v", err)
		}

		//tx.DeleteBucket(backup.BackupsBucket)
		_, err = tx.CreateBucketIfNotExists(backup.BackupsBucket)
		if err != nil {
			return fmt.Errorf("could not create backups bucket: %#v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %#v", err)
	}
	fmt.Println("DB Setup Done")

	return db, nil
}

// itob converts int to []byte
func itob(v int) []byte {
	return []byte(strconv.FormatInt(int64(v), 10))
}
