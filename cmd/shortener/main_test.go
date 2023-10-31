package main

import (
	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/handlers"
	"github.com/lammer90/shortener/internal/handlers/middleware/auth"
	"github.com/lammer90/shortener/internal/handlers/middleware/compressor"
	"github.com/lammer90/shortener/internal/handlers/middleware/logginer"
	"github.com/lammer90/shortener/internal/handlers/ping"
	"github.com/lammer90/shortener/internal/logger"
	"github.com/lammer90/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type userAndValue struct {
	UserId string
	Value  string
}

type testStorage map[string]*userAndValue

func (m testStorage) Save(id string, value string, userId string) error {
	m[id] = &userAndValue{userId, value}
	return nil
}

func (m testStorage) SaveBatch(shorts []*models.BatchToSave) error {
	for _, short := range shorts {
		m[short.ShortURL] = &userAndValue{short.UserID, short.OriginalURL}
	}
	return nil
}

func (m testStorage) Find(id string) (string, bool, error) {
	if val, ok := m[id]; ok {
		return val.Value, ok, nil
	} else {
		return "", ok, nil
	}
}

func (m testStorage) FindByUserID(userID string) (map[string]string, error) {
	result := make(map[string]string, 0)
	for key, val := range m {
		if val.UserId == userID {
			result[key] = val.Value
		}
	}
	return result, nil
}

var testStorageImpl testStorage = make(map[string]*userAndValue)

type mockGenerator struct{}

type testUserStorage struct {
	arr []string
}

var testUserStorageImpl testUserStorage = testUserStorage{make([]string, 0)}

func (t testUserStorage) Save(name string) error {
	t.arr = append(t.arr, name)
	return nil
}

func (t testUserStorage) Find(name string) (string, bool, error) {
	for _, val := range t.arr {
		if val == name {
			return val, true, nil
		}
	}
	return "", false, nil
}

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
	ts := httptest.NewServer(shortenerRouter(auth.New(testUserStorageImpl, "test", compressor.New(logginer.New(handlers.New(testStorageImpl, mockGenerator{}, "http://localhost:8080")))), ping.New(pingMock{})))
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
			if testStorageImpl[test.want.storageKey] == nil {
				assert.Equal(t, "", test.want.storageValue)
			} else {
				assert.Equal(t, testStorageImpl[test.want.storageKey].Value, test.want.storageValue)
			}
			if test.want.response != "it_does_not_matter" {
				assert.Equal(t, string(resBody), test.want.response)
				assert.Equal(t, resp.Header.Get(test.want.headerName), test.want.headerValue)
			}
		})
	}
}
