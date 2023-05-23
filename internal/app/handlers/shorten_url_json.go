package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenURLJSON(c RequestContext) {
	var req ShortenURLRequest
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Errorf("error while reading request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while reading request"})
		return
	}
	jsonErr := json.Unmarshal(requestBytes, &req)
	if jsonErr != nil {
		logger.Log.Errorf("error while marshalling json data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	result, err := h.shortener.ShortenURL(req.URL)
	if err != nil {
		logger.Log.Errorf("error while shortening url: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while shortening url"})
		return
	}
	response := ShortenURLResponse{Result: h.config.BaseAddress + "/" + result}
	logger.Log.Infof("Запрос на сокращение URL: %s, результат: %s", string(requestBytes), response.Result)
	c.JSON(http.StatusCreated, response)
}
