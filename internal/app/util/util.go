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

func ValidPostUrl(url string) bool {
	return url == "/"
}

func ValidGetUrl(url string) bool {
	return strings.Count(url, "/") == 1
}
