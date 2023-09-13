package handlers

import (
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/util"
	"io"
	"net/http"
	"strings"
)

type ShortenerProvider interface {
	Save(string, string)
	Find(string) (string, bool)
}

func SaveShortUrl(provider ShortenerProvider) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil || !util.CheckContentHeader(req) || !util.ValidPostURL(req.URL.String()) || len(body) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		provider.Save("EwHXdJfB", string(body[:]))
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(config.BaseURL + "/" + "EwHXdJfB"))
	}

}

func FindByShortUrl(provider ShortenerProvider) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		arr := strings.Split(req.URL.String(), "/")
		address, ok := provider.Find(arr[len(arr)-1])
		if !ok || !util.ValidGetURL(req.URL.String()) {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", address)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
