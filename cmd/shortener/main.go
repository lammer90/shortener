package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/compressor"
	"github.com/lammer90/shortener/internal/handlers/logginer"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/storage/filestorage"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"net/http"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")
	http.ListenAndServe(config.ServAddress, shortenerRouter(
		compressor.New(
			logginer.New(
				handlers.NewShortenerHandler(filestorage.New(), base64generator.New(), config.BaseURL)))))
}

func shortenerRouter(handler handlers.Shortener) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	r.Post("/api/shorten", handler.SaveShortURLApi)
	return r
}
