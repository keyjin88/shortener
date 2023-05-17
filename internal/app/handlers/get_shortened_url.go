package handlers

import (
	"fmt"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) GetShortenedURL(context RequestContext) {
	id := context.Param("id")
	originalURL, ok := h.shortener.GetShortenedURLByID(id)
	if !ok {
		logger.Log.Infof("URL not found by id: %s. Error while Api.GetShortenedURLByID()", id)
		context.String(http.StatusBadRequest, fmt.Sprintf("URL not found by id: %s", id))
		return
	} else {
		logger.Log.Infof("Запрос на получение URL по id: %s, originalURL: %s", id, originalURL)
		context.Redirect(http.StatusTemporaryRedirect, originalURL)
		return
	}
}
