package handlers

import (
	"github.com/gin-gonic/gin"
)

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

func (h *Handler) shortenURLJSON(c RequestContext) {
	var req ShortenURLRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := h.shortener.ShortenURL(req.URL)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ShortenURLResponse{Result: h.config.BaseAddress + "/" + result})
}
