package inmemory

import "github.com/lammer90/shortener/internal/storage"

type mockStorage map[string]string

func (m mockStorage) Save(id string, value string) {
	m[id] = value
}

func (m mockStorage) Find(id string) (string, bool) {
	val, ok := m[id]
	return val, ok
}

var mockStorageImpl mockStorage = make(map[string]string)

func New() storage.Repository {
	return mockStorageImpl
}
