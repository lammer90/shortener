package filestorage

import (
	"encoding/json"
	"github.com/lammer90/shortener/internal/storage"
	"os"
	"strings"
)

type fileStorage struct {
	storage storage.Repository
	file    *os.File
}

func (f fileStorage) Save(id string, value string) error {
	if savedValue, ok, err := f.storage.Find(id); err != nil || !ok || savedValue != value {
		f.storage.Save(id, value)
		return saveToFile(id, value, f.file)
	}
	return nil
}

func (f fileStorage) Find(id string) (string, bool, error) {
	return f.storage.Find(id)
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
				storage.Save(fileModel.ShortURL, fileModel.OriginalURL)
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
}

func newFileModel(UUID string, ShortURL string, OriginalURL string) fileModel {
	return fileModel{
		UUID:        UUID,
		ShortURL:    ShortURL,
		OriginalURL: OriginalURL,
	}
}
