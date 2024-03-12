package inmemoryuser

import (
	"github.com/lammer90/shortener/internal/userstorage"
)

type userStorage struct {
	arr []string
}

var userStorageImpl userStorage = userStorage{make([]string, 0)}

// Save Сохранить пользователя.
func (t userStorage) Save(name string) error {
	t.arr = append(t.arr, name)
	return nil
}

// Find Найти пользователя.
func (t userStorage) Find(name string) (string, bool, error) {
	for _, val := range t.arr {
		if val == name {
			return val, true, nil
		}
	}
	return "", false, nil
}

// New userStorageImpl конструткор.
func New() userstorage.Repository {
	return userStorageImpl
}
