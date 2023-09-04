package handlers

import (
	"github.com/lammer90/shortener/internal/app/storage"
	"github.com/lammer90/shortener/internal/app/util"
	"io"
	"net/http"
	"strings"
)

func ToShort(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && util.CheckContentHeader(req) {
		post(&res, req)
		return
	}
	if req.Method == http.MethodGet {
		get(&res, req)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func post(res *http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		(*res).WriteHeader(http.StatusBadRequest)
		return
	}
	storage.MockStorageImpl.Save("EwHXdJfB", string(body[:]))
	(*res).Header().Set("content-type", "text/plain")
	(*res).WriteHeader(http.StatusCreated)
	(*res).Write([]byte("http://localhost:8080/" + "EwHXdJfB"))
}

func get(res *http.ResponseWriter, req *http.Request) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok := storage.MockStorageImpl.Find(arr[len(arr)-1])
	if !ok {
		(*res).WriteHeader(http.StatusBadRequest)
		return
	}
	(*res).Header().Set("Location", address)
	(*res).WriteHeader(http.StatusTemporaryRedirect)

}
