package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/lammer90/shortener/internal/handlers"
)

// Compressor фильтр
type Compressor struct {
	shortener handlers.ShortenerRestProviderWithContext
}

// New Compressor констуктор
func New(shortener handlers.ShortenerRestProviderWithContext) handlers.ShortenerRestProviderWithContext {
	return Compressor{
		shortener,
	}
}

// SaveShortURL сократить оригинальную ссылку(ссылка в параметре), в ответ будет возвращена сокращенная.
func (c Compressor) SaveShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.SaveShortURL)
}

// FindByShortURL найти оригинальную ссылку по сокращенной.
func (c Compressor) FindByShortURL(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.FindByShortURL)
}

// SaveShortURLApi сократить оригинальную ссылку(ссылка в теле запроса), в ответ будет возвращена сокращенная.
func (c Compressor) SaveShortURLApi(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.SaveShortURLApi)
}

// SaveShortURLBatch сократить несколько ссылок батчом, в ответ будет возвращена сокращенная.
func (c Compressor) SaveShortURLBatch(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.SaveShortURLBatch)
}

// FindURLByUser найти все ссылки сокращенные пользователем.
func (c Compressor) FindURLByUser(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.FindURLByUser)
}

// Delete Удалить созраненные ссылки.
func (c Compressor) Delete(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext) {
	compress(res, req, ctx, c.shortener.Delete)
}

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter compressWriter Конструтор
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

// Write Записать
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Close Закрыть
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	*gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		Reader: zr,
	}, nil
}

func compress(res http.ResponseWriter, req *http.Request, ctx *handlers.RequestContext, f func(http.ResponseWriter, *http.Request, *handlers.RequestContext)) {
	ow := res
	acceptEncoding := req.Header.Get("Accept-Encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		cw := newCompressWriter(res)
		ow = cw
		ow.Header().Set("Content-Encoding", "gzip")
		defer cw.Close()
	}

	contentEncoding := req.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		cr, err := newCompressReader(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Body = cr
		defer cr.Close()
	}
	f(ow, req, ctx)
}
