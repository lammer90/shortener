package inmemory

import (
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
)

type userAndValue struct {
	UserId string
	Value  string
}

type mockStorage map[string]*userAndValue

func (m mockStorage) Save(id string, value string, userID string) error {
	m[id] = &userAndValue{userID, value}
	return nil
}

func (m mockStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		m[short.ShortURL] = &userAndValue{short.UserID, short.OriginalURL}
	}
	return nil
}

func (m mockStorage) Find(id string) (string, bool, error) {
	if val, ok := m[id]; ok {
		return val.Value, ok, nil
	} else {
		return "", ok, nil
	}
}

func (m mockStorage) FindByUserID(userID string) (map[string]string, error) {
	result := make(map[string]string, 0)
	for key, val := range m {
		if val.UserId == userID {
			result[key] = val.Value
		}
	}
	return result, nil
}

var mockStorageImpl mockStorage = make(map[string]*userAndValue)

func New() storage.Repository {
	return mockStorageImpl
}
