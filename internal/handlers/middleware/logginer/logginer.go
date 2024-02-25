package logginer

import (
	"net/http"
	"time"

	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/logger"
	"go.uber.org/zap"
)

type Logging struct {
	shortener handlers.ShortenerRestProviderWithContext
}

func New(shortener handlers.ShortenerRestProviderWithContext) handlers.ShortenerRestProviderWithContext {
	return Logging{
		shortener,
	}
}

func (l Logging) SaveShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURL)
}

func (l Logging) FindByShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.FindByShortURL)
}

func (l Logging) SaveShortURLApi(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURLApi)
}

func (l Logging) SaveShortURLBatch(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.SaveShortURLBatch)
}

func (l Logging) FindURLByUser(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	log(res, req, ctx, l.shortener.FindURLByUser)
}

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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

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
