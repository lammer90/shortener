package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/handlers/ping"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/storage/dbstorage"
	"github.com/lammer90/shortener/internal/storage/filestorage"
	"github.com/lammer90/shortener/internal/storage/inmemory"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"net/http"
	"os"
)

func main() {
	config.InitConfig()
	db := dbstorage.InitDB("pgx", config.DataSource)
	defer db.Close()
	logger.InitLogger("info")
	file := openFile(config.FileStoragePath)
	defer file.Close()
	http.ListenAndServe(config.ServAddress, shortenerRouter(
		compressor.New(
			logginer.New(
				handlers.New(filestorage.New(inmemory.New(), file), base64generator.New(), config.BaseURL))), ping.New(db)))
}

func shortenerRouter(handler handlers.ShortenerRestProvider, ping ping.Ping) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	r.Post("/api/shorten", handler.SaveShortURLApi)
	r.Get("/ping", ping.Ping)
	return r
}

func openFile(path string) *os.File {
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	return file
}
