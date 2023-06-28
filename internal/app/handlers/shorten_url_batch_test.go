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

func TestHandler_ShortenURLBatch(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		jsonErr            error
		jsonCallParam      interface{}
		jsonCallCount      int
		serviceReturn      []storage.ShortenURLBatchResponse
		serviceErr         error
		serviceCallCount   int
		expectedStatusCode int
		getStringCallCount int
	}{
		{
			name:    "save successfully",
			jsonErr: nil,
			serviceReturn: []storage.ShortenURLBatchResponse{
				{CorrelationID: "1", ShortURL: "https://localhost:8080/shorten1"},
				{CorrelationID: "2", ShortURL: "https://localhost:8080/shorten2"},
			},
			jsonCallParam: []storage.ShortenURLBatchResponse{
				{CorrelationID: "1", ShortURL: "https://localhost:8080/shorten1"},
				{CorrelationID: "2", ShortURL: "https://localhost:8080/shorten2"},
			},
			jsonCallCount:      1,
			serviceErr:         nil,
			serviceCallCount:   1,
			getStringCallCount: 1,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "error while shortened",
			jsonErr:            nil,
			jsonCallParam:      gin.H{"error": "Error while shortening url"},
			jsonCallCount:      1,
			serviceReturn:      nil,
			serviceErr:         errors.New("some error"),
			serviceCallCount:   1,
			getStringCallCount: 1,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "error while marshalling response",
			jsonErr:            errors.New("some error"),
			jsonCallParam:      gin.H{"error": "Error while marshalling json"},
			jsonCallCount:      1,
			serviceReturn:      nil,
			serviceErr:         nil,
			serviceCallCount:   0,
			getStringCallCount: 0,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequestContext := mocks.NewMockRequestContext(ctrl)
			mockService := mocks.NewMockShortenService(ctrl)

			mockRequestContext.EXPECT().
				BindJSON(gomock.Any()).
				Return(tt.jsonErr)
			mockRequestContext.EXPECT().
				JSON(tt.expectedStatusCode, tt.jsonCallParam).
				Times(tt.jsonCallCount)
			mockService.EXPECT().
				ShortenURLBatch(gomock.Any(), gomock.Any()).
				Return(tt.serviceReturn, tt.serviceErr).
				Times(tt.serviceCallCount)
			mockRequestContext.EXPECT().
				GetString(gomock.Any()).
				Times(tt.getStringCallCount)

			h := &Handler{
				shortener: mockService,
			}
			h.ShortenURLBatch(mockRequestContext)
		})
	}
}
