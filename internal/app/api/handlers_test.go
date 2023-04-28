package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getAPI() *API {
	api := New()
	api.setupRouter()
	api.configureShortenerService()
	return api
}

func TestAPI_ShortenURL(t *testing.T) {
	api := getAPI()
	router := api.router
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("www.ya.ru")))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	assert.EqualValues(t, http.StatusCreated, w.Code)
	assert.EqualValues(t, len(w.Body.String()), 8)
}

func TestAPI_GetShortenedURL(t *testing.T) {
	api := getAPI()
	router := api.router
	request := httptest.NewRequest(http.MethodGet, "/shortedURL", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}
