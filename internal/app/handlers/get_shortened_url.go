package handlers

import (
	"fmt"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

func (h *Handler) GetShortenedURL(context RequestContext) {
	id := context.Param("id")
	shortenUrl, err := h.shortener.GetShortenedURLByID(id)
	if err != nil {
		logger.Log.Infof("URL not found by id: %s. Error while Api.GetShortenedURLByID()", id)
		context.String(http.StatusBadRequest, fmt.Sprintf("URL not found by id: %s", id))
		return
	}
	if shortenUrl.IsDeleted {
		context.AbortWithStatus(http.StatusGone)
	}
	context.Redirect(http.StatusTemporaryRedirect, shortenUrl.OriginalURL)
	return
}
