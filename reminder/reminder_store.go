package reminder

import (
	"encoding/json"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var RemindersBucket = []byte("reminders")

type ReminderStore struct {
	db *bolt.DB
}

func NewReminderStore(db *bolt.DB) *ReminderStore {
	return &ReminderStore{db: db}
}

// CreateUser saves u to the store. The new user ID is set on u once the data is persisted.
func (s *ReminderStore) CreateReminder(r *Reminder) (int, error) {
	var ID int

	err := s.db.Update(func(tx *bolt.Tx) error {
		reminderBucket := tx.Bucket(RemindersBucket)
		ownerBucket := reminderBucket.Bucket(itob(r.OwnerID))

		id, err := ownerBucket.NextSequence()
		if err != nil {
			return err
		}
		r.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		return ownerBucket.Put(itob(r.ID), buf)
	})
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *ReminderStore) UpdateReminder(r *Reminder) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		reminderBucket := tx.Bucket(RemindersBucket)
		ownerBucket := reminderBucket.Bucket(itob(r.OwnerID))

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		return ownerBucket.Put(itob(r.ID), buf)
	})
}

func (s *ReminderStore) GetAllRemindersByUserAndGroup() (map[int][]Reminder, error) {
	remindersByUserAndGroup := map[int][]Reminder{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RemindersBucket)

		s.Debug()

		return b.ForEach(func(ownerID, v []byte) error {
			fmt.Printf("ownerID=%s, value=%s\n", ownerID, v)
			userGroupBucket := b.Bucket(ownerID)

			return userGroupBucket.ForEach(func(k1, v1 []byte) error {
				var reminder Reminder
				fmt.Printf("key=%s, value=%s\n", ownerID, v)

				err := json.Unmarshal(v1, &reminder)
				if err != nil {
					return err
				}

				remindersByUserAndGroup[btoi(ownerID)] = append(remindersByUserAndGroup[btoi(ownerID)], reminder)
				return nil
			})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("remindersByUserAndGroup: %#v\n", remindersByUserAndGroup)

	return remindersByUserAndGroup, nil
}

func (s *ReminderStore) GetAllRemindersByOwnerID(ownerID int) ([]Reminder, error) {
	remindersByUserAndGroup := []Reminder{}

	err := s.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		ownerBucket := rootBucket.Bucket(itob(ownerID))

		return ownerBucket.ForEach(func(k1, v1 []byte) error {
			var reminder Reminder

			err := json.Unmarshal(v1, &reminder)
			if err != nil {
				return err
			}

			remindersByUserAndGroup = append(remindersByUserAndGroup, reminder)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return remindersByUserAndGroup, nil
}

func (s *ReminderStore) GetReminder(ownerID, ID int) (*Reminder, error) {
	var reminder Reminder

	err := s.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		ownerBucket := rootBucket.Bucket(itob(ownerID))
		v := ownerBucket.Get(itob(ID))

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

func (s *ReminderStore) DeleteReminder(ownerID, ID int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(RemindersBucket)
		ownerBucket := rootBucket.Bucket(itob(ownerID))

		return ownerBucket.Delete(itob(ID))
	})
}

func (s *ReminderStore) Debug() error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(RemindersBucket)

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})

		return nil
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
