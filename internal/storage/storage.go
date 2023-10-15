package storage

import "github.com/lammer90/shortener/internal/models"

type Repository interface {
	Save(string, string) error
	SaveBatch([]*models.BatchToSave) error
	Find(string) (string, bool, error)
}
