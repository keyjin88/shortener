package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseAddress   string `env:"BASE_URL"`
}

func NewConfig() *Config {
	return &Config{}
}

// InitConfig обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (config *Config) InitConfig() {
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.BaseAddress, "b", "http://localhost:8080", "base address for shortened url")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	// Пробуем распарсить переменные окружения, если их не будет, то оставляем значения по уиолчанию из флагов
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting server with configs: ServerAddress {%s}, BaseAddress {%s}", config.ServerAddress, config.BaseAddress)
}
