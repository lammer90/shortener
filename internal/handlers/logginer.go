package handlers

import (
	"github.com/lammer90/shortener/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoggingHandler struct {
	shortener Shortener
}

func NewLoggingHandler(shortener Shortener) Shortener {
	return LoggingHandler{
		shortener,
	}
}

func (l LoggingHandler) SaveShortURL(res http.ResponseWriter, req *http.Request) {
	l.log(res, req, l.shortener.SaveShortURL)
}

func (l LoggingHandler) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	l.log(res, req, l.shortener.FindByShortURL)
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

func (l LoggingHandler) log(res http.ResponseWriter, req *http.Request, f func(http.ResponseWriter, *http.Request)) {
	start := time.Now()

	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := loggingResponseWriter{
		ResponseWriter: res,
		responseData:   responseData,
	}

	f(&lw, req)

	duration := time.Since(start)

	logger.Log.Info("Server received new request",
		zap.String("uri", req.RequestURI),
		zap.String("method", req.Method),
		zap.Int("status", responseData.status),
		zap.Int("duration", int(duration.Milliseconds())),
		zap.Int("size", responseData.size),
	)
}
