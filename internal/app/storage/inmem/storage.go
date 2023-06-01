package inmem

import (
	"errors"
	"github.com/keyjin88/shortener/internal/app/storage"
	"strconv"
)

type URLRepositoryInMem struct {
	inMemStorage map[string]string
}

func NewURLRepositoryInMem() *URLRepositoryInMem {
	return &URLRepositoryInMem{
		inMemStorage: make(map[string]string),
	}
}

func (ur *URLRepositoryInMem) Save(shortURL string, url string) (storage.ShortenedURL, error) {
	ur.inMemStorage[shortURL] = url
	shortenedURL := storage.ShortenedURL{
		UUID:        strconv.Itoa(len(ur.inMemStorage)),
		OriginalURL: url,
		ShortURL:    shortURL,
	}
	return shortenedURL, nil
}

func (ur *URLRepositoryInMem) FindByShortenedURL(shortURL string) (string, error) {
	url, ok := ur.inMemStorage[shortURL]
	if !ok {
		return "", errors.New("URL not found: " + shortURL)
	}
	return url, nil
}

// RestoreData восстанавливает состояние БД
func (ur *URLRepositoryInMem) RestoreData(data []storage.ShortenedURL) {
	for _, e := range data {
		ur.inMemStorage[e.ShortURL] = e.OriginalURL
	}
}

func (ur *URLRepositoryInMem) Close() {
	//нужен для реализации интерфейса
}
