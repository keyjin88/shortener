package service

import (
	"github.com/google/uuid"
	"github.com/keyjin88/shortener/internal/app/storage"
)

type ShortenService struct {
	urlRepository storage.URLRepository
}

func NewShortenService(urlRepository storage.URLRepository) *ShortenService {
	return &ShortenService{urlRepository: urlRepository}
}

func (s *ShortenService) ShortenString(url string) (string, error) {
	for {
		u, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		keyStr := u.String()[:8]
		_, ok := s.GetShortenedURLByID(keyStr)
		if !ok {
			s.urlRepository.Create(keyStr, url)
			return keyStr, nil
		}
	}
}

func (s *ShortenService) GetShortenedURLByID(id string) (string, bool) {
	return s.urlRepository.FindByShortenedString(id)
}
