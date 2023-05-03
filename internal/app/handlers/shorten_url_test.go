package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ShortenURL(t *testing.T) {
	type fields struct {
		shortener *service.ShortenService
		config    *config.Config
	}

	tests := []struct {
		name         string
		fields       fields
		request      *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Create successfully",
			fields:       fields{shortener: service.NewShortenService(storage.NewURLRepositoryInMem()), config: config.NewConfig()},
			request:      httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://ya.ru"))),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Create Bad Request",
			fields:       fields{shortener: service.NewShortenService(storage.NewURLRepositoryInMem()), config: config.NewConfig()},
			request:      httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("htt=.ru"))),
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				shortener: tt.fields.shortener,
				config:    tt.fields.config,
			}
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			context.Request = tt.request
			h.ShortenURL(context)
			assert.EqualValues(t, tt.expectedCode, w.Code)
		})
	}
}
