package handlers

import (
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
	"net/url"
)

// ShortenURLText is a handler function that shortens a given URL string and returns the shortened URL.
// It takes a RequestContext c and sets the Content-Type header to "text/plain".
// It reads the raw request data using c.GetRawData() and attempts to parse it as a URL.
// If there is an error while reading the request data, it logs the error and returns a response with status code 400
// and the message "Invalid request body.".
// If there is an error while parsing the URL, it logs the error and returns a response with status code 400 and the
// message "Invalid url string.".
// It then retrieves the "uid" value from the RequestContext and calls h.shortener.ShortenURL to shorten the URL.
// If there is an error while shortening the URL and the error message is "URL already exists", it returns a response
// with status code 409 and the shortened URL.
// If there is any other error while shortening the URL, it logs the error and returns a response with status code 400
// and the message "Trouble while shortening url.".
// If there are no errors, it returns a response with status code 201 and the shortened URL.
func (h *Handler) ShortenURLText(c RequestContext) {
	c.Header("Content-Type", "text/plain")
	requestBytes, err := c.GetRawData()
	if err != nil {
		logger.Log.Infof("error while reading request: %v", err)
		c.String(http.StatusBadRequest, "Invalid request body.")
		return
	}
	uri, err := url.ParseRequestURI(string(requestBytes))
	if err != nil {
		logger.Log.Infof("error while parsing URL: %v", err)
		c.String(http.StatusBadRequest, "Invalid url string.")
		return
	}
	uid := c.GetString("uid")
	shortenURL, err := h.shortener.ShortenURL(uri.String(), uid)
	if err != nil {
		if err.Error() == "URL already exists" {
			c.String(http.StatusConflict, shortenURL)
			return
		}
		logger.Log.Infof("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		c.String(http.StatusBadRequest, "Trouble while shortening url.")
		return
	}
	c.String(http.StatusCreated, shortenURL)
}
