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
	return &DeleteMessage{urls}
}

type DeleteProvider interface {
	Delete(*DeleteMessage)
}
