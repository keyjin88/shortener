package handlers

import (
	"net/http"

	"github.com/keyjin88/shortener/internal/app/logger"
)

// DeleteURLs deletes the URLs associated with the given user ID.
//
// The function checks if the user ID is empty and aborts the request with
// a 401 Unauthorized status if it is. It then binds the request body as JSON,
// and if there is an error during the marshaling process, aborts the request
// with a 400 Bad Request status. Next, it invokes the shortener's DeleteURLs
// method to delete the URLs associated with the user ID, and if there is an
// error during this process, aborts the request with a 500 Internal Server Error
// status. Finally, it responds with HTTP status 202 Accepted and an empty response body.
func (h *Handler) DeleteURLs(context RequestContext) {
	uid := context.GetString(key)
	if uid == "" {
		logger.Log.Infof(template)
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var req []string
	jsonErr := context.BindJSON(&req)
	if jsonErr != nil {
		logger.Log.Infof(marshalErrorTemplate, jsonErr)
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
