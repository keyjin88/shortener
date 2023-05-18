package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

func (h *Handler) shortenURLJSON(c RequestContext) {
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
	c.JSON(http.StatusCreated, ShortenURLResponse{Result: h.config.BaseAddress + "/" + result})
}
