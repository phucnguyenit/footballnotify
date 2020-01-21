package store

import "github.com/syndtr/goleveldb/leveldb"

// Store ...
type Store struct {
	db *leveldb.DB
}

// New new store
func New(dbPath string) *Store {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	return &Store{db: db}
}

// SetFootballAPIKey ...
func (s *Store) SetFootballAPIKey(apiKey string) error {
	return s.db.Put(getFootballAPIKey(), []byte(apiKey), nil)
}

// GetFootballAPIKey ...
func (s *Store) GetFootballAPIKey() string {
	b, err := s.db.Get(getFootballAPIKey(), nil)
	if err != nil {
		return ""
	}
	return string(b)
}

func getFootballAPIKey() []byte {
	return []byte("football_api_key")
}
