package util

import (
	"net/http"
)

func CheckContentHeader(req *http.Request) bool {
	val, ok := req.Header["Content-Type"]
	if ok {
		return "text/plain" == val[0]
	}
	return false
}
