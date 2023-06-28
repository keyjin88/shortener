package handlers

import (
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
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := h.shortener.DeleteURLs(&req, uid)
	if err != nil {
		logger.Log.Infof("error while deleting urls: %v", err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusAccepted, nil)
}
