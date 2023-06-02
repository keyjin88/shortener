package inmem

import (
	"errors"
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

func (ur *URLRepositoryInMem) Save(shortURL string, url string) (storage.ShortenedURL, error) {
	shortenedURL := storage.ShortenedURL{
		UUID:        strconv.Itoa(len(ur.inMemStorage)),
		OriginalURL: url,
		ShortURL:    shortURL,
	}
	ur.inMemStorage[shortURL] = shortenedURL
	return shortenedURL, nil
}

func (ur *URLRepositoryInMem) FindByShortenedURL(shortURL string) (string, error) {
	url, ok := ur.inMemStorage[shortURL]
	if !ok {
		return "", errors.New("URL not found: " + shortURL)
	}
	return url.OriginalURL, nil
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
