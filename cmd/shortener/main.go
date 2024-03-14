package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware"
	"github.com/lammer90/shortener/internal/handlers/middleware/auth"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/handlers/ping"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/service/deleter/async"
	"github.com/lammer90/shortener/internal/storage"
	"github.com/lammer90/shortener/internal/storage/dbstorage"
	"github.com/lammer90/shortener/internal/storage/filestorage"
	"github.com/lammer90/shortener/internal/storage/inmemory"
	"github.com/lammer90/shortener/internal/urlgenerator/base64generator"
	"github.com/lammer90/shortener/internal/userstorage"
	"github.com/lammer90/shortener/internal/userstorage/dbuserstorage"
	"github.com/lammer90/shortener/internal/userstorage/inmemoryuser"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("version=%s, time=%s, commit=%s\n", buildVersion, buildDate, buildCommit)

	config.InitConfig()
	logger.InitLogger("info")
	st, userSt, cl, db := getActualStorage()
	delProvider, ch1, ch2 := async.New(st, 3)
	defer cl.Close()
	defer close(ch1)
	defer close(ch2)
	r := shortenerRouter(
		auth.New(userSt, config.PrivateKey,
			compressor.New(
				logginer.New(
					handlers.New(st, base64generator.New(), config.BaseURL, delProvider)))),
		ping.New(db))
	http.ListenAndServe(config.ServAddress, r)
}

func shortenerRouter(handler handlers.ShortenerRestProvider, ping ping.Ping) chi.Router {
	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Post("/", handler.SaveShortURL)
	r.Get("/{short}", handler.FindByShortURL)
	r.Post("/api/shorten", handler.SaveShortURLApi)
	r.Post("/api/shorten/batch", handler.SaveShortURLBatch)
	r.Get("/ping", ping.Ping)
	r.Get("/api/user/urls", handler.FindURLByUser)
	r.Delete("/api/user/urls", handler.Delete)
	return r
}

func getActualStorage() (storage.Repository, userstorage.Repository, io.Closer, *sql.DB) {
	if config.DataSource != "" {
		db := initDB("pgx", config.DataSource)
		return dbstorage.New(db), dbuserstorage.New(db), db, db
	} else {
		file := openFile(config.FileStoragePath)
		return filestorage.New(inmemory.New(), file), inmemoryuser.New(), file, nil
	}
}

func openFile(path string) *os.File {
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	return file
}

func initDB(driverName, dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	return db
}
