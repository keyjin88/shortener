package handlers

import (
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) DBPing(c RequestContext) {
	err := h.connectionChecker.Ping()
	if err != nil {
		logger.Log.Errorf("Unable to connect to database: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		c.String(http.StatusOK, "pong")
	}
}
