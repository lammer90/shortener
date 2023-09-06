package handlers

import (
	"github.com/lammer90/shortener/internal/app/storage"
	"github.com/lammer90/shortener/internal/app/util"
	"io"
	"net/http"
	"strings"
)

type shortenerHandler struct {
	repository storage.Repository
}

func GetShortenerHandler(repository storage.Repository) http.Handler {
	return shortenerHandler{
		repository: repository,
	}
}

func (s shortenerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		post(&res, req, s.repository)
		return
	}
	if req.Method == http.MethodGet {
		get(&res, req, s.repository)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func post(res *http.ResponseWriter, req *http.Request, repository storage.Repository) {
	body, err := io.ReadAll(req.Body)
	if err != nil || !util.CheckContentHeader(req) || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
		(*res).WriteHeader(http.StatusBadRequest)
		return
	}
	repository.Save("EwHXdJfB", string(body[:]))
	(*res).Header().Set("content-type", "text/plain")
	(*res).WriteHeader(http.StatusCreated)
	(*res).Write([]byte("http://localhost:8080/" + "EwHXdJfB"))
}

func get(res *http.ResponseWriter, req *http.Request, repository storage.Repository) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok := repository.Find(arr[len(arr)-1])
	if !ok || !util.ValidGetURL(req.URL.String()) {
		(*res).WriteHeader(http.StatusBadRequest)
		return
	}
	(*res).Header().Set("Location", address)
	(*res).WriteHeader(http.StatusTemporaryRedirect)

}
