package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/acme/autocert"

	"github.com/lammer90/shortener/internal/logger"
	"go.uber.org/zap"

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
	logger.InitLogger("info")
	logger.Log.Info("Starting shortener app",
		zap.String("version", buildVersion),
		zap.String("time", buildDate),
		zap.String("commit", buildCommit))

	err := config.InitConfig()
	if err != nil {
		logger.Log.Fatal("Ошибка загрузки конфигурации: " + err.Error())
	}

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

	var server *http.Server
	if config.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}
		server = &http.Server{
			Addr:      ":443",
			Handler:   r,
			TLSConfig: manager.TLSConfig(),
		}
	} else {
		server = &http.Server{
			Addr:    config.ServAddress,
			Handler: r,
		}
	}

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Log.Error("HTTP server Shutdown: %v", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if config.EnableHTTPS {
		server.ListenAndServeTLS("", "")
	} else {
		server.ListenAndServe()
	}

	<-idleConnsClosed
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
