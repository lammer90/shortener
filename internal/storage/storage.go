package storage

import (
	"fmt"

	"github.com/lammer90/shortener/internal/models"
)

// ErrConflictDB Тип ошибки для констрейнта в бд.
type ErrConflictDB struct {
	ShortURL string
	Err      error
}

// Error ErrConflictDB.
func (e *ErrConflictDB) Error() string {
	return fmt.Sprintf("%v : original url is %s", e.Err, e.ShortURL)
}

// Unwrap ErrConflictDB.
func (e *ErrConflictDB) Unwrap() error {
	return e.Err
}

// NewErrConflict ErrConflictDB Констуктор.
func NewErrConflict(shortURL string, err error) error {
	return &ErrConflictDB{
		ShortURL: shortURL,
		Err:      err,
	}
}

// Repository Репозиторий для работы с хранилищем ссылок.
type Repository interface {

	// Save  сохранить ссылку с параметрами: key, value, userID.
	Save(string, string, string) error

	// SaveBatch  сохранить батч с ссылками
	SaveBatch([]*models.BatchToSave) error

	// Find  Найти оригинальную ссылку по сокращенной
	Find(string) (string, bool, error)

	// Find  Найти оригинальную ссылки по владельцу
	FindByUserID(string) (map[string]string, error)

	// Delete  Удалить ссылки
	Delete([]string, string) error
}
