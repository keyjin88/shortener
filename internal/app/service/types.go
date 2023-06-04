package service

import "github.com/keyjin88/shortener/internal/app/storage"

type ShortenService struct {
	urlRepository URLRepository
	config        *Config
}

type Config struct {
	PathToStorageFile string //путь до фпйла для резервного хранения
	BaseAddress       string //base address for shortened url
}

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedURL(shortURL string) (string, error)
	FindByOriginalURL(originalURL string) (string, error)
	Save(shortURL string, url string) (storage.ShortenedURL, error)
	Close()
	SaveBatch(urls *[]storage.ShortenedURL) error
}
