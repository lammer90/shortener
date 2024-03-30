package config

import (
	"flag"
	"os"
)

// ServAddress Адрес старта веб-сервера
var ServAddress string

// BaseURL Адрес для запросос в короткой сслыке
var BaseURL string

// FileStoragePath Папка для хранения данных по ссылкам
var FileStoragePath string

// DataSource Строка подключения к бд
var DataSource string

// PrivateKey Приватный ключ для подписи jwt токена
var PrivateKey string

// EnableHTTPS Флаг включения HTTPS сервера
var EnableHTTPS bool

// InitConfig Инизиализация всех параметров
func InitConfig() {
	initFlags()
	initEnv()
}

func initFlags() {
	flag.StringVar(&ServAddress, "a", ":8080", "Request URL")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Response URL")
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&DataSource, "d", "", "DataSource path")
	flag.StringVar(&PrivateKey, "p", "privateKey", "PrivateKey for jwt auth")
	flag.BoolVar(&EnableHTTPS, "s", false, "Enable HTTPS")
	flag.Parse()
}

func initEnv() {
	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		ServAddress = envServAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		BaseURL = envBaseURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}

	if envDataSource := os.Getenv("DATABASE_DSN"); envDataSource != "" {
		DataSource = envDataSource
	}

	if privateKey := os.Getenv("PRIVATE_KEY"); privateKey != "" {
		PrivateKey = privateKey
	}

	if enableHTTPS := os.Getenv(" ENABLE_HTTPS"); enableHTTPS != "" {
		EnableHTTPS = true
	}
}
