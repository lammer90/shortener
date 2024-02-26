package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"
)

func ExampleSaveShortURL() {
	request := httptest.NewRequest("POST", "/", strings.NewReader("https://practicum.yandex.ru/"))
	request.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	handler := New(testStorageImpl, mockGeneratorImpl, "http://localhost:8080", nil)
	handler.SaveShortURL(w, request, &RequestContext{""})
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	// Output:
	// 201
}
