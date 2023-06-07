package inmem

import (
	"errors"
	"fmt"
	"github.com/keyjin88/shortener/internal/app/storage"
	"strconv"
)

type URLRepositoryInMem struct {
	inMemStorage map[string]storage.ShortenedURL
}

func NewURLRepositoryInMem() *URLRepositoryInMem {
	return &URLRepositoryInMem{
		inMemStorage: make(map[string]storage.ShortenedURL),
	}
}

func (ur *URLRepositoryInMem) SaveBatch(urls *[]storage.ShortenedURL) error {
	for _, url := range *urls {
		ur.inMemStorage[url.ShortURL] = url
	}
	return nil
}

func (ur *URLRepositoryInMem) Save(shortenedURL *storage.ShortenedURL) error {
	shortenedURL.UUID = strconv.Itoa(len(ur.inMemStorage))
	ur.inMemStorage[shortenedURL.ShortURL] = *shortenedURL
	return nil
}

func (ur *URLRepositoryInMem) FindByShortenedURL(shortURL string) (string, error) {
	url, ok := ur.inMemStorage[shortURL]
	if !ok {
		return "", fmt.Errorf("URL not found: %v", shortURL)
	}
	return url.OriginalURL, nil
}

func (ur *URLRepositoryInMem) FindByOriginalURL(originalURL string) (string, error) {
	for key, value := range ur.inMemStorage {
		if value.OriginalURL == originalURL {
			return key, nil
		}
	}
	return "", errors.New("URL not found: " + originalURL)
}

// RestoreData восстанавливает состояние БД
func (ur *URLRepositoryInMem) RestoreData(data []storage.ShortenedURL) {
	for _, e := range data {
		ur.inMemStorage[e.ShortURL] = e
	}
}

func (ur *URLRepositoryInMem) Close() {
	//нужен для реализации интерфейса
}

func (r *URLRepositoryInMem) Ping() error {
	//нужен для реализации интерфейса
	return nil
}
