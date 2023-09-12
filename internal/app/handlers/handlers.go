package handlers

import (
	"github.com/lammer90/shortener/config/flags"
	"github.com/lammer90/shortener/internal/app/storage"
	"github.com/lammer90/shortener/internal/app/util"
	"io"
	"net/http"
	"strings"
)

type shortenerHandler struct {
	repository storage.Repository
}

func GetShortenerHandler(repository storage.Repository) shortenerHandler {
	return shortenerHandler{
		repository: repository,
	}
}

func (s shortenerHandler) Post(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || !util.CheckContentHeader(req) || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	s.repository.Save("EwHXdJfB", string(body[:]))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(flags.BaseURL + "/" + "EwHXdJfB"))
}

func (s shortenerHandler) Get(res http.ResponseWriter, req *http.Request) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok := s.repository.Find(arr[len(arr)-1])
	if !ok || !util.ValidGetURL(req.URL.String()) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)

}
