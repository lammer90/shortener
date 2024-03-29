package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/service/deleter"
	"github.com/lammer90/shortener/internal/storage"
	"github.com/lammer90/shortener/internal/util"
)

type urlGeneratorProvider interface {
	GenerateURL(string) string
}

// ShortenerHandler реализация ShortenerRestProvider
type ShortenerHandler struct {
	storage        storage.Repository
	generator      urlGeneratorProvider
	baseURL        string
	deleteProvider deleter.DeleteProvider
}

// New ShortenerHandler конструктор
func New(storage storage.Repository, generator urlGeneratorProvider, baseURL string, del deleter.DeleteProvider) ShortenerRestProviderWithContext {
	return ShortenerHandler{
		storage:        storage,
		generator:      generator,
		baseURL:        baseURL,
		deleteProvider: del,
	}
}

// SaveShortURL сократить оригинальную ссылку(ссылка в параметре), в ответ будет возвращена сокращенная.
func (s ShortenerHandler) SaveShortURL(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	body, err := io.ReadAll(req.Body)
	if err != nil || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := s.generator.GenerateURL(string(body[:]))
	err = s.storage.Save(shortURL, string(body[:]), ctx.UserID)
	if err != nil {
		target := new(storage.ErrConflictDB)
		if errors.As(err, &target) {
			res.WriteHeader(http.StatusConflict)
			res.Header().Set("content-type", "text/plain")
			res.Write([]byte(s.baseURL + "/" + target.ShortURL))
			return
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(s.baseURL + "/" + shortURL))
}

// FindByShortURL найти оригинальную ссылку по сокращенной.
func (s ShortenerHandler) FindByShortURL(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	arr := strings.Split(req.URL.String(), "/")
	address, ok, err := s.storage.Find(arr[len(arr)-1])
	if address == "" || err != nil || !util.ValidGetURL(req.URL.String()) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ok {
		res.WriteHeader(http.StatusGone)
		return
	}
	res.Header().Set("Location", address)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// SaveShortURLApi сократить оригинальную ссылку(ссылка в теле запроса), в ответ будет возвращена сокращенная.
func (s ShortenerHandler) SaveShortURLApi(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	var request models.Request
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&request)
	if err != nil || request.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := s.generator.GenerateURL(request.URL)
	err = s.storage.Save(shortURL, request.URL, ctx.UserID)
	if err != nil {
		target := new(storage.ErrConflictDB)
		if errors.As(err, &target) {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusConflict)
			enc := json.NewEncoder(res)
			if err := enc.Encode(models.NewResponse(s.baseURL + "/" + shortURL)); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	if err := enc.Encode(models.NewResponse(s.baseURL + "/" + shortURL)); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// SaveShortURLBatch сократить несколько ссылок батчом, в ответ будет возвращена сокращенная.
func (s ShortenerHandler) SaveShortURLBatch(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
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
		toSave = append(toSave, models.NewBatchToSave(shortURL, short.OriginalURL, ctx.UserID))
		response = append(response, models.NewBatchResponse(short.CorrelationID, s.baseURL+"/"+shortURL))
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

// FindURLByUser найти все ссылки сокращенные пользователем.
func (s ShortenerHandler) FindURLByUser(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	results, err := s.storage.FindByUserID(ctx.UserID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(results) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]*models.UserResult, 0)
	for key, val := range results {
		response = append(response, models.NewUserResult(s.baseURL+"/"+key, val))
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(response); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Delete Удалить созраненные ссылки.
func (s ShortenerHandler) Delete(res http.ResponseWriter, req *http.Request, ctx *RequestContext) {
	var urls []string
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&urls)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var deletingURL []*deleter.DeletingURL
	for _, url := range urls {
		deletingURL = append(deletingURL, deleter.NewDeletingURL(url, ctx.UserID))
	}
	go s.deleteProvider.Delete(deleter.NewDeleteMessage(deletingURL))
	res.WriteHeader(http.StatusAccepted)
}
