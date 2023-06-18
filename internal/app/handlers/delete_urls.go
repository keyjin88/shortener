package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) DeleteURLs(context RequestContext) {
	uid := context.GetString("uid")
	if uid == "" {
		logger.Log.Infof("uid is empty")
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var req []string
	jsonErr := context.BindJSON(&req)
	if jsonErr != nil {
		logger.Log.Infof("error while marshalling json data: %v", jsonErr)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	logger.Log.Infof("received delete request: %v", req)
	h.shortener.DeleteURLs(&req, uid)
	context.JSON(http.StatusAccepted, nil)
}
