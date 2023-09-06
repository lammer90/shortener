package main

import (
	"github.com/lammer90/shortener/internal/app/handlers"
	"github.com/lammer90/shortener/internal/app/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(`/`, handlers.GetShortenerHandler(storage.GetStorage()))

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
