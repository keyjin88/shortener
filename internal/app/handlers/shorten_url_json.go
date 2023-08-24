package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
)

// ShortenURLJSON is a method that handles shortening the given URL in JSON format.
// It reads the request data, unmarshals it into a ShortenURLRequest object, and then calls the shortener's ShortenURL method.
// If the URL already exists, it returns an HTTP status code of 409 (Conflict) with the existing short URL.
// If there are any errors during the process, it returns an appropriate HTTP status code and error message.
// Finally, it returns the shortened URL in JSON format.
func (h *Handler) ShortenURLJSON(c RequestContext) {
	var req storage.ShortenURLRequest
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading request"})
		return
	}
	jsonErr := json.Unmarshal(requestBytes, &req)
	if jsonErr != nil {
		logger.Log.Infof("error while marshalling json data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	uid := c.GetString("uid")
	result, err := h.shortener.ShortenURL(req.URL, uid)
	response := storage.ShortenURLResponse{Result: result}
	if err != nil {
		if err.Error() == "URL already exists" {
			logger.Log.Infof("error while shortening url: %v", err)
			c.JSON(http.StatusConflict, response)
			return
		}
		logger.Log.Infof("error while shortening url: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while shortening url"})
		return
	}
	c.JSON(http.StatusCreated, response)
}
