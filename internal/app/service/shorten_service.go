package service

import (
	"errors"
	"github.com/google/uuid"
)

type ShortenService struct {
	urlRepository URLRepository
}

func NewShortenService(urlRepository URLRepository) *ShortenService {
	return &ShortenService{urlRepository: urlRepository}
}

func (s *ShortenService) GetShortenedURLByID(id string) (string, bool) {
	return s.urlRepository.FindByShortenedString(id)
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
		if _, ok := s.GetShortenedURLByID(keyStr); !ok {
			err := s.urlRepository.Create(keyStr, url)
			if err != nil {
				return "", err
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
