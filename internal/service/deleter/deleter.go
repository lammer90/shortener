package deleter

type DeletingURL struct {
	DeletingURL string
	UserID      string
}

func NewDeletingURL(deletingURL, userID string) *DeletingURL {
	return &DeletingURL{deletingURL, userID}
}

type DeleteMessage struct {
	DeletingURLs []*DeletingURL
}

func NewDeleteMessage(urls []*DeletingURL) *DeleteMessage {
	var result []*DeletingURL
	for _, url := range urls {
		result = append(result, url)
	}
	return &DeleteMessage{result}
}

type DeleteProvider interface {
	Delete(*DeleteMessage)
}
