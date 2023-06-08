package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keyjin88/shortener/internal/app/storage"
	"time"
)

func NewShortenService(urlRepository URLRepository, baseAddress string) *ShortenService {
	return &ShortenService{
		urlRepository: urlRepository,
		config: &Config{
			BaseAddress: baseAddress,
		},
	}
}

func (s *ShortenService) GetShortenedURLByID(id string) (string, error) {
	return s.urlRepository.FindByShortenedURL(id)
}

func (s *ShortenService) ShortenURL(url string) (string, error) {
	keyStr, err := s.generateShortenURL()
	if err != nil {
		return "", err
	}
	shortURL := storage.ShortenedURL{
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

func (s *ShortenService) ShortenURLBatch(request storage.ShortenURLBatchRequest) ([]storage.ShortenURLBatchResponse, error) {
	var urlArray []storage.ShortenedURL
	for _, url := range request {
		shortenURL, err := s.generateShortenURL()
		if err != nil {
			return nil, err
		}
		shortenedURL := storage.ShortenedURL{
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

func (s *ShortenService) PingDB() error {
	return s.urlRepository.Ping()
}
