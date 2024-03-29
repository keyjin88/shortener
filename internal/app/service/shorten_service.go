package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/shortener/internal/app/storage"
	"time"
)

// NewShortenService creates a new shortener.
func NewShortenService(urlRepository URLRepository, baseAddress string) *ShortenService {
	return &ShortenService{
		urlRepository: urlRepository,
		config: &Config{
			BaseAddress: baseAddress,
		},
	}
}

// GetShortenedURLByID returns a URL by given ID.
func (s *ShortenService) GetShortenedURLByID(id string) (storage.ShortenedURL, error) {
	return s.urlRepository.FindByShortenedURL(id)
}

// GetShortenedURLByUserID returns a URL by given user ID.
func (s *ShortenService) GetShortenedURLByUserID(userID string) ([]storage.UsersURLResponse, error) {
	usersURLResponses, err := s.urlRepository.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	for i, u := range usersURLResponses {
		usersURLResponses[i].ShortURL = s.config.BaseAddress + "/" + u.ShortURL
	}
	return usersURLResponses, nil
}

// ShortenURL returns shorten URL by given original URL and user ID.
func (s *ShortenService) ShortenURL(url string, userID string) (string, error) {
	keyStr, err := s.generateShortenURL()
	if err != nil {
		return "", err
	}
	shortURL := storage.ShortenedURL{
		UserID:      userID,
		ShortURL:    keyStr,
		OriginalURL: url,
	}
	err = s.urlRepository.Save(&shortURL)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			shortenedURL, err := s.urlRepository.FindByOriginalURL(url)
			if err != nil {
				return "", err
			}
			return s.config.BaseAddress + "/" + shortenedURL, errors.New("URL already exists")
		} else {
			return "", err
		}
	}

	return s.config.BaseAddress + "/" + shortURL.ShortURL, err
}

// ShortenURLBatch returns a short URL batch for the given URLs
func (s *ShortenService) ShortenURLBatch(request storage.ShortenURLBatchRequest, userID string) ([]storage.ShortenURLBatchResponse, error) {
	var urlArray []storage.ShortenedURL
	for _, url := range request {
		shortenURL, err := s.generateShortenURL()
		if err != nil {
			return nil, err
		}
		shortenedURL := storage.ShortenedURL{
			UserID:        userID,
			CreatedAt:     time.Now(),
			OriginalURL:   url.OriginalURL,
			ShortURL:      shortenURL,
			CorrelationID: url.CorrelationID,
		}
		urlArray = append(urlArray, shortenedURL)
	}

	err := s.urlRepository.SaveBatch(&urlArray)
	if err != nil {
		return nil, err
	}
	var result []storage.ShortenURLBatchResponse
	for _, url := range urlArray {
		result = append(result, storage.ShortenURLBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      s.config.BaseAddress + "/" + url.ShortURL})
	}
	return result, nil
}

// DeleteURLs deletes the specified URLs from the repository
func (s *ShortenService) DeleteURLs(req *[]string, userID string) error {
	return s.urlRepository.Delete(*req, userID)
}

func (s *ShortenService) generateShortenURL() (string, error) {
	var attemptCounter = 0
	var maxAttempts = 500
	for {
		randomUUID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		keyStr := randomUUID.String()[:8]
		// Проверяем что такого URL нет у нас в хранилище
		if _, err := s.GetShortenedURLByID(keyStr); err != nil {
			return keyStr, nil
		} else {
			if attemptCounter > maxAttempts {
				return "", errors.New("too many attempts to generate short url")
			}
			attemptCounter += 1
		}
	}
}
