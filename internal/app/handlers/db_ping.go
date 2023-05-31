package handlers

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) DBPing(c RequestContext) {
	conn, err := pgx.Connect(context.Background(), h.config.DataBaseDSN)
	if err != nil {
		logger.Log.Errorf("Unable to connect to database: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		c.String(http.StatusOK, "pong")
	}
	defer conn.Close(context.Background())
}
