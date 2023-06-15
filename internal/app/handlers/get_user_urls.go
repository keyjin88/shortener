package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) GetUserURL(context RequestContext) {
	uid := context.GetString("uid")
	if uid == "" {
		logger.Log.Infof("uid is empty")
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	originalURL, err := h.shortener.GetShortenedURLByUserID(uid)
	if err != nil {
		logger.Log.Infof("error while shortening url: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error while getting shortened url"})
		return
	}
	if len(originalURL) == 0 {
		logger.Log.Infof("urls not found")
		context.AbortWithStatus(http.StatusNoContent)
		return
	}
	context.JSON(http.StatusOK, originalURL)
}
