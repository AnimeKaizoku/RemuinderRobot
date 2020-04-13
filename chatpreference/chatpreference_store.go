package chatpreference

//go:generate mockgen -source=$GOFILE -destination=$PWD/chatpreference/mocks/${GOFILE} -package=mocks

import (
	"encoding/json"
	"errors"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

var ChatPreferencesBucket = []byte("chatpreferences")
var ErrNotFound = errors.New("chat preference not found")

type Storer interface {
	GetChatPreference(chatID int) (*ChatPreference, error)
	CreateChatPreference(*ChatPreference) error
}

type Store struct {
	db *bolt.DB
}

func NewStore(db *bolt.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateChatPreference(chatPreference *ChatPreference) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ChatPreferencesBucket)

		buf, err := json.Marshal(chatPreference)
		if err != nil {
			return err
		}

		return bucket.Put(itob(chatPreference.ChatID), buf)
	})
}

func (s *Store) GetChatPreference(chatID int) (*ChatPreference, error) {
	var chatPreference ChatPreference

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ChatPreferencesBucket)
		v := bucket.Get(itob(chatID))
		if v == nil {
			return ErrNotFound
		}

		err := json.Unmarshal(v, &chatPreference)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &chatPreference, nil
}

// itob converts int to []byte
func itob(v int) []byte {
	return []byte(strconv.FormatInt(int64(v), 10))
}
