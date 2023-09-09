package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
)

// ShortenURLBatch handles a batch request to shorten multiple URLs.
// It receives the request body as JSON, unmarshals it into a ShortenURLBatchRequest object, and checks for any JSON
// marshalling errors.
// If there are any errors, it logs the error message and returns a Bad Request response with an error message.
// If there are no errors, it retrieves the user ID from the request context, and calls the shortener service to shorten
// the batch of URLs.
// If there is an error while shortening the URLs, it logs the error message and returns a Bad Request response with an
// error message.
// If there are no errors, it returns a Created response with the shortened URLs batch.
func (h *Handler) ShortenURLBatch(c RequestContext) {
	var req storage.ShortenURLBatchRequest
	jsonErr := c.BindJSON(&req)
	if jsonErr != nil {
		logger.Log.Infof(marshalErrorTemplate, jsonErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while marshalling json"})
		return
	}
	uid := c.GetString(key)
	batch, err := h.shortener.ShortenURLBatch(req, uid)
	if err != nil {
		logger.Log.Infof(shorteningErrorTemplate, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": shorteningErrorTemplate})
		return
	}
	c.JSON(http.StatusCreated, batch)
}
