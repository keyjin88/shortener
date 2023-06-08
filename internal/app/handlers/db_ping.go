package handlers

import (
	"context"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) DBPing(c RequestContext) {
	err := h.pinger.Ping(context.Background())
	if err != nil {
		logger.Log.Errorf("Unable to connect to database: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "pong")
}
