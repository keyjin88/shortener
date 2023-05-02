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
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			if tt.name == "Get successfully" {
				context.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://ya.ru")))
				h.ShortenURL(context)
				shortenUrl := w.Body.String()
				logger.Info("Shorten url: ", shortenUrl)

				// Обновляем контекст и responseRecorder
				w = httptest.NewRecorder()
				context, _ = gin.CreateTestContext(w)
				context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
				context.Params = []gin.Param{
					{
						Key:   "id",
						Value: shortenUrl[1:],
					},
				}
				h.GetShortenedURL(context)
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

//func TestAPI_GetShortenedURL(t *testing.T) {
//	api := getAPI()
//	w, context := getGinTestContext()
//	context.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("http://ya.ru")))
//	api.ShortenURL(context)
//	shortenUrl := w.Body.String()
//	logger.Info("Shorten url: ", shortenUrl)
//
//	context.Request = httptest.NewRequest(http.MethodGet, shortenUrl, nil)
//	context.Params = []gin.Param{
//		{
//			Key:   "id",
//			Value: shortenUrl[1:],
//		},
//	}
//	api.GetShortenedURL(context)
//	// я совершенно не понимаю, почему я проставляю статус в методе api,
//	// но во время выполнения тестов он не меняется  ((
//	// при этом после первого запроса статус во врайтере меняется на 201
//	//assert.EqualValues(t, http.StatusTemporaryRedirect, w.Code)
//	assert.EqualValues(t, "http://ya.ru", w.Header().Get("Location"))
//}

//func TestAPI_GetShortenedURL(t *testing.T) {
//	api := getAPI()
//	w, context := getGinTestContext()
//
//	shortedUrl, _ := api.storage.Urls().Create("http://www.ya.ru")
//	logger.Info("Shorten url: /", shortedUrl)
//
//	context.Request = httptest.NewRequest(http.MethodGet, "/"+shortedUrl, nil)
//	context.Params = []gin.Param{
//		{
//			Key:   "id",
//			Value: shortedUrl,
//		},
//	}
//	api.GetShortenedURL(context)
//	//assert.EqualValues(t, http.StatusTemporaryRedirect, w.Code)
//	assert.EqualValues(t, "http://www.ya.ru", w.Header().Get("Location"))
//}

//func TestAPI_GetShortenedURL_BadRequest(t *testing.T) {
//	api := getAPI()
//	w, context := getGinTestContext()
//	context.Request = httptest.NewRequest(http.MethodGet, "/shortedURL", nil)
//	api.GetShortenedURL(context)
//	assert.EqualValues(t, http.StatusBadRequest, w.Code)
//}
