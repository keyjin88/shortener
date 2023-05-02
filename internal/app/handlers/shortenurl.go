package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

func (h *Handler) ShortenURL(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Error while read request body. Error from handlers.ShortenURL() :", err)
		c.Status(http.StatusBadRequest)
		_, err = c.Writer.Write([]byte("Invalid request body."))
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		return
	}
	uri, err := url.ParseRequestURI(string(requestBytes))
	if err != nil {
		logger.Error("Invalid url string. Error from url.ParseRequestURI() :", err)
		c.Status(http.StatusBadRequest)
		_, err = c.Writer.Write([]byte("Invalid url string."))
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		return
	}
	shortenString, err := h.shortener.ShortenString(uri.String())
	if err != nil {
		logger.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		c.Status(http.StatusBadRequest)
		_, err = c.Writer.Write([]byte("Trouble while shortening url."))
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		return
	}
	logger.Infof("Запрос на сокращение URL: %s", string(requestBytes))
	c.Status(http.StatusCreated)
	_, err = c.Writer.Write([]byte(h.config.BaseAddress + "/" + shortenString))
	if err != nil {
		return
	}
}
