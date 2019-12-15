package healthcheck

import (
	"encoding/json"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var HealthchecksBucket = []byte("healthchecks")

type HealthcheckStore struct {
	db *bolt.DB
}

func NewHealthcheckStore(db *bolt.DB) *HealthcheckStore {
	return &HealthcheckStore{db: db}
}

func (s *HealthcheckStore) CreateHealthcheck(h *Healthcheck) (int, error) {
	var ID int

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		h.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(h)
		if err != nil {
			return err
		}

		return b.Put(itob(h.ID), buf)
	})
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *HealthcheckStore) UpdateHealthcheck(h *Healthcheck) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)

		buf, err := json.Marshal(h)
		if err != nil {
			return err
		}

		return b.Put(itob(h.ID), buf)
	})
}

func (s *HealthcheckStore) GetHealthcheck(ID int) (*Healthcheck, error) {
	var healthcheck Healthcheck

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)
		v := b.Get(itob(ID))

		err := json.Unmarshal(v, &healthcheck)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &healthcheck, nil
}

func (s *HealthcheckStore) DeleteHealthcheck(ID int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)

		err := b.Delete(itob(ID))
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *HealthcheckStore) GetAllHealtchecks() ([]Healthcheck, error) {
	var hcList []Healthcheck

	return hcList, s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)

		return b.ForEach(func(k, v []byte) error {
			var hc Healthcheck

			err := json.Unmarshal(v, &hc)
			if err != nil {
				return err
			}

			hcList = append(hcList, hc)
			return nil
		})
	})
}

func (s *HealthcheckStore) Debug() error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(HealthchecksBucket)

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})

		return nil
	})
}

// func (s *ReminderByDateStore) GetRemindersByDate(key time.Time) ([]ReminderByDate, error) {
// 	var reminderByDates []ReminderByDate
// 	ID := []byte(key.Format(time.RFC3339))

// 	err := s.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket(RemindersByDateBucket)
// 		v := b.Get(ID)

// 		err := json.Unmarshal(v, &reminderByDates)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return reminderByDates, nil
// }

// itob converts int to []byte
func itob(v int) []byte {
	return []byte(strconv.FormatInt(int64(v), 10))
}
