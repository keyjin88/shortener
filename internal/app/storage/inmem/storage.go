package inmem

import (
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"strconv"
)

func (ur *URLRepositoryInMem) Create(shortURL string, url string) error {
	ur.inMemStorage[shortURL] = url
	if ur.config.PathToStorageFile != "" {
		err := file.SaveURLJSONToFile(ur.config.PathToStorageFile, file.URLJSON{
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

// RestoreData восстанавливает состояние БД
func (ur *URLRepositoryInMem) RestoreData(data []file.URLJSON) {
	for _, e := range data {
		ur.inMemStorage[e.ShortURL] = e.OriginalURL
	}
}
