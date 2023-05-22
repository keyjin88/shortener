package storage

import (
	"github.com/keyjin88/shortener/internal/app/logger"
	"strconv"
)

type URLRepositoryInMem struct {
	config Config
}

func NewURLRepositoryInMem(pathToStorageFile string) *URLRepositoryInMem {
	return &URLRepositoryInMem{
		config: Config{
			PathToStorageFile: pathToStorageFile,
		},
	}
}

var (
	inMemStorage map[string]string
)

func init() {
	inMemStorage = make(map[string]string)
}

func (ur *URLRepositoryInMem) Create(shortURL string, url string) error {
	inMemStorage[shortURL] = url
	if ur.config.PathToStorageFile != "" {
		err := SaveUrlJsonToFile(ur.config.PathToStorageFile, UrlJson{
			UUID:        strconv.Itoa(len(inMemStorage)),
			OriginalURL: url,
			ShortURL:    shortURL,
		})
		if err != nil {
			logger.Log.Errorf("Error while saving to file: %v", err)
			return err
		}
	}
	return nil
}

func (ur *URLRepositoryInMem) FindByShortenedString(id string) (string, bool) {
	url, ok := inMemStorage[id]
	return url, ok
}

func (ur *URLRepositoryInMem) RestoreFromFile() error {
	if ur.config.PathToStorageFile != "" {
		result, err := RestoreFromFile(ur.config.PathToStorageFile)
		if err != nil {
			logger.Log.Errorf("error while restoring from file file: %v", err)
			return err
		}
		inMemStorage = result
	} else {
		inMemStorage = make(map[string]string)
	}
	return nil
}
