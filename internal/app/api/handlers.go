package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/api/helpers"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	logger = logrus.New()
)

// Вспомогательная структура для формирования сообщений
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func (api *API) ShortenURL(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Invalid url string. Error while ShortenURL() :", err)
		helpers.RespondJSON(c, 400, "Invalid url string.")
		return
	}
	shortenString, err := api.shortener.ShortenString(string(requestBytes))
	if err != nil {
		logger.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		helpers.RespondJSON(c, 400, "Trouble while shortening url.")
		return
	}
	logger.Infof("Запрос на сокращение URL: %s", string(requestBytes))
	c.Status(http.StatusCreated)
	_, err = c.Writer.Write([]byte(api.config.Flags.BaseAddr + shortenString))
	if err != nil {
		return
	}
}
func (api *API) GetShortenedURL(c *gin.Context) {
	id := c.Params.ByName("id")
	originalURL, ok, err := api.shortener.GetShortenedURL(id)
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
