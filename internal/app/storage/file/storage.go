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

type URLRepositoryFile struct {
	file *os.File
}

func NewURLRepositoryFile(filePath *string) (*URLRepositoryFile, error) {
	file, err := os.OpenFile(*filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Infof("error while opening file: %v", err)
		return nil, err
	}
	urlRepositoryFile := URLRepositoryFile{file: file}
	return &urlRepositoryFile, nil
}

func (r *URLRepositoryFile) FindByShortenedURL(shortURL string) (string, error) {
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
		if temp.ShortURL == shortURL {
			return temp.OriginalURL, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("URL not found: %v", shortURL)
}

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

func (r *URLRepositoryFile) Close() {
	err := r.file.Close()
	if err != nil {
		logger.Log.Infof("error while closing file: %s", r.file.Name())
		return
	}
}

func (r *URLRepositoryFile) Ping(ctx context.Context) error {
	if _, err := os.Stat(r.file.Name()); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", r.file.Name())
	} else {
		return nil
	}
}
