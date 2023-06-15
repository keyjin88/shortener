package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
)

func (h *Handler) ShortenURLBatch(c RequestContext) {
	var req storage.ShortenURLBatchRequest
	jsonErr := c.BindJSON(&req)
	if jsonErr != nil {
		logger.Log.Infof("error while marshalling json data: %v", jsonErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	uid := c.GetString("uid")
	batch, err := h.shortener.ShortenURLBatch(req, uid)
	if err != nil {
		logger.Log.Infof("error while shortening url: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while shortening url"})
		return
	}
	c.JSON(http.StatusCreated, batch)
}
