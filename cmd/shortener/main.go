package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/storage/inmemory"
	"net/http"
)

func main() {
	config.InitConfig()
	handler := handlers.New(inmemory.New())
	http.ListenAndServe(config.ServAddress, shortenerRouter(handler.SaveShortURL, handler.FindByShortURL))
}

func shortenerRouter(postFunc func(http.ResponseWriter, *http.Request), getFunc func(http.ResponseWriter, *http.Request)) chi.Router {
	r := chi.NewRouter()
	r.Post("/", postFunc)
	r.Get("/{short}", getFunc)
	return r
}
