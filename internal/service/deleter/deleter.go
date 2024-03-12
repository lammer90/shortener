package deleter

// DeletingURL Ссылка для удаления.
type DeletingURL struct {
	DeletingURL string
	UserID      string
}

// NewDeletingURL DeletingURL конструктор.
func NewDeletingURL(deletingURL, userID string) *DeletingURL {
	return &DeletingURL{deletingURL, userID}
}

// DeleteMessage Список ссылок для удаления.
type DeleteMessage struct {
	DeletingURLs []*DeletingURL
}

// NewDeleteMessage DeleteMessage конструткор.
func NewDeleteMessage(urls []*DeletingURL) *DeleteMessage {
	return &DeleteMessage{urls}
}

// DeleteProvider Обработчик для удаления ссылок.
type DeleteProvider interface {

	// Delete Удалить ссылки.
	Delete(*DeleteMessage)
}
