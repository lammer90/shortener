package config

import (
	"encoding/json"
	"flag"
	"fmt"
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

// FileConfig Флаг включения HTTPS сервера
var FileConfig string

// InitConfig Инизиализация всех параметров
func InitConfig() error {
	initFlags()
	return initEnv()
}

func initFlags() {
	flag.StringVar(&ServAddress, "a", ":8080", "Request URL")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Response URL")
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&DataSource, "d", "", "DataSource path")
	flag.StringVar(&PrivateKey, "p", "privateKey", "PrivateKey for jwt auth")
	flag.BoolVar(&EnableHTTPS, "s", false, "Enable HTTPS")
	flag.StringVar(&FileConfig, "m", "", "FileConfig path")
	flag.Parse()
}

func initEnv() error {
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

	if fileConfig := os.Getenv("CONFIG"); fileConfig != "" {
		FileConfig = fileConfig
	}

	if FileConfig != "" {
		return readConfigFromFile(FileConfig)
	}
	return nil
}

type configStruct struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

func readConfigFromFile(fileConfig string) error {
	data, err := os.ReadFile(fileConfig)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	var config configStruct
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("не удалось спарсить файл конфигурации: %w", err)
	}

	if ServAddress == "" && config.ServerAddress != "" {
		ServAddress = config.ServerAddress
	}

	if BaseURL == "" && config.BaseURL != "" {
		BaseURL = config.BaseURL
	}

	if FileStoragePath == "" && config.FileStoragePath != "" {
		FileStoragePath = config.FileStoragePath
	}

	if DataSource == "" && config.DatabaseDSN != "" {
		DataSource = config.DatabaseDSN
	}

	if !EnableHTTPS && config.EnableHTTPS {
		EnableHTTPS = config.EnableHTTPS
	}
	return nil
}
