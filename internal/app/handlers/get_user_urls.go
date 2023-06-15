package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) GetUserURL(context RequestContext) {
	uid := context.GetString("uid")
	originalURL, err := h.shortener.GetShortenedURLByUserID(uid)
	if err != nil {
		logger.Log.Infof("error while shortening url: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error while getting shortened url"})
		return
	}
	context.JSON(http.StatusOK, originalURL)
}
