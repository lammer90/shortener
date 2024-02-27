package async

import (
	"time"

	"github.com/lammer90/shortener/internal/service/deleter"
	"github.com/lammer90/shortener/internal/storage"
)

type batch struct {
	URLs   []string
	UserID string
}

func newBatch(urls []string, userID string) *batch {
	return &batch{urls, userID}
}

// Deleter Обработчик для удаления ссылок.
type Deleter struct {
	repository storage.Repository
	jobs       chan *deleter.DeletingURL
	batch      chan *batch
}

// New Deleter Конструктор.
func New(repository storage.Repository, workers int) (deleter.DeleteProvider, chan *deleter.DeletingURL, chan *batch) {
	del := Deleter{
		repository,
		make(chan *deleter.DeletingURL, 5),
		make(chan *batch, 5),
	}
	go del.initJobsWorker()
	for w := 1; w <= workers; w++ {
		go del.initBatchWorker()
	}
	return del, del.jobs, del.batch
}

// Delete Удалить ссылки.
func (d Deleter) Delete(message *deleter.DeleteMessage) {
	for _, url := range message.DeletingURLs {
		d.jobs <- url
	}
}

func (d Deleter) initJobsWorker() {
	ticker := time.NewTicker(2 * time.Second)

	var urls []string
	var previousUserID string

	for {
		select {
		case url := <-d.jobs:
			if url.UserID != previousUserID && previousUserID != "" {
				d.batch <- newBatch(urls, previousUserID)
				urls = nil
			}
			urls = append(urls, url.DeletingURL)
			if len(urls) == 4 {
				d.batch <- newBatch(urls, previousUserID)
				urls = nil
			}
			previousUserID = url.UserID
		case <-ticker.C:
			if len(urls) == 0 {
				continue
			}
			d.batch <- newBatch(urls, previousUserID)
			urls = nil
		}
	}
}

func (d Deleter) initBatchWorker() {
	for b := range d.batch {
		d.repository.Delete(b.URLs, b.UserID)
	}
}
