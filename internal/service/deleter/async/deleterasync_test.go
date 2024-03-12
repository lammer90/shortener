package async

import (
	"testing"

	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/service/deleter"
)

type userAndValue struct {
	UserID string
	Value  string
}

type testStorage map[string]*userAndValue

func (m testStorage) Save(id string, value string, userID string) error {
	m[id] = &userAndValue{userID, value}
	return nil
}

func (m testStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		m[short.ShortURL] = &userAndValue{short.UserID, short.OriginalURL}
	}
	return nil
}

func (m testStorage) Find(id string) (string, bool, error) {
	if val, ok := m[id]; ok {
		return val.Value, ok, nil
	} else {
		return "", ok, nil
	}
}

func (m testStorage) FindByUserID(userID string) (map[string]string, error) {
	result := make(map[string]string, 0)
	for key, val := range m {
		if val.UserID == userID {
			result[key] = val.Value
		}
	}
	return result, nil
}

func (m testStorage) Delete(keys []string, userID string) error {
	return nil
}

var testStorageImpl testStorage = make(map[string]*userAndValue)

func BenchmarkDeleter_InitJobsWorker(b *testing.B) {
	delProvider3, _, _ := New(testStorageImpl, 3)

	urls := make([]*deleter.DeletingURL, 0)
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "1"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "2"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "3"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "4"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "5"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "6"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "7"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "8"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "9"})
	urls = append(urls, &deleter.DeletingURL{DeletingURL: "1234", UserID: "10"})
	message := deleter.NewDeleteMessage(urls)
	b.ResetTimer()

	b.Run("3_workers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			delProvider3.Delete(message)
		}
	})

	b.StopTimer()
	delProvider5, _, _ := New(testStorageImpl, 5)
	b.StartTimer()

	b.Run("5_workers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			delProvider5.Delete(message)
		}
	})
}
