package filestorage

import (
	"encoding/json"
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
	"os"
	"strings"
)

type fileStorage struct {
	storage storage.Repository
	file    *os.File
}

func (f fileStorage) Save(id string, value string, userId string) error {
	if savedValue, ok, err := f.storage.Find(id); err != nil || !ok || savedValue != value {
		if err := f.storage.Save(id, value, userId); err != nil {
			return err
		}
		return saveToFile(id, value, f.file)
	}
	return nil
}

func (f fileStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		if savedValue, ok, err := f.storage.Find(short.ShortURL); err != nil || !ok || savedValue != short.OriginalURL {
			if err := f.storage.Save(short.ShortURL, short.OriginalURL, short.UserID); err != nil {
				return err
			}
			return saveToFile(short.ShortURL, short.OriginalURL, f.file)
		}
	}
	return nil
}

func (f fileStorage) Find(id string) (string, bool, error) {
	return f.storage.Find(id)
}

func (f fileStorage) FindByUserID(userID string) (map[string]string, error) {
	return f.storage.FindByUserID(userID)
}

func New(storage storage.Repository, file *os.File) storage.Repository {
	initStorage(storage, file)
	return fileStorage{
		storage: storage,
		file:    file,
	}
}

func initStorage(storage storage.Repository, file *os.File) {
	data, err := os.ReadFile(file.Name())
	if err == nil {
		fileData := string(data)
		fileModel := fileModel{}

		for _, line := range strings.Split(fileData, "\n") {
			err := json.Unmarshal([]byte(line), &fileModel)
			if err == nil {
				storage.Save(fileModel.ShortURL, fileModel.OriginalURL, fileModel.UserId)
			}
		}
	}
}

func saveToFile(id string, value string, file *os.File) error {
	fileModel := newFileModel(id, id, value)
	data, err := json.Marshal(&fileModel)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

type fileModel struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserId      string `json:"user_id"`
}

func newFileModel(UUID string, ShortURL string, OriginalURL string) fileModel {
	return fileModel{
		UUID:        UUID,
		ShortURL:    ShortURL,
		OriginalURL: OriginalURL,
	}
}
