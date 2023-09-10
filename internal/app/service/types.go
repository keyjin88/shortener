package service

import (
	"context"
	"github.com/keyjin88/shortener/internal/app/storage"
)

const findURLErrorTemplate = "failed to find URL from repository"

// ShortenService is a service for shortening URLs.
type ShortenService struct {
	urlRepository URLRepository
	config        *Config
}

// Config is a cinfig for shortening URLs.
type Config struct {
	BaseAddress string // base address for shortened url
}

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedURL(shortURL string) (storage.ShortenedURL, error)
	FindByOriginalURL(originalURL string) (string, error)
	FindAllByUserID(userID string) ([]storage.UsersURLResponse, error)
	Save(*storage.ShortenedURL) error
	SaveBatch(urls *[]storage.ShortenedURL) error
	Close()
	Ping(ctx context.Context) error
	Delete(ids []string, userID string) error
}
