package inmem

import (
	"errors"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"strconv"
)

type URLRepositoryInMem struct {
	config       storage.Config
	inMemStorage map[string]string
}

func NewURLRepositoryInMem(pathToStorageFile string) *URLRepositoryInMem {
	return &URLRepositoryInMem{
		config: storage.Config{
			PathToStorageFile: pathToStorageFile,
		},
		inMemStorage: make(map[string]string),
	}
}

func (ur *URLRepositoryInMem) Save(shortURL string, url string) error {
	ur.inMemStorage[shortURL] = url
	if ur.config.PathToStorageFile != "" {
		err := file.SaveURLJSONToFile(ur.config.PathToStorageFile, storage.ShortenedURL{
			UUID:        strconv.Itoa(len(ur.inMemStorage)),
			OriginalURL: url,
			ShortURL:    shortURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ur *URLRepositoryInMem) FindByShortenedURL(shortURL string) (string, error) {
	url, ok := ur.inMemStorage[shortURL]
	if !ok {
		return "", errors.New("URL not found: " + shortURL)
	}
	return url, nil
}

// RestoreData восстанавливает состояние БД
func (ur *URLRepositoryInMem) RestoreData(data []storage.ShortenedURL) {
	for _, e := range data {
		ur.inMemStorage[e.ShortURL] = e.OriginalURL
	}
}
