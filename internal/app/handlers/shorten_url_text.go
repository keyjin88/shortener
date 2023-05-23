package handlers

import (
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
	"net/url"
)

func (h *Handler) ShortenURLText(c RequestContext) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.String(http.StatusBadRequest, "Invalid request body.")
		return
	}
	uri, err := url.ParseRequestURI(string(requestBytes))
	if err != nil {
		logger.Log.Infof("error while parsing URL: %v", err)
		c.String(http.StatusBadRequest, "Invalid url string.")
		return
	}
	shortenString, err := h.shortener.ShortenURL(uri.String())
	if err != nil {
		logger.Log.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		c.String(http.StatusBadRequest, "Trouble while shortening url.")
		return
	}
	logger.Log.Infof("Запрос на сокращение URL: %s, результат: %s", string(requestBytes), shortenString)
	c.String(http.StatusCreated, h.config.BaseAddress+"/"+shortenString)
}
