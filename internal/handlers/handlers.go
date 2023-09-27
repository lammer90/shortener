package handlers

import (
	"github.com/lammer90/shortener/internal/util"
	"io"
	"net/http"
	"strings"
)

type shortenerStorageProvider interface {
	Save(string, string)
	Find(string) (string, bool)
}

type urlGeneratorProvider interface {
	GenerateURL(string) string
}

type ShortenerHandler struct {
	shortenerStorageProvider
	urlGeneratorProvider
	baseURL string
}

func New(storage shortenerStorageProvider, generator urlGeneratorProvider, baseURL string) ShortenerHandler {
	return ShortenerHandler{
		storage,
		generator,
		baseURL,
	}
}

func (s ShortenerHandler) SaveShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || !util.CheckContentHeader(req) || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := s.GenerateURL(string(body[:]))
	s.Save(shortURL, string(body[:]))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(s.baseURL + "/" + shortURL))
}

func (s ShortenerHandler) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok := s.Find(arr[len(arr)-1])
	if !ok || !util.ValidGetURL(req.URL.String()) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
