package handlers

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lammer90/shortener/internal/config"
	"github.com/lammer90/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userAndValue struct {
	UserID string
	Value  string
}

type testStorage map[string]*userAndValue

func (m testStorage) Save(id string, value string, userID string) error {
	m[id] = &userAndValue{userID, value}
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
		if val.UserID == userID {
			result[key] = val.Value
		}
	}
	return result, nil
}

func (m testStorage) Delete(keys []string, userID string) error {
	for _, key := range keys {
		m[key] = nil
	}
	return nil
}

var testStorageImpl testStorage = make(map[string]*userAndValue)

type mockGenerator struct{}

func (m mockGenerator) GenerateURL(data string) string {
	return "EwHXdJfB"
}

var mockGeneratorImpl = mockGenerator{}

func TestGetShortenerHandler(t *testing.T) {
	config.InitConfig()
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
				code:         400,
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
				code:         400,
				response:     "",
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
				code:         307,
				response:     "",
				headerName:   "Location",
				headerValue:  "https://practicum.yandex.ru/",
				storageValue: "https://practicum.yandex.ru/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.requestMethod, test.request.requestURL, strings.NewReader(test.request.requestBody))
			request.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			handler := New(testStorageImpl, mockGeneratorImpl, "http://localhost:8080", nil)
			if test.request.requestMethod == "GET" {
				handler.FindByShortURL(w, request, &RequestContext{""})
			} else if test.request.requestMethod == "POST" {
				handler.SaveShortURL(w, request, &RequestContext{""})
			}

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if testStorageImpl["EwHXdJfB"] == nil {
				assert.Equal(t, "", test.want.storageValue)
			} else {
				assert.Equal(t, testStorageImpl["EwHXdJfB"].Value, test.want.storageValue)
			}
			assert.Equal(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get(test.want.headerName), test.want.headerValue)
		})
	}
}
