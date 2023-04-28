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

func TestGetShortenedURL(t *testing.T) {
	router := gin.Default()
	router.GET("/:id", GetShortenedURL)
	request := httptest.NewRequest(http.MethodGet, "/UUID", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestShortenURL(t *testing.T) {
	router := gin.Default()
	router.POST("/", ShortenURL)
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("www.ya.ru")))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	assert.EqualValues(t, http.StatusCreated, w.Code)
	assert.True(t, strings.HasPrefix(w.Body.String(), "http://localhost"))
}
