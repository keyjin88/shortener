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
	uid := c.GetString("uid")
	shortenURL, err := h.shortener.ShortenURL(uri.String(), uid)
	if err != nil {
		if err.Error() == "URL already exists" {
			c.String(http.StatusConflict, shortenURL)
			return
		}
		logger.Log.Infof("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		c.String(http.StatusBadRequest, "Trouble while shortening url.")
		return
	}
	c.String(http.StatusCreated, shortenURL)
}
