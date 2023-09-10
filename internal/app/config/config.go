package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/keyjin88/shortener/internal/app/logger"
	"os"
)

// Config represents a configuration of application.
type Config struct {
	ConfigFromJSON  string `env:"CONFIG"`
	ServerAddress   string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseAddress     string `env:"BASE_URL" json:"base_url"`
	GinReleaseMode  bool   `env:"GIN_MODE"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DataBaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	SecretKey       string `env:"SECRET_KEY"`
	HTTPSEnable     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	PathToCert      string `env:"PATH_TO_CERT"`
	PathToKey       string `env:"PATH_TO_KEY"`
}

// NewConfig creates a new Config instance with the given parameters.
func NewConfig() *Config {
	return &Config{}
}

// InitConfig обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных.
func (config *Config) InitConfig() {
	flag.StringVar(&config.ConfigFromJSON, "c", "", "Path to config file")
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseAddress, "b", "http://localhost:8080", "base address for shortened url")
	flag.BoolVar(&config.GinReleaseMode, "grm", false, "gin release mode")
	flag.StringVar(&config.LogLevel, "ll", "info", "log level")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json", "path to storage")
	flag.StringVar(&config.SecretKey, "sk", "abcdefghijklmnopqrstuvwxyz123456", "secret key for cryptographic")
	flag.StringVar(&config.DataBaseDSN, "d", "", "database dsn")
	flag.BoolVar(&config.HTTPSEnable, "s", false, "https mode")
	flag.StringVar(&config.PathToCert, "ptc", "path/to/cert.pem", "Path to Certificate")
	flag.StringVar(&config.PathToKey, "ptk", "path/to/key.pem", "Path to Key")

	// Оставил для локальных тестов
	// flag.StringVar(&config.DataBaseDSN, "d", "postgres://pgadmin:postgres@localhost:5432/shortener", "database dsn")
	// парсим переданные серверу аргументы в зарегистрированные переменные.
	flag.Parse()
	if config.ConfigFromJSON != "" {
		configData, err := os.ReadFile(config.ConfigFromJSON)
		if err != nil {
			logger.Log.Errorf("failed to read JSON configuration file: %v", err)
		}

		// Применение значений из файла конфигурации к структуре Config
		if err := json.Unmarshal(configData, &config); err != nil {
			logger.Log.Errorf("failed to parse JSON configuration file: %v", err)
		}
	}
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по умолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		logger.Log.Errorf("Error parsing config: %v", err)
	}
}
