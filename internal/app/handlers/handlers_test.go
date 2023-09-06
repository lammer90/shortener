package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
	type request struct {
		requestMethod string
		requestUrl    string
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
				requestUrl:    "/test",
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
				requestUrl:    "/",
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
				requestUrl:    "/",
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
				requestUrl:    "/EwHXdJfB/test",
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
				requestUrl:    "/EwHXdJf",
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
				requestUrl:    "/EwHXdJfB",
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
			request := httptest.NewRequest(test.request.requestMethod, test.request.requestUrl, strings.NewReader(test.request.requestBody))
			request.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			handler := GetShortenerHandler(testStorageImpl)
			handler.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, testStorageImpl["EwHXdJfB"], test.want.storageValue)
			assert.Equal(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get(test.want.headerName), test.want.headerValue)
		})
	}
}
