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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jsonErr := json.Unmarshal(requestBytes, &req)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	result, err := h.shortener.ShortenURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := ShortenURLResponse{Result: h.config.BaseAddress + "/" + result}
	logger.Log.Infof("Запрос на сокращение URL: %s, результат: %s", string(requestBytes), response.Result)
	c.JSON(http.StatusCreated, response)
}
