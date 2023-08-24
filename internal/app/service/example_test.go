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
	_, err := service.ShortenURL("http://example.com", "24")
	fmt.Printf("Error: %v\n", err)
	// Output:
	// Error: <nil>
}
