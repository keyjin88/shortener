package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"os"
)

// URLRepositoryFile is in file repository
type URLRepositoryFile struct {
	file *os.File
}

// NewURLRepositoryFile creates a new URLRepositoryFile
func NewURLRepositoryFile(filePath *string) (*URLRepositoryFile, error) {
	file, err := os.OpenFile(*filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Infof("error while opening file: %v", err)
		return nil, err
	}
	urlRepositoryFile := URLRepositoryFile{file: file}
	return &urlRepositoryFile, nil
}

// FindByShortenedURL find URL by given shortened string in file
func (r *URLRepositoryFile) FindByShortenedURL(shortURL string) (storage.ShortenedURL, error) {
	_, err := r.file.Seek(0, 0)
	if err != nil {
		return storage.ShortenedURL{}, err
	}
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		line := scanner.Text()
		var temp storage.ShortenedURL
		err := json.Unmarshal([]byte(line), &temp)
		if err != nil {
			return storage.ShortenedURL{}, err
		}
		if temp.ShortURL == shortURL {
			return temp, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return storage.ShortenedURL{}, err
	}
	return storage.ShortenedURL{}, fmt.Errorf("URL not found: %v", shortURL)
}

// FindByOriginalURL find shortened URL by original URL
func (r *URLRepositoryFile) FindByOriginalURL(originalURL string) (string, error) {
	_, err := r.file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		line := scanner.Text()
		var temp storage.ShortenedURL
		err := json.Unmarshal([]byte(line), &temp)
		if err != nil {
			return "", err
		}
		if temp.OriginalURL == originalURL {
			return temp.ShortURL, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("URL not found: %v", originalURL)
}

// FindAllByUserID find URLs by user ID
func (r *URLRepositoryFile) FindAllByUserID(userID string) ([]storage.UsersURLResponse, error) {
	_, err := r.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	var userURLs []storage.UsersURLResponse
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		line := scanner.Text()
		var temp storage.ShortenedURL
		err := json.Unmarshal([]byte(line), &temp)
		if err != nil {
			return nil, err
		}
		if temp.UserID == userID {
			userURLs = append(userURLs, storage.UsersURLResponse{ShortURL: temp.ShortURL, OriginalURL: temp.OriginalURL})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return userURLs, nil
}

// SaveBatch saves a batch of USRs to storage
func (r *URLRepositoryFile) SaveBatch(urls *[]storage.ShortenedURL) error {
	_, err := r.file.Seek(0, 0)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(r.file)
	for _, shortURL := range *urls {
		err := json.NewEncoder(writer).Encode(shortURL)
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Save method for saving URL in storage
func (r *URLRepositoryFile) Save(data *storage.ShortenedURL) error {
	urlJSONAsBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(r.file)
	defer writer.Flush()
	_, err = writer.Write(urlJSONAsBytes)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

// Close method closes the repository
func (r *URLRepositoryFile) Close() {
	err := r.file.Close()
	if err != nil {
		logger.Log.Infof("error while closing file: %s", r.file.Name())
		return
	}
}

// Ping method pings storage
func (r *URLRepositoryFile) Ping(ctx context.Context) error {
	if _, err := os.Stat(r.file.Name()); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", r.file.Name())
	} else {
		return nil
	}
}

// Delete method deleted URLs by given IDs
func (r *URLRepositoryFile) Delete(ids []string, userID string) error {
	// Открыть файл для чтения и записи
	file, err := os.OpenFile(r.file.Name(), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Прочитать содержимое файла в память
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// Проанализировать каждую строку в файле
	for i, line := range lines {
		// Распаковать JSON в структуру
		var url storage.ShortenedURL
		if err := json.Unmarshal([]byte(line), &url); err != nil {
			return err
		}

		// Проверить, соответствует ли идентификатор URL одному из идентификаторов, переданных в `ids`
		for _, id := range ids {
			if url.ShortURL == id {
				// Пометить запись как удаленную
				url.IsDeleted = true

				// Записать измененную структуру обратно в JSON
				updatedLine, err := json.Marshal(url)
				if err != nil {
					return err
				}

				lines[i] = string(updatedLine)
				break
			}
		}
	}

	// Записать изменения обратно в файл
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
