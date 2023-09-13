package main

import (
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testStorage map[string]string

func (t testStorage) Save(id string, value string) {
	t[id] = value
}

func (t testStorage) Find(id string) (string, bool) {
	val, ok := t[id]
	return val, ok
}

var testStorageImpl testStorage = make(map[string]string)

func TestGetShortenerHandler(t *testing.T) {
	ts := httptest.NewServer(ShortenerRouter(handlers.SaveShortUrl(testStorageImpl), handlers.FindByShortUrl(testStorageImpl)))
	config.InitConfig()
	defer ts.Close()

	type request struct {
		requestMethod string
		requestURL    string
		requestBody   string
	}
	type want struct {
		code         int
		response     string
		headerName   string
		headerValue  string
		storageValue string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "negative test POST 1",
			request: request{
				requestMethod: "POST",
				requestURL:    "/test",
				requestBody:   "https://practicum.yandex.ru/",
			},
			want: want{
				code:         405,
				response:     "",
				headerName:   "",
				headerValue:  "",
				storageValue: "",
			},
		},
		{
			name: "negative test POST 2",
			request: request{
				requestMethod: "POST",
				requestURL:    "/",
				requestBody:   "",
			},
			want: want{
				code:         400,
				response:     "",
				headerName:   "",
				headerValue:  "",
				storageValue: "",
			},
		},
		{
			name: "positive test POST",
			request: request{
				requestMethod: "POST",
				requestURL:    "/",
				requestBody:   "https://practicum.yandex.ru/",
			},
			want: want{
				code:         201,
				response:     "http://localhost:8080/EwHXdJfB",
				headerName:   "Content-Type",
				headerValue:  "text/plain",
				storageValue: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "negative test GET 1",
			request: request{
				requestMethod: "GET",
				requestURL:    "/EwHXdJfB/test",
				requestBody:   "",
			},
			want: want{
				code:         404,
				response:     "404 page not found\n",
				headerName:   "",
				headerValue:  "",
				storageValue: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "negative test GET 2",
			request: request{
				requestMethod: "GET",
				requestURL:    "/EwHXdJf",
				requestBody:   "",
			},
			want: want{
				code:         400,
				response:     "",
				headerName:   "",
				headerValue:  "",
				storageValue: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "positive test GET",
			request: request{
				requestMethod: "GET",
				requestURL:    "/EwHXdJfB",
				requestBody:   "",
			},
			want: want{
				code:         200,
				response:     "it_does_not_matter",
				headerName:   "Location",
				headerValue:  "https://practicum.yandex.ru/",
				storageValue: "https://practicum.yandex.ru/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.request.requestMethod, ts.URL+test.request.requestURL, strings.NewReader(test.request.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, resp.StatusCode, test.want.code)
			defer resp.Body.Close()
			resBody, err := io.ReadAll(resp.Body)

			require.NoError(t, err)
			assert.Equal(t, testStorageImpl["EwHXdJfB"], test.want.storageValue)
			if test.want.response != "it_does_not_matter" {
				assert.Equal(t, string(resBody), test.want.response)
				assert.Equal(t, resp.Header.Get(test.want.headerName), test.want.headerValue)
			}
		})
	}
}
