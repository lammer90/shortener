package main

import (
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/handlers/ping"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testStorage map[string]string

func (t testStorage) Save(id string, value string) error {
	t[id] = value
	return nil
}

func (t testStorage) Find(id string) (string, bool, error) {
	val, ok := t[id]
	return val, ok, nil
}

var testStorageImpl testStorage = make(map[string]string)

type mockGenerator struct{}

func (m mockGenerator) GenerateURL(data string) string {
	if data == "https://practicum.yandex.ru/" {
		return "EwHXdJfB"
	}
	return "EwHXdJfJ"
}

type pingMock struct{}

func (p pingMock) Ping() error {
	return nil
}

func TestGetShortenerHandler(t *testing.T) {
	ts := httptest.NewServer(shortenerRouter(compressor.New(logginer.New(handlers.New(testStorageImpl, mockGenerator{}, "http://localhost:8080"))), ping.New(pingMock{})))
	config.InitConfig()
	logger.InitLogger("info")
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
		storageKey   string
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
				storageKey:   "EwHXdJfB",
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
				storageKey:   "EwHXdJfB",
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
				storageKey:   "EwHXdJfB",
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
				storageKey:   "EwHXdJfB",
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
				storageKey:   "EwHXdJfB",
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
				storageKey:   "EwHXdJfB",
			},
		},
		{
			name: "negative api test POST 1",
			request: request{
				requestMethod: "POST",
				requestURL:    "/api/shorten",
				requestBody:   "",
			},
			want: want{
				code:         400,
				response:     "",
				headerName:   "",
				headerValue:  "",
				storageValue: "",
				storageKey:   "EwHXdJfJ",
			},
		},
		{
			name: "negative api test POST 2",
			request: request{
				requestMethod: "POST",
				requestURL:    "/api/shorten",
				requestBody:   "{\"test\": \"https://practicum.yandex.com\"}",
			},
			want: want{
				code:         400,
				response:     "",
				headerName:   "",
				headerValue:  "",
				storageValue: "",
				storageKey:   "EwHXdJfJ",
			},
		},
		{
			name: "positive api test POST",
			request: request{
				requestMethod: "POST",
				requestURL:    "/api/shorten",
				requestBody:   "{\"url\": \"https://practicum.yandex.com/\"}",
			},
			want: want{
				code:         201,
				response:     "{\"result\":\"http://localhost:8080/EwHXdJfJ\"}\n",
				headerName:   "Content-Type",
				headerValue:  "application/json",
				storageValue: "https://practicum.yandex.com/",
				storageKey:   "EwHXdJfJ",
			},
		},
		{
			name: "positive api test GET",
			request: request{
				requestMethod: "GET",
				requestURL:    "/EwHXdJfJ",
				requestBody:   "",
			},
			want: want{
				code:         200,
				response:     "it_does_not_matter",
				headerName:   "Location",
				headerValue:  "https://practicum.yandex.com/",
				storageValue: "https://practicum.yandex.com/",
				storageKey:   "EwHXdJfJ",
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
			assert.Equal(t, testStorageImpl[test.want.storageKey], test.want.storageValue)
			if test.want.response != "it_does_not_matter" {
				assert.Equal(t, string(resBody), test.want.response)
				assert.Equal(t, resp.Header.Get(test.want.headerName), test.want.headerValue)
			}
		})
	}
}
