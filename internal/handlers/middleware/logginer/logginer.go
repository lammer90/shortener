package logginer

import (
	"net/http"
	"time"

	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/logger"
	"go.uber.org/zap"
)

// Logging фильтр
type Logging struct {
	shortener handlers.ShortenerRestProviderWithContext
}

// New Logging констуктор
func New(shortener handlers.ShortenerRestProviderWithContext) handlers.ShortenerRestProviderWithContext {
	return Logging{
		shortener,
	}
}

// SaveShortURL сократить оригинальную ссылку(ссылка в параметре), в ответ будет возвращена сокращенная.
func (l Logging) SaveShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURL)
}

// FindByShortURL найти оригинальную ссылку по сокращенной.
func (l Logging) FindByShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.FindByShortURL)
}

// SaveShortURLApi сократить оригинальную ссылку(ссылка в теле запроса), в ответ будет возвращена сокращенная.
func (l Logging) SaveShortURLApi(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURLApi)
}

// SaveShortURLBatch сократить несколько ссылок батчом, в ответ будет возвращена сокращенная.
func (l Logging) SaveShortURLBatch(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURLBatch)
}

// FindURLByUser найти все ссылки сокращенные пользователем.
func (l Logging) FindURLByUser(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.FindURLByUser)
}

// Delete Удалить созраненные ссылки.
func (l Logging) Delete(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.Delete)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write Записать
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader Записать statusCode
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func log(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext, f func(http.ResponseWriter, *http.Request, *handlers.RequestContext)) {
	start := time.Now()

	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := loggingResponseWriter{
		ResponseWriter: res,
		responseData:   responseData,
	}

	f(&lw, req, ctx)

	duration := time.Since(start)

	logger.Log.Info("Server received new request",
		zap.String("uri", req.RequestURI),
		zap.String("method", req.Method),
		zap.Int("status", responseData.status),
		zap.Int("duration", int(duration.Milliseconds())),
		zap.Int("size", responseData.size),
	)
}
