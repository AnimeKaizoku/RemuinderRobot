package reminder

//go:generate mockgen -destination=./mocks/mock_Storer.go -package=mocks github.com/enrico5b1b4/telegram-bot/reminder Storer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var RemindersBucket = []byte("reminders")
var ErrNotFound = errors.New("reminder not found")

type Storer interface {
	CreateReminder(r *Reminder) (int, error)
	UpdateReminder(r *Reminder) error
	DeleteReminder(chatID, ID int) error
	GetReminder(chatID, ID int) (*Reminder, error)
	GetAllRemindersByChat() (map[int][]Reminder, error)
	GetAllRemindersByChatID(chatID int) ([]Reminder, error)
}

type Store struct {
	db *bolt.DB
}

func NewStore(db *bolt.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateReminder(r *Reminder) (int, error) {
	var ID int

	err := s.db.Update(func(tx *bolt.Tx) error {
		reminderBucket := tx.Bucket(RemindersBucket)
		chatBucket := reminderBucket.Bucket(itob(r.ChatID))

		id, err := chatBucket.NextSequence()
		if err != nil {
			return err
		}
		r.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		return chatBucket.Put(itob(r.ID), buf)
	})
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *Store) UpdateReminder(r *Reminder) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		reminderBucket := tx.Bucket(RemindersBucket)
		chatBucket := reminderBucket.Bucket(itob(r.ChatID))

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		return chatBucket.Put(itob(r.ID), buf)
	})
}

func (s *Store) GetAllRemindersByChat() (map[int][]Reminder, error) {
	remindersByChat := map[int][]Reminder{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RemindersBucket)

		return b.ForEach(func(chatID, v []byte) error {
			chatBucket := b.Bucket(chatID)

			return chatBucket.ForEach(func(k1, v1 []byte) error {
				var reminder Reminder

				err := json.Unmarshal(v1, &reminder)
				if err != nil {
					return err
				}

				remindersByChat[btoi(chatID)] = append(remindersByChat[btoi(chatID)], reminder)
				return nil
			})
		})
	})
	if err != nil {
		return nil, err
	}

	return remindersByChat, nil
}

func (s *Store) GetAllRemindersByChatID(chatID int) ([]Reminder, error) {
	remindersByChat := []Reminder{}

	err := s.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		chatBucket := rootBucket.Bucket(itob(chatID))

		return chatBucket.ForEach(func(k1, v1 []byte) error {
			var reminder Reminder

			err := json.Unmarshal(v1, &reminder)
			if err != nil {
				return err
			}

			remindersByChat = append(remindersByChat, reminder)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return remindersByChat, nil
}

func (s *Store) GetReminder(chatID, id int) (*Reminder, error) {
	var reminder Reminder

	err := s.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		chatBucket := rootBucket.Bucket(itob(chatID))
		v := chatBucket.Get(itob(id))
		if v == nil {
			return ErrNotFound
		}

		err := json.Unmarshal(v, &reminder)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (s *Store) DeleteReminder(chatID, id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		chatBucket := rootBucket.Bucket(itob(chatID))

		return chatBucket.Delete(itob(id))
	})
}

func (s *Store) Debug() error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RemindersBucket)

		return b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
	})
}

// itob converts int to []byte
func itob(v int) []byte {
	return []byte(strconv.FormatInt(int64(v), 10))
}

func btoi(v []byte) int {
	byteToInt, _ := strconv.Atoi(string(v))
	return byteToInt
}
