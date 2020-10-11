package db

import (
	"fmt"
	"strconv"

	"github.com/enrico5b1b4/telegram-bot/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"go.etcd.io/bbolt"
)

// SetupDB creates a root reminders bucket
// Buckets are then created in the root bucket for each chat
func SetupDB(filename string, chats []int) (*bbolt.DB, error) {
	db, err := bbolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %#v", err)
	}

	updateErr := db.Update(func(tx *bbolt.Tx) error {
		rootReminderBucket, err := tx.CreateBucketIfNotExists(reminder.RemindersBucket)
		if err != nil {
			return fmt.Errorf("could not create reminders bucket: %#v", err)
		}

		// create individual buckets for chats
		for i := range chats {
			_, err = rootReminderBucket.CreateBucketIfNotExists(itob(chats[i]))
			if err != nil {
				return fmt.Errorf("could not create reminders bucket for chat: %d %#v", chats[i], err)
			}
		}

		_, err = tx.CreateBucketIfNotExists(chatpreference.ChatPreferencesBucket)
		if err != nil {
			return fmt.Errorf("could not create chat preferences bucket: %#v", err)
		}

		return nil
	})
	if updateErr != nil {
		return nil, fmt.Errorf("could not set up buckets, %#v", updateErr)
	}

	return db, nil
}

// itob converts int to []byte
func itob(v int) []byte {
	return []byte(strconv.FormatInt(int64(v), 10))
}
