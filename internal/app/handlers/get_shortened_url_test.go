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

func TestHandler_GetShortenedURL(t *testing.T) {
	type fields struct {
		shortener *service.ShortenService
		config    *config.Config
	}

	tests := []struct {
		name         string
		fields       fields
		requests     []*http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Get successfully",
			fields:       fields{shortener: service.NewShortenService(storage.NewURLRepositoryInMem()), config: config.NewConfig()},
			requests:     []*http.Request{},
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "Get Bad Request",
			fields:       fields{shortener: service.NewShortenService(storage.NewURLRepositoryInMem()), config: config.NewConfig()},
			requests:     []*http.Request{httptest.NewRequest(http.MethodGet, "/shortedURL", nil)},
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
			if tt.name == "Get successfully" {
				context.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://ya.ru")))
				logger.Info("Status before ShortenURL(): ", context.Writer.Status())
				h.ShortenURL(context)
				logger.Info("Status after ShortenURL(): ", context.Writer.Status())
				shortenURL := w.Body.String()

				// Обновляем контекст и responseRecorder
				w = httptest.NewRecorder()
				context, _ = gin.CreateTestContext(w)
				context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
				context.Params = []gin.Param{
					{
						Key:   "id",
						Value: shortenURL[1:],
					},
				}
				logger.Info("Status before GetShortenedURL(): ", context.Writer.Status())
				h.GetShortenedURL(context)
				logger.Info("Status after GetShortenedURL(): ", context.Writer.Status())
				assert.EqualValues(t, tt.expectedCode, w.Code)
				assert.EqualValues(t, "https://ya.ru", w.Header().Get("Location"))
			} else {
				context.Request = tt.requests[0]
				h.GetShortenedURL(context)
				assert.EqualValues(t, tt.expectedCode, w.Code)
			}
		})
	}
}
