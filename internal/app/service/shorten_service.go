package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/keyjin88/shortener/internal/app/storage/file"
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
	var attemptCounter = 0
	var maxAttempts = 500
	for {
		randomUUID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		keyStr := randomUUID.String()[:8]
		if _, err := s.GetShortenedURLByID(keyStr); err == nil {
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
			return keyStr, nil
		} else {
			if attemptCounter > maxAttempts {
				return "", errors.New("too many attempts")
			}
			attemptCounter += 1
		}
	}
}

func (s *ShortenService) saveToFile(url storage.ShortenedURL, pathToSave string) error {
	err := file.SaveURLJSONToFile(pathToSave, url)
	if err != nil {
		return err
	}
	return nil
}
