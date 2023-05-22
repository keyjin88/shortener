package storage

import (
	"bufio"
	"encoding/json"
	"github.com/keyjin88/shortener/internal/app/logger"
	"os"
)

type UrlJson struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func SaveUrlJsonToFile(filePath string, data UrlJson) error {
	logger.Log.Infof("Saving to file: %s, data: %s", filePath, data)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Errorf("error opening file: %v", err)
		return err
	}
	defer file.Close()

	urlJsonAsBytes, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("error marshalling data to json: %v", err)
		return err
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.Write(urlJsonAsBytes)
	if err != nil {
		logger.Log.Errorf("error writing to file: %v", err)
		return err
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		logger.Log.Errorf("error writing carriage return to file: %v", err)
		return err
	}
	return nil
}

func RestoreFromFile(filePath string) (map[string]string, error) {
	logger.Log.Infof("restoring from file: %s", filePath)
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Errorf("can't open file: %v", err)
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		var temp UrlJson
		err := json.Unmarshal([]byte(line), &temp)
		if err != nil {
			logger.Log.Errorf("can't open file: %v", err)
			return nil, err
		}
		result[temp.ShortURL] = temp.OriginalURL
	}
	if err := scanner.Err(); err != nil {
		logger.Log.Errorf("error while reading file: %v", err)
		return nil, err
	}
	return result, nil
}
