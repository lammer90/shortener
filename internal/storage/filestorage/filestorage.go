package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/lammer90/shortener/internal/storage"
	"os"
)

type fileStorage struct {
	filePath string
}

func (f fileStorage) Save(id string, value string) error {
	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
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
	file.Close()
	return nil
}

func (f fileStorage) Find(id string) (string, bool, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return "", false, err
	}
	scanner := bufio.NewScanner(file)
	fileModel := fileModel{}

	for scanner.Scan() {
		data := scanner.Bytes()
		err := json.Unmarshal(data, &fileModel)
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
