package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/storage/filestorage"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"net/http"
	"os"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")
	file := openFile(config.FileStoragePath)
	defer file.Close()
	http.ListenAndServe(config.ServAddress, shortenerRouter(
		compressor.New(
			logginer.New(
				handlers.New(filestorage.New(file), base64generator.New(), config.BaseURL)))))
}

func shortenerRouter(handler handlers.ShortenerRestProvider) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	r.Post("/api/shorten", handler.SaveShortURLApi)
	return r
}

func openFile(path string) *os.File {
	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return file
}
