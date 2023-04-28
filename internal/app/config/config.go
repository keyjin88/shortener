package config

import (
	"flag"
)

type Config struct {
	Flags *Flags
}

type Flags struct {
	ServerAddress string
	BaseAddr      string
}

func NewConfig() *Config {
	return &Config{Flags: &Flags{}}
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (config *Config) ParseFlags() {
	flag.StringVar(&config.Flags.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.Flags.BaseAddr, "b", "http://localhost:8080/", "base address for shortened url")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
