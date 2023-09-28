package handlers

import "net/http"

type Shortener interface {
	SaveShortURL(http.ResponseWriter, *http.Request)
	FindByShortURL(http.ResponseWriter, *http.Request)
	SaveShortURLApi(http.ResponseWriter, *http.Request)
}
