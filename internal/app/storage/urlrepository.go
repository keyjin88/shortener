package storage

import (
	"errors"
	"github.com/google/uuid"
)

type UrlRepository struct {
	storage *Storage
}

var (
	inmemStorage map[string]string
)

func init() {
	inmemStorage = make(map[string]string)
}

func (ur *UrlRepository) Create(url string) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuidStr := u.String()[:8]
	inmemStorage[uuidStr] = url
	return uuidStr, nil
}

func (ur *UrlRepository) FindByShortenedString(id string) (string, bool, error) {
	for shortenStr, url := range inmemStorage {
		if shortenStr == id {
			return url, true, nil
		} else {
			return "", false, nil
		}
	}
	return "", false, errors.New("something went wrong")
}
