package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseAddress     string `env:"BASE_URL"`
	GinReleaseMode  bool   `env:"GIN_MODE"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"SECRET_KEY"`
}

func NewConfig() *Config {
	return &Config{}
}

// InitConfig обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (config *Config) InitConfig() {
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseAddress, "b", "http://localhost:8080", "base address for shortened url")
	flag.BoolVar(&config.GinReleaseMode, "grm", false, "gin release mode")
	flag.StringVar(&config.LogLevel, "ll", "info", "log level")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json", "path to storage")
	flag.StringVar(&config.SecretKey, "sk", "abcdefghijklmnopqrstuvwxyz123456", "secret key for cryptographic")
	//flag.StringVar(&config.DataBaseDSN, "d", "", "database dsn")
	//Оставил для локальных тестов
	flag.StringVar(&config.DataBaseDSN, "d", "postgres://pgadmin:postgres@localhost:5432/shortener", "database dsn")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по умолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
}
