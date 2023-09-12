package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/shortener/config/flags"
	"github.com/lammer90/shortener/internal/app/handlers"
	"github.com/lammer90/shortener/internal/app/storage"
	"net/http"
)

func main() {
	flags.InitFlags()
	handler := handlers.GetShortenerHandler(storage.GetStorage())
	http.ListenAndServe(flags.ServAddress, ShortenerRouter(handler.Post, handler.Get))
}

func ShortenerRouter(postFunc func(http.ResponseWriter, *http.Request), getFunc func(http.ResponseWriter, *http.Request)) chi.Router {
	r := chi.NewRouter()
	r.Post("/", postFunc)
	r.Get("/{short}", getFunc)
	return r
}
