package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getApi() *API {
	api := New()
	api.setupRouter()
	api.configureShortenerService()
	api.config.ParseFlags()
	return api
}

func TestGetShortenedURL(t *testing.T) {
	api := getApi()
	w := httptest.NewRecorder()
	api.GetShortenedURL(gin.CreateTestContextOnly(w, api.router))
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestShortenURL(t *testing.T) {
	api := getApi()
	router := api.router
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("www.ya.ru")))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	assert.EqualValues(t, http.StatusCreated, w.Code)
	assert.True(t, strings.HasPrefix(w.Body.String(), "http://localhost"))
}
