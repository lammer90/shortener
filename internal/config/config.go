package config

import (
	"flag"
	"os"
)

var ServAddress string
var BaseURL string
var FileStoragePath string
var DataSource string
var PrivateKey string

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
}
