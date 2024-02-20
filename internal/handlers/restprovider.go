package handlers

import "net/http"

type RequestContext struct {
	UserID string
}

type ShortenerRestProvider interface {
	SaveShortURL(http.ResponseWriter, *http.Request)
	FindByShortURL(http.ResponseWriter, *http.Request)
	SaveShortURLApi(http.ResponseWriter, *http.Request)
	SaveShortURLBatch(http.ResponseWriter, *http.Request)
	FindURLByUser(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type ShortenerRestProviderWithContext interface {
	SaveShortURL(http.ResponseWriter, *http.Request, *RequestContext)
	FindByShortURL(http.ResponseWriter, *http.Request, *RequestContext)
	SaveShortURLApi(http.ResponseWriter, *http.Request, *RequestContext)
	SaveShortURLBatch(http.ResponseWriter, *http.Request, *RequestContext)
	FindURLByUser(http.ResponseWriter, *http.Request, *RequestContext)
	Delete(http.ResponseWriter, *http.Request, *RequestContext)
}
