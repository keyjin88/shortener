package storage

import (
	"bufio"
	"encoding/json"
	"github.com/keyjin88/shortener/internal/app/logger"
	"os"
)

func saveURLJSONToFile(filePath string, data URLJSON) error {
	logger.Log.Infof("Saving to file: %s, data: %s", filePath, data)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	urlJSONAsBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
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

func restoreFromFile(filePath string) (map[string]string, error) {
	logger.Log.Infof("restoring from file: %s", filePath)
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		var temp URLJSON
		err := json.Unmarshal([]byte(line), &temp)
		if err != nil {
			return nil, err
		}
		result[temp.ShortURL] = temp.OriginalURL
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
