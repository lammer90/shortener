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
}

func NewResponse(response string) *Response {
	return &Response{
		Result: response,
	}
}

func NewBatchToSave(ShortURL, OriginalURL string) *BatchToSave {
	return &BatchToSave{
		ShortURL:    ShortURL,
		OriginalURL: OriginalURL,
	}
}

func NewBatchResponse(correlationID, shortURL string) *BatchResponse {
	return &BatchResponse{
		CorrelationID: correlationID,
		ShortURL:      shortURL,
	}
}
