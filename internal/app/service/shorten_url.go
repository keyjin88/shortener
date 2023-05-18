package service

import (
	"github.com/google/uuid"
)

func (s *ShortenService) ShortenURL(url string) (string, error) {
	for {
		randomUUID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		keyStr := randomUUID.String()[:8]
		_, ok := s.GetShortenedURLByID(keyStr)
		if !ok {
			s.urlRepository.Create(keyStr, url)
			return keyStr, nil
		}
	}
}
