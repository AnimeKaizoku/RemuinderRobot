package backup

import (
	"encoding/json"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var BackupsBucket = []byte("backups")

type BackupStore struct {
	db *bolt.DB
}

func NewBackupStore(db *bolt.DB) *BackupStore {
	return &BackupStore{db: db}
}

// CreateUser saves u to the store. The new user ID is set on u once the data is persisted.
func (s *BackupStore) CreateBackup(r *Backup) (int, error) {
	var ID int

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BackupsBucket)

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		r.ID = int(id)
		ID = int(id)

		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		return b.Put(itob(r.ID), buf)
	})
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (s *BackupStore) GetBackup(ID int) (*Backup, error) {
	var backup Backup

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BackupsBucket)
		v := b.Get(itob(ID))

		err := json.Unmarshal(v, &backup)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &backup, nil
}

func (s *BackupStore) DeleteBackup(ID int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BackupsBucket)

		err := b.Delete(itob(ID))
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *BackupStore) Debug() error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BackupsBucket)

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
