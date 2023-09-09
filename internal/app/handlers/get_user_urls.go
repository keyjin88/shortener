package handlers

import (
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
)

// GetUserURL is a function that handles a request to retrieve the original URL associated with a user.
// It expects a RequestContext object as the input parameter.
//
// If the "uid" parameter in the request context is empty, it logs an informational message and returns an unauthorized
// status.
//
// If there is an error while retrieving the shortened URL associated with the provided user ID,
// it logs an error message and returns a bad request status.
//
// If there are no original URLs found for the user, it logs an informational message and returns a no content status.
//
// If everything is successful, it returns the original URL associated with the user in JSON format with a status OK.
func (h *Handler) GetUserURL(context RequestContext) {
	uid := context.GetString(key)
	if uid == "" {
		logger.Log.Infof("uid is empty")
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	originalURL, err := h.shortener.GetShortenedURLByUserID(uid)
	if err != nil {
		logger.Log.Infof(shorteningErrorTemplate, err)
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(originalURL) == 0 {
		logger.Log.Infof("urls not found")
		context.AbortWithStatus(http.StatusNoContent)
		return
	}
	context.JSON(http.StatusOK, originalURL)
}
