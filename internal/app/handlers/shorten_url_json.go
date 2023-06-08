package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
)

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
	result, err := h.shortener.ShortenURL(req.URL)
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
