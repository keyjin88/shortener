package inmem

import (
	"context"
	"errors"
	"fmt"
	"github.com/keyjin88/shortener/internal/app/storage"
	"strconv"
)

// URLRepositoryInMem is in memory repository.
type URLRepositoryInMem struct {
	inMemStorage map[string]storage.ShortenedURL
}

// NewURLRepositoryInMem creates a new URLRepositoryInMem.
func NewURLRepositoryInMem() *URLRepositoryInMem {
	const repositoryCapacity = 1000
	return &URLRepositoryInMem{
		inMemStorage: make(map[string]storage.ShortenedURL, repositoryCapacity),
	}
}

// SaveBatch saves a batch of USRs to storage.
func (r *URLRepositoryInMem) SaveBatch(urls *[]storage.ShortenedURL) error {
	for i := 0; i < len(*urls); i++ {
		url := &((*urls)[i])
		r.inMemStorage[url.ShortURL] = *url
	}
	return nil
}

// Save method for saving URL in storage.
func (r *URLRepositoryInMem) Save(shortenedURL *storage.ShortenedURL) error {
	shortenedURL.UUID = strconv.Itoa(len(r.inMemStorage))
	r.inMemStorage[shortenedURL.ShortURL] = *shortenedURL
	return nil
}

// FindByShortenedURL find URL by given shortened string in memory.
func (r *URLRepositoryInMem) FindByShortenedURL(shortURL string) (storage.ShortenedURL, error) {
	url, ok := r.inMemStorage[shortURL]
	if !ok {
		return storage.ShortenedURL{}, fmt.Errorf("URL not found: %v", shortURL)
	}
	return url, nil
}

// FindByOriginalURL find shortened URL by original URL.
func (r *URLRepositoryInMem) FindByOriginalURL(originalURL string) (string, error) {
	for key, value := range r.inMemStorage {
		if value.OriginalURL == originalURL {
			return key, nil
		}
	}
	return "", errors.New("URL not found: " + originalURL)
}

// FindAllByUserID find URLs by user ID.
func (r *URLRepositoryInMem) FindAllByUserID(userID string) ([]storage.UsersURLResponse, error) {
	var userURLs []storage.UsersURLResponse
	for _, value := range r.inMemStorage {
		if value.UserID == userID {
			userURLs = append(userURLs, storage.UsersURLResponse{ShortURL: value.ShortURL, OriginalURL: value.OriginalURL})
		}
	}
	return userURLs, nil
}

// RestoreData восстанавливает состояние БД.
func (r *URLRepositoryInMem) RestoreData(data []storage.ShortenedURL) {
	for _, e := range data {
		r.inMemStorage[e.ShortURL] = e
	}
}

// Close method closes the repository.
func (r *URLRepositoryInMem) Close() {
	//нужен для реализации интерфейса
}

// Ping method pings storage.
func (r *URLRepositoryInMem) Ping(_ context.Context) error {
	if r.inMemStorage == nil {
		return errors.New("storage is not initialized")
	}
	return nil
}

// Delete method deleted URLs by given IDs.
func (r *URLRepositoryInMem) Delete(ids []string, _ string) error {
	for _, id := range ids {
		delete(r.inMemStorage, id)
	}
	return nil
}
