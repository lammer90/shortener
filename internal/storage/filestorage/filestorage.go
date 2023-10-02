package filestorage

import (
	"encoding/json"
	"github.com/lammer90/shortener/internal/storage"
	"os"
	"strings"
)

type fileStorage struct {
	filePath string
}

func (f fileStorage) Save(id string, value string) error {
	fileModel := newFileModel(id, id, value)
	data, err := json.Marshal(&fileModel)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	err = os.WriteFile(f.filePath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (f fileStorage) Find(id string) (string, bool, error) {
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return "", false, err
	}
	fileData := string(data)
	fileModel := fileModel{}

	for _, line := range strings.Split(fileData, "\n") {
		err := json.Unmarshal([]byte(line), &fileModel)
		if err != nil {
			return "", false, err
		}
		if fileModel.ShortURL == id {
			return fileModel.OriginalURL, true, nil
		}
	}
	return "", false, nil
}

func New(fileStoragePath string) storage.Repository {
	return fileStorage{
		filePath: fileStoragePath,
	}
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
