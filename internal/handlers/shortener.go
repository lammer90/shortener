package handlers

import (
	"encoding/json"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/util"
	"io"
	"net/http"
	"strings"
)

type shortenerStorageProvider interface {
	Save(string, string) error
	SaveBatch([]*models.BatchToSave) error
	Find(string) (string, bool, error)
}

type urlGeneratorProvider interface {
	GenerateURL(string) string
}

type ShortenerHandler struct {
	storage   shortenerStorageProvider
	generator urlGeneratorProvider
	baseURL   string
}

func New(storage shortenerStorageProvider, generator urlGeneratorProvider, baseURL string) ShortenerRestProvider {
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
	err = s.storage.Save(shortURL, string(body[:]))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(s.baseURL + "/" + shortURL))
}

func (s ShortenerHandler) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	logger.Log.Info("< FindByShortURL")
	arr := strings.Split(req.URL.String(), "/")
	address, ok, err := s.storage.Find(arr[len(arr)-1])
	if !ok || err != nil || !util.ValidGetURL(req.URL.String()) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
	logger.Log.Info("> FindByShortURL")
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

func (s ShortenerHandler) SaveShortURLBatch(res http.ResponseWriter, req *http.Request) {
	shorts := make([]models.BatchRequest, 0)
	toSave := make([]*models.BatchToSave, 0)
	response := make([]*models.BatchResponse, 0)
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&shorts)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, short := range shorts {
		shortURL := s.generator.GenerateURL(short.OriginalURL)
		toSave = append(toSave, models.NewBatchToSave(shortURL, short.OriginalURL))
		response = append(response, models.NewBatchResponse(short.CorrelationID, shortURL))
	}
	s.storage.SaveBatch(toSave)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	if err := enc.Encode(response); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
