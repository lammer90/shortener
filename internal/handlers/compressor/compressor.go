package compressor

import (
	"compress/gzip"
	"github.com/lammer90/shortener/internal/handlers"
	"io"
	"net/http"
	"strings"
)

type Compressor struct {
	shortener handlers.Shortener
}

func New(shortener handlers.Shortener) handlers.Shortener {
	return Compressor{
		shortener,
	}
}

func (c Compressor) SaveShortURL(res http.ResponseWriter, req *http.Request) {
	compress(res, req, c.shortener.SaveShortURL)
}

func (c Compressor) FindByShortURL(res http.ResponseWriter, req *http.Request) {
	compress(res, req, c.shortener.FindByShortURL)
}

func (c Compressor) SaveShortURLApi(res http.ResponseWriter, req *http.Request) {
	compress(res, req, c.shortener.SaveShortURLApi)
}

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

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

func compress(res http.ResponseWriter, req *http.Request, f func(http.ResponseWriter, *http.Request)) {
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
	f(ow, req)
}
