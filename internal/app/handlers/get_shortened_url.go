package handlers

import (
	"fmt"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

// GetShortenedURL retrieves the shortened URL by ID from the Handler's shortener instance
// and performs necessary actions based on the retrieved data.
//
// If the URL is not found by the given ID, it logs a corresponding error message,
// responds with a Bad Request status and a message indicating the URL was not found.
//
// If the retrieved URL is marked as deleted, it aborts the request and responds with a Gone status.
//
// Otherwise, it redirects the request to the original URL associated with the retrieved shortened URL.
func (h *Handler) GetShortenedURL(context RequestContext) {
	id := context.Param("id")
	shortenURL, err := h.shortener.GetShortenedURLByID(id)
	if err != nil {
		logger.Log.Infof("URL not found by id: %s. Error while Api.GetShortenedURLByID()", id)
		context.String(http.StatusBadRequest, fmt.Sprintf("URL not found by id: %s", id))
		return
	}
	if shortenURL.IsDeleted {
		context.AbortWithStatus(http.StatusGone)
	}
	context.Redirect(http.StatusTemporaryRedirect, shortenURL.OriginalURL)
}
