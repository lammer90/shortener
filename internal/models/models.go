package models

// Request Тело запросв для метода SaveShortURLApi.
type Request struct {
	URL string `json:"url"`
}

// Response Ответ для метода SaveShortURLApi.
type Response struct {
	Result string `json:"result"`
}

// BatchRequest Тело запросв для метода SaveShortURLBatch.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse Ответ для метода SaveShortURLBatch.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// BatchToSave Модельдля сохранения батчевого запроса.
type BatchToSave struct {
	ShortURL    string
	OriginalURL string
	UserID      string
}

// UserResult Ответ для метода FindURLByUser.
type UserResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewResponse(response string) *Response {
	return &Response{
		Result: response,
	}
}

func NewBatchToSave(shortURL, originalURL, userID string) *BatchToSave {
	return &BatchToSave{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
}

func NewBatchResponse(correlationID, shortURL string) *BatchResponse {
	return &BatchResponse{
		CorrelationID: correlationID,
		ShortURL:      shortURL,
	}
}

func NewUserResult(shortURL, originalURL string) *UserResult {
	return &UserResult{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
