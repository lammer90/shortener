package inmemory

import (
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
)

type mockStorage map[string]string

func (m mockStorage) Save(id string, value string) error {
	m[id] = value
	return nil
}

func (m mockStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		m[short.ShortURL] = short.OriginalURL
	}
	return nil
}

func (m mockStorage) Find(id string) (string, bool, error) {
	val, ok := m[id]
	return val, ok, nil
}

var mockStorageImpl mockStorage = make(map[string]string)

func New() storage.Repository {
	return mockStorageImpl
}
