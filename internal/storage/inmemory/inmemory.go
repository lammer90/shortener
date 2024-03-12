package inmemory

import (
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
)

type userAndValue struct {
	UserID string
	Value  string
}

type mockStorage map[string]*userAndValue

// Save  сохранить ссылку с параметрами: key, value, userID.
func (m mockStorage) Save(id string, value string, userID string) error {
	m[id] = &userAndValue{userID, value}
	return nil
}

// SaveBatch  сохранить батч с ссылками
func (m mockStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		m[short.ShortURL] = &userAndValue{short.UserID, short.OriginalURL}
	}
	return nil
}

// Find  Найти оригинальную ссылку по сокращенной
func (m mockStorage) Find(id string) (string, bool, error) {
	if val, ok := m[id]; ok {
		return val.Value, ok, nil
	} else {
		return "", ok, nil
	}
}

// Find  Найти оригинальную ссылки по владельцу
func (m mockStorage) FindByUserID(userID string) (map[string]string, error) {
	result := make(map[string]string, 0)
	for key, val := range m {
		if val.UserID == userID {
			result[key] = val.Value
		}
	}
	return result, nil
}

// Delete  Удалить ссылки
func (m mockStorage) Delete(keys []string, userID string) error {
	for _, key := range keys {
		m[key] = nil
	}
	return nil
}

var mockStorageImpl mockStorage = make(map[string]*userAndValue)

// New mockStorage конструктор.
func New() storage.Repository {
	return mockStorageImpl
}
