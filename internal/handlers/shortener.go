package handlers

import (
	"encoding/json"
	"github.com/lammer90/shortener/internal/models"
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
	storage   shortenerStorageProvider
	generator urlGeneratorProvider
	baseURL   string
}

func NewShortenerHandler(storage shortenerStorageProvider, generator urlGeneratorProvider, baseURL string) Shortener {
	return ShortenerHandler{
		storage,
		generator,
		baseURL,
	}
}

func (s ShortenerHandler) SaveShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := s.generator.GenerateURL(string(body[:]))
	s.storage.Save(shortURL, string(body[:]))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(s.baseURL + "/" + shortURL))
}

func (s ShortenerHandler) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok := s.storage.Find(arr[len(arr)-1])
	if !ok || !util.ValidGetURL(req.URL.String()) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s ShortenerHandler) SaveShortURLApi(res http.ResponseWriter, req *http.Request) {
	var request models.Request
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&request)
	if err != nil || request.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := s.generator.GenerateURL(request.URL)
	s.storage.Save(shortURL, request.URL)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	if err := enc.Encode(models.NewResponse(s.baseURL + "/" + shortURL)); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
