package service

import "github.com/keyjin88/shortener/internal/app/storage"

type ShortenService struct {
	storage *storage.Storage
}

func NewShortenService(storage *storage.Storage) *ShortenService {
	return &ShortenService{
		storage: storage,
	}
}

func (s *ShortenService) ShortenString(url string) (string, error) {
	uid, err := s.storage.Urls().Create(url)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (s *ShortenService) GetShortenedURL(id string) (string, bool, error) {
	return s.storage.Urls().FindByShortenedString(id)
}
