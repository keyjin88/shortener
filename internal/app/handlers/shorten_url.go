package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
)

func (h *Handler) ShortenURL(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Error while read request body. Error from handlers.ShortenURL() :", err)
		c.String(http.StatusBadRequest, "Invalid request body.")
		return
	}
	uri, err := url.ParseRequestURI(string(requestBytes))
	if err != nil {
		logger.Error("Invalid url string. Error from url.ParseRequestURI() :", err)
		c.String(http.StatusBadRequest, "Invalid url string.")
		return
	}
	shortenString, err := h.shortener.ShortenString(uri.String())
	if err != nil {
		logger.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		c.String(http.StatusBadRequest, "Trouble while shortening url.")
		return
	}
	logger.Infof("Запрос на сокращение URL: %s, результат: %s", string(requestBytes), shortenString)
	c.String(http.StatusCreated, h.config.BaseAddress+"/"+shortenString)
}
