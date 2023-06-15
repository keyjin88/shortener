package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"net/http"
	"testing"
)

func TestHandler_GetUserURL(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type getShortenURLReturn struct {
		result []storage.UsersURLResponse
		error  error
	}

	tests := []struct {
		name                     string
		serviceReturn            getShortenURLReturn
		getStringReturn          string
		statusCode               int
		getShortenedURLCallCount int
		abortWithStatusCallCount int
		jsonCallCount            int
	}{
		{
			name: "success",
			serviceReturn: getShortenURLReturn{
				result: []storage.UsersURLResponse{
					{ShortURL: "http://localhost:8080/2116a093", OriginalURL: "https://www.rambler.ru/%20"},
					{ShortURL: "http://localhost:8080/4f259096", OriginalURL: "https://www.rambler.ru/%20"},
					{ShortURL: "http://localhost:8080/6351b00c", OriginalURL: "https://www.rambler.ru/%20"},
				},
				error: nil,
			},
			getStringReturn:          "userID",
			statusCode:               http.StatusOK,
			getShortenedURLCallCount: 1,
			abortWithStatusCallCount: 0,
			jsonCallCount:            1,
		},
		{
			name:                     "uid is empty",
			getStringReturn:          "",
			statusCode:               http.StatusUnauthorized,
			getShortenedURLCallCount: 0,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
		{
			name:            "error wile shortening",
			getStringReturn: "userId",
			serviceReturn: getShortenURLReturn{
				result: nil,
				error:  errors.New("some error"),
			},
			statusCode:               http.StatusBadRequest,
			getShortenedURLCallCount: 1,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
		{
			name:            "urls not found",
			getStringReturn: "userId",
			serviceReturn: getShortenURLReturn{
				result: nil,
				error:  nil,
			},
			statusCode:               http.StatusNoContent,
			getShortenedURLCallCount: 1,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockShortenService(ctrl)
			mockService.EXPECT().GetShortenedURLByUserID(tt.getStringReturn).
				Return(tt.serviceReturn.result, tt.serviceReturn.error).Times(tt.getShortenedURLCallCount)

			mockRequestContext := mocks.NewMockRequestContext(ctrl)
			mockRequestContext.EXPECT().GetString("uid").Return(tt.getStringReturn)
			mockRequestContext.EXPECT().AbortWithStatus(tt.statusCode).
				Times(tt.abortWithStatusCallCount)
			mockRequestContext.EXPECT().JSON(tt.statusCode, gomock.Any()).
				Times(tt.jsonCallCount)

			h := &Handler{
				shortener: mockService,
			}
			h.GetUserURL(mockRequestContext)
		})
	}
}
