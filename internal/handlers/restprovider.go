package handlers

import "net/http"

type ShortenerRestProvider interface {
	SaveShortURL(http.ResponseWriter, *http.Request)
	FindByShortURL(http.ResponseWriter, *http.Request)
	SaveShortURLApi(http.ResponseWriter, *http.Request)
	SaveShortURLBatch(http.ResponseWriter, *http.Request)
}
