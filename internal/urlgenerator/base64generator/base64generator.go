package base64generator

import (
	"encoding/base64"
)

// Base64Generator генератор значений
type Base64Generator struct{}

// GenerateURL сгенерировать уникальное значение
func (b Base64Generator) GenerateURL(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

var base64GeneratorImpl = Base64Generator{}

// New Base64Generator конструктор
func New() Base64Generator {
	return base64GeneratorImpl
}
