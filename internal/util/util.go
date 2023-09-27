package util

import (
	"net/http"
	"strings"
)

func CheckContentHeader(req *http.Request) bool {
	val, ok := req.Header["Content-Type"]
	if ok {
		return strings.Contains(val[0], "text/plain")
	}
	return false
}

func ValidPostURL(url string) bool {
	return url == "/"
}

func ValidGetURL(url string) bool {
	return strings.Count(url, "/") == 1
}
