package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/pkg/errors"
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
	url, err := s.urlRepository.FindByShortenedURL(id)
	if err != nil {
		return storage.ShortenedURL{}, errors.Wrap(err, findURLErrorTemplate)
	}
	return url, nil
}

// GetShortenedURLByUserID returns a URL by given user ID.
func (s *ShortenService) GetShortenedURLByUserID(userID string) ([]storage.UsersURLResponse, error) {
	usersURLResponses, err := s.urlRepository.FindAllByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(err, findURLErrorTemplate)
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
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == pgerrcode.UniqueViolation {
			shortenedURL, err := s.urlRepository.FindByOriginalURL(url)
			if err != nil {
				return "", fmt.Errorf("failed to find original URL: %w", err)
			}
			return s.config.BaseAddress + "/" + shortenedURL, errors.New("URL already exists")
		} else {
			return "", fmt.Errorf("failed to save URL: %w", err)
		}
	}

	return s.config.BaseAddress + "/" + shortURL.ShortURL, nil
}

// ShortenURLBatch returns a short URL batch for the given URLs.
func (s *ShortenService) ShortenURLBatch(request storage.ShortenURLBatchRequest,
	userID string) ([]storage.ShortenURLBatchResponse, error) {
	var urlArray = make([]storage.ShortenedURL, 0, len(request))

	for index := 0; index < len(request); index++ {
		url := request[index]
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
		return nil, fmt.Errorf("failed to save batch: %w", err)
	}

	var result = make([]storage.ShortenURLBatchResponse, 0, len(urlArray))
	for _, u := range urlArray {
		result = append(result, storage.ShortenURLBatchResponse{
			CorrelationID: u.CorrelationID,
			ShortURL:      s.config.BaseAddress + "/" + u.ShortURL,
		})
	}
	return result, nil
}

// DeleteURLs deletes the specified URLs from the repository.
func (s *ShortenService) DeleteURLs(req *[]string, userID string) error {
	err := s.urlRepository.Delete(*req, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete UR")
	}
	return nil
}

func (s *ShortenService) generateShortenURL() (string, error) {
	var attemptCounter = 0
	var maxAttempts = 500
	for {
		randomUUID, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("failed to generate new random UUID: %w", err)
		}
		keyStr := randomUUID.String()[:8]
		// Проверяем что такого URL нет у нас в хранилище
		if _, err := s.GetShortenedURLByID(keyStr); err != nil {
			return keyStr, nil //nolint:nilerr
		} else {
			if attemptCounter > maxAttempts {
				return "", errors.New("too many attempts to generate short url")
			}
			attemptCounter++
		}
	}
}
