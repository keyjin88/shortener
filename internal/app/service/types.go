package service

import (
	"context"
	"github.com/keyjin88/shortener/internal/app/storage"
)

type ShortenService struct {
	urlRepository URLRepository
	config        *Config
}

type Config struct {
	BaseAddress string //base address for shortened url
}

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedURL(shortURL string) (string, error)
	FindByOriginalURL(originalURL string) (string, error)
	Save(*storage.ShortenedURL) error
	SaveBatch(urls *[]storage.ShortenedURL) error
	Close()
	Ping(ctx context.Context) error
}
