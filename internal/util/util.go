package util

import (
	"strings"
)

// ValidPostURL валидировать post url.
func ValidPostURL(url string) bool {
	return url == "/"
}

// ValidGetURL валидировать get url.
func ValidGetURL(url string) bool {
	return strings.Count(url, "/") == 1
}
