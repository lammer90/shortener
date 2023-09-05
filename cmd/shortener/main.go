package main

import (
	"github.com/lammer90/shortener/internal/app/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(`/`, handlers.GetShortenerHandler())

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
