package async

import (
	"github.com/lammer90/shortener/internal/service/deleter"
	"github.com/lammer90/shortener/internal/storage"
	"time"
)

type Batch struct {
	URLs   []string
	UserID string
}

func newBatch(urls []string, userID string) *Batch {
	return &Batch{urls, userID}
}

type Deleter struct {
	repository storage.Repository
	jobs       chan *deleter.DeletingURL
	batch      chan *Batch
}

func New(repository storage.Repository) (deleter.DeleteProvider, chan *deleter.DeletingURL, chan *Batch) {
	del := Deleter{
		repository,
		make(chan *deleter.DeletingURL, 5),
		make(chan *Batch, 5),
	}
	go del.InitJobsWorker()
	for w := 1; w <= 3; w++ {
		go del.InitBatchWorker()
	}
	return del, del.jobs, del.batch
}

func (d Deleter) Delete(message *deleter.DeleteMessage) {
	for _, url := range message.DeletingURLs {
		d.jobs <- url
	}
}

func (d Deleter) InitJobsWorker() {
	ticker := time.NewTicker(2 * time.Second)

	var urls []string
	var previousUserID string

	for {
		select {
		case url := <-d.jobs:
			urls = append(urls, url.DeletingURL)
			if len(urls) == 4 || (url.UserID != previousUserID && previousUserID != "") {
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

func (d Deleter) InitBatchWorker() {
	for b := range d.batch {
		d.repository.Delete(b.URLs, b.UserID)
	}
}
