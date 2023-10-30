package storage

import (
	"fmt"
	"github.com/lammer90/shortener/internal/models"
)

type ErrConflictDB struct {
	ShortURL string
	Err      error
}

func (e *ErrConflictDB) Error() string {
	return fmt.Sprintf("%v : original url is %s", e.Err, e.ShortURL)
}

func (e *ErrConflictDB) Unwrap() error {
	return e.Err
}

func NewErrConflict(shortURL string, err error) error {
	return &ErrConflictDB{
		ShortURL: shortURL,
		Err:      err,
	}
}

type Repository interface {
	Save(string, string, string) error
	SaveBatch([]*models.BatchToSave) error
	Find(string) (string, bool, error)
	FindByUserId(string) (map[string]string, error)
}
