package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/api/helpers"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	logger    = logrus.New()
	shortener = service.NewShortenService(storage.NewStorage())
)

// Вспомогательная структура для формирования сообщений
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func ShortenURL(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Invalid url string. Error while ShortenURL() :", err)
		helpers.RespondJSON(c, 400, "Invalid url string.")
		return
	}
	shortenString, err := shortener.ShortenString(string(requestBytes))
	if err != nil {
		logger.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		helpers.RespondJSON(c, 400, "Trouble while shortening url.")
		return
	}
	logger.Infof("Запрос на сокращение URL: %s", string(requestBytes))
	c.Status(http.StatusCreated)
	_, err = c.Writer.Write([]byte("http://localhost:8080/" + shortenString))
	if err != nil {
		return
	}
}
func GetShortenedURL(c *gin.Context) {
	id := c.Params.ByName("id")
	originalURL, ok, err := shortener.GetShortenedURL(id)
	if err != nil {
		logger.Error("Trouble while getting shortened url. Error while shortener.GetShortenedURL() :", err)
		helpers.RespondJSON(c, 400, "Trouble while getting shortened url.")
		return
	}
	if !ok {
		logger.Infof("URL not found by id: %s. Error while Api.GetShortenedURL()", id)
		helpers.RespondJSON(c, 400, fmt.Sprintf("URL not found by id: %s", id))
		return
	} else {
		logger.Infof("Запрос на получение URL по id: %s", id)
		c.Header("Location", originalURL)
		c.Status(http.StatusTemporaryRedirect)
		return
	}
}
