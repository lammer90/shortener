package handlers

import "net/http"

// RequestContext Контекст авторизованного пользователя.
type RequestContext struct {
	UserID string
}

// ShortenerRestProvider Обработчик без контекста авторизованного пользовател
type ShortenerRestProvider interface {

	// SaveShortURL сократить оригинальную ссылку, в ответ будет возвращена сокращенная.
	SaveShortURL(http.ResponseWriter, *http.Request)

	// FindByShortURL найти оригинальную ссылку по сокращенной.
	FindByShortURL(http.ResponseWriter, *http.Request)

	// SaveShortURLApi сократить оригинальную ссылку(ссылка в теле запроса), в ответ будет возвращена сокращенная.
	SaveShortURLApi(http.ResponseWriter, *http.Request)

	// SaveShortURLBatch сократить несколько ссылок батчом, в ответ будет возвращена сокращенная.
	SaveShortURLBatch(http.ResponseWriter, *http.Request)

	// FindURLByUser найти все ссылки сокращенные пользователем.
	FindURLByUser(http.ResponseWriter, *http.Request)

	// Delete Удалить созраненные ссылки.
	Delete(http.ResponseWriter, *http.Request)
}

// ShortenerRestProviderWithContext Обработчик принимает контекст авторизованного пользователя
type ShortenerRestProviderWithContext interface {

	// SaveShortURL сократить оригинальную ссылку(ссылка в параметре), в ответ будет возвращена сокращенная.
	SaveShortURL(http.ResponseWriter, *http.Request, *RequestContext)

	// FindByShortURL найти оригинальную ссылку по сокращенной.
	FindByShortURL(http.ResponseWriter, *http.Request, *RequestContext)

	// SaveShortURLApi сократить оригинальную ссылку(ссылка в теле запроса), в ответ будет возвращена сокращенная.
	SaveShortURLApi(http.ResponseWriter, *http.Request, *RequestContext)

	// SaveShortURLBatch сократить несколько ссылок батчом, в ответ будет возвращена сокращенная.
	SaveShortURLBatch(http.ResponseWriter, *http.Request, *RequestContext)

	// FindURLByUser найти все ссылки сокращенные пользователем.
	FindURLByUser(http.ResponseWriter, *http.Request, *RequestContext)

	// Delete Удалить созраненные ссылки.
	Delete(http.ResponseWriter, *http.Request, *RequestContext)
}
