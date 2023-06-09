package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
	"testing"
)

func TestHandler_shortenURLJSON(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		url                string
		getRowDataReturn   getRowDataReturn
		serviceReturn      shortenURLReturn
		expectedCode       int
		expectedBody       interface{}
		expectedStringCall int
		shortenStringCall  int
		expectedJSONCall   int
		getStringCallCount int
	}{
		{
			name:               "Save successfully",
			url:                "https://www.yandex.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte(`{"url":"https://www.yandex.ru"}`), error: nil},
			serviceReturn:      shortenURLReturn{result: "/SHORTEN", error: nil},
			expectedCode:       http.StatusCreated,
			expectedBody:       storage.ShortenURLResponse{Result: "/SHORTEN"},
			shortenStringCall:  1,
			expectedJSONCall:   1,
			getStringCallCount: 1,
		},
		{
			name:               "Invalid request body",
			url:                "ww.bad url",
			getRowDataReturn:   getRowDataReturn{result: nil, error: errors.New("error from GetRowData()")},
			serviceReturn:      shortenURLReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       gin.H{"error": "Error while reading request"},
			shortenStringCall:  0,
			expectedJSONCall:   1,
			getStringCallCount: 0,
		},
		{
			name:               "Error while json unmarshal",
			url:                "htt=bad_url.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("htt=bad_url.ru"), error: nil},
			serviceReturn:      shortenURLReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       gin.H{"error": "Error while marshalling json"},
			shortenStringCall:  0,
			expectedJSONCall:   1,
			getStringCallCount: 0,
		},
		{
			name:               "Save successfully",
			url:                "https://www.yandex.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte(`{"url":"https://www.yandex.ru"}`), error: nil},
			serviceReturn:      shortenURLReturn{result: "", error: errors.New("error from shorten service")},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       gin.H{"error": "Error while shortening url"},
			shortenStringCall:  1,
			expectedJSONCall:   1,
			getStringCallCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockShortenService(ctrl)
			mockService.EXPECT().ShortenURL(tt.url, gomock.Any()).
				Times(tt.shortenStringCall).
				Return(tt.serviceReturn.result, tt.serviceReturn.error)

			mockRequestContext := mocks.NewMockRequestContext(ctrl)
			mockRequestContext.EXPECT().String(tt.expectedCode, tt.expectedBody).
				Times(tt.expectedStringCall)
			mockRequestContext.EXPECT().GetRawData().
				Return(tt.getRowDataReturn.result, tt.getRowDataReturn.error)
			mockRequestContext.EXPECT().JSON(tt.expectedCode, tt.expectedBody).
				Times(tt.expectedJSONCall)
			mockRequestContext.EXPECT().GetString(gomock.Any()).
				Times(tt.getStringCallCount)

			h := &Handler{
				shortener: mockService,
			}
			h.ShortenURLJSON(mockRequestContext)
		})
	}
}
