package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"time"
)

type ShortenService struct {
	urlRepository URLRepository
	config        Config
}

func NewShortenService(urlRepository URLRepository, pathToStorageFile string) *ShortenService {
	return &ShortenService{
		urlRepository: urlRepository,
		config: Config{
			PathToStorageFile: pathToStorageFile,
		},
	}
}

func (s *ShortenService) GetShortenedURLByID(id string) (string, error) {
	return s.urlRepository.FindByShortenedURL(id)
}

func (s *ShortenService) ShortenURL(url string) (string, error) {
	keyStr, err := s.generateShortenURL(url)
	if err != nil {
		return "", err
	}
	shortURL, err := s.urlRepository.Save(keyStr, url)
	if err != nil {
		return "", err
	}
	if s.config.PathToStorageFile != "" {
		err := s.saveToFile(shortURL, s.config.PathToStorageFile)
		if err != nil {
			return "", err
		}
	}

	return shortURL.ShortURL, err
}

func (s *ShortenService) ShortenURLBatch(request storage.ShortenURLBatchRequest) ([]storage.ShortenURLBatchResponse, error) {
	var urlArray []storage.ShortenedURL
	for _, url := range request {
		shortenURL, err := s.generateShortenURL(url.OriginalURL)
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
		if s.config.PathToStorageFile != "" {
			err := s.saveToFile(url, s.config.PathToStorageFile)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, storage.ShortenURLBatchResponse{CorrelationID: url.CorrelationID, ShortURL: url.ShortURL})
	}
	return result, nil
}

func (s *ShortenService) saveToFile(url storage.ShortenedURL, pathToSave string) error {
	err := file.SaveURLJSONToFile(pathToSave, url)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShortenService) generateShortenURL(originalURL string) (string, error) {
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
