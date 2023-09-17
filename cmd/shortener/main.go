package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/storage/inmemory"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"net/http"
)

func main() {
	config.InitConfig()
	http.ListenAndServe(config.ServAddress, shortenerRouter(handlers.New(inmemory.New(), base64generator.New())))
}

func shortenerRouter(handler handlers.ShortenerHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	return r
}
