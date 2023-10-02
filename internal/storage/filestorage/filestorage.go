package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/lammer90/shortener/internal/storage"
	"os"
)

type fileStorage struct {
	file *os.File
}

func (f fileStorage) Save(id string, value string) error {
	fileModel := newFileModel(id, id, value)
	data, err := json.Marshal(&fileModel)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = f.file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (f fileStorage) Find(id string) (string, bool, error) {
	scanner := bufio.NewScanner(f.file)
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

func New(file *os.File) storage.Repository {
	return fileStorage{
		file: file,
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
