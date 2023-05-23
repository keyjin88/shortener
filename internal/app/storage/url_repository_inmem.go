package storage

import (
	"strconv"
)

type URLRepositoryInMem struct {
	config       Config
	inMemStorage map[string]string
}

func NewURLRepositoryInMem(pathToStorageFile string) *URLRepositoryInMem {
	return &URLRepositoryInMem{
		config: Config{
			PathToStorageFile: pathToStorageFile,
		},
		inMemStorage: make(map[string]string),
	}
}

func (ur *URLRepositoryInMem) Create(shortURL string, url string) error {
	ur.inMemStorage[shortURL] = url
	if ur.config.PathToStorageFile != "" {
		err := saveURLJSONToFile(ur.config.PathToStorageFile, URLJSON{
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

func (ur *URLRepositoryInMem) FindByShortenedString(id string) (string, bool) {
	url, ok := ur.inMemStorage[id]
	return url, ok
}

// RestoreDataFromFile восстанавливает состояние БД из файла,
// если произошла ошибка оставляем кеш пустым
func (ur *URLRepositoryInMem) RestoreDataFromFile(filePath string) error {
	result, err := restoreFromFile(filePath)
	if err != nil {
		return err
	} else {
		ur.inMemStorage = result
	}
	return nil
}
