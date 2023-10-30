package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchToSave struct {
	ShortURL    string
	OriginalURL string
	UserId      string
}

type UserResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewResponse(response string) *Response {
	return &Response{
		Result: response,
	}
}

func NewBatchToSave(shortURL, originalURL, userId string) *BatchToSave {
	return &BatchToSave{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserId:      userId,
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
