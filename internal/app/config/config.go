package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"log"
)

type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS"`
	BaseAddress    string `env:"BASE_URL"`
	GinReleaseMode bool   `env:"GIN_MODE"`
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
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по уиолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Infof("Starting server with configs: ServerAddress {%s}, BaseAddress {%s}\n", config.ServerAddress,
		config.BaseAddress)
}
