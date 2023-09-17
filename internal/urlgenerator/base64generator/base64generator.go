package base64generator

import (
	"encoding/base64"
)

type Base64Generator struct{}

func (b Base64Generator) GenerateURL(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

var base64GeneratorImpl = Base64Generator{}

func New() Base64Generator {
	return base64GeneratorImpl
}
