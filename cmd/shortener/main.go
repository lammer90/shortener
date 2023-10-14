package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/handlers/ping"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/storage"
	"github.com/lammer90/shortener/internal/storage/dbstorage"
	"github.com/lammer90/shortener/internal/storage/filestorage"
	"github.com/lammer90/shortener/internal/storage/inmemory"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"io"
	"net/http"
	"os"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")
	st, cl, db := getActualStorage()
	defer cl.Close()
	http.ListenAndServe(config.ServAddress, shortenerRouter(
		compressor.New(
			logginer.New(
				handlers.New(st, base64generator.New(), config.BaseURL))),
		ping.New(db)))
}

func shortenerRouter(handler handlers.ShortenerRestProvider, ping ping.Ping) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	r.Post("/api/shorten", handler.SaveShortURLApi)
	r.Get("/ping", ping.Ping)
	return r
}

func getActualStorage() (storage.Repository, io.Closer, *sql.DB) {
	if config.DataSource != "" {
		db := InitDB("pgx", config.DataSource)
		return dbstorage.New(db), db, db
	} else {
		file := openFile(config.FileStoragePath)
		return filestorage.New(inmemory.New(), file), file, nil
	}
}

func openFile(path string) *os.File {
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	return file
}

func InitDB(driverName, dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	return db
}
