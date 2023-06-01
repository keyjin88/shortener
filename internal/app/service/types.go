package service

import "github.com/keyjin88/shortener/internal/app/storage"

type Config struct {
	PathToStorageFile string //путь до фпйла для резервного хранения
}

//go:generate mockgen -destination=mocks/url_repository.go -package=mocks . URLRepository
type URLRepository interface {
	FindByShortenedURL(shortURL string) (string, error)
	Save(shortURL string, url string) (storage.ShortenedURL, error)
	Close()
}
