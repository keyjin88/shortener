package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
	"testing"
)

func TestHandler_DeleteURLs(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                     string
		serviceReturn            error
		serviceCallCount         int
		status                   int
		abortWithStatusCallCount int
		jsonCallCount            int
		bindJsonReturn           error
		bindJsonCallCount        int
		getStringReturn          string
	}{
		{
			name:                     "success",
			getStringReturn:          "userID",
			serviceReturn:            nil,
			serviceCallCount:         1,
			bindJsonReturn:           nil,
			bindJsonCallCount:        1,
			status:                   http.StatusAccepted,
			abortWithStatusCallCount: 0,
			jsonCallCount:            1,
		},
		{
			name:                     "empty userID",
			getStringReturn:          "",
			serviceReturn:            nil,
			serviceCallCount:         0,
			bindJsonReturn:           nil,
			bindJsonCallCount:        0,
			status:                   http.StatusUnauthorized,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
		{
			name:                     "error binding json",
			getStringReturn:          "userID",
			serviceReturn:            nil,
			serviceCallCount:         0,
			bindJsonReturn:           errors.New("some error"),
			bindJsonCallCount:        1,
			status:                   http.StatusBadRequest,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
		{
			name:                     "error while deleting urls",
			getStringReturn:          "userID",
			serviceReturn:            errors.New("some error"),
			serviceCallCount:         1,
			bindJsonReturn:           nil,
			bindJsonCallCount:        1,
			status:                   http.StatusInternalServerError,
			abortWithStatusCallCount: 1,
			jsonCallCount:            0,
		},
	}
	for _, tt := range tests {
		mockService := mocks.NewMockShortenService(ctrl)
		mockService.EXPECT().DeleteURLs(gomock.Any(), gomock.Any()).
			Return(tt.serviceReturn).
			Times(tt.serviceCallCount)

		mockRequestContext := mocks.NewMockRequestContext(ctrl)
		mockRequestContext.EXPECT().GetString(gomock.Any()).Return(tt.getStringReturn).Times(1)
		mockRequestContext.EXPECT().AbortWithStatus(tt.status).Times(tt.abortWithStatusCallCount)
		mockRequestContext.EXPECT().JSON(tt.status, gomock.Any()).Times(tt.jsonCallCount)
		mockRequestContext.EXPECT().BindJSON(gomock.Any()).Return(tt.bindJsonReturn).Times(tt.bindJsonCallCount)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				shortener: mockService,
			}
			h.DeleteURLs(mockRequestContext)
		})
	}
}
