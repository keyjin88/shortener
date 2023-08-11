package service

import (
	"fmt"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
)

func ExampleShortenService_ShortenURL() {
	service := &ShortenService{
		urlRepository: inmem.NewURLRepositoryInMem(),
		config:        &Config{},
	}

	// Вызываем функцию ShortenURL с примером
	_, err := service.ShortenURL("http://example.com", "24")

	// Проверяем результаты
	fmt.Printf("Error: %v\n", err)

	// Output:
	// Error: <nil>
}
