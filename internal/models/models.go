package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func NewResponse(response string) *Response {
	return &Response{
		Result: response,
	}
}
