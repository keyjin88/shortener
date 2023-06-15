package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"github.com/keyjin88/shortener/internal/app/logger"
	"net/http"
	"testing"
)

type shortenURLReturn struct {
	result string
	error  error
}

type getRowDataReturn struct {
	result []byte
	error  error
}

func TestHandler_ShortenURLWithMock(t *testing.T) {
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
		expectedBody       string
		expectedStringCall int
		shortenStringCall  int
		getStringCallCount int
	}{
		{
			name:               "Save successfully",
			url:                "https://ya.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("https://ya.ru"), error: nil},
			serviceReturn:      shortenURLReturn{result: "/SHORTURL", error: nil},
			expectedCode:       http.StatusCreated,
			expectedBody:       "/SHORTURL",
			expectedStringCall: 1,
			shortenStringCall:  1,
			getStringCallCount: 1,
		},
		{
			name:               "Invalid request body",
			url:                "ww.bad url",
			getRowDataReturn:   getRowDataReturn{result: nil, error: errors.New("error from GetRowData()")},
			serviceReturn:      shortenURLReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Invalid request body.",
			expectedStringCall: 1,
			shortenStringCall:  0,
		},
		{
			name:               "Invalid url string.",
			url:                "htt=bad_url.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("htt=bad_url.ru"), error: nil},
			serviceReturn:      shortenURLReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Invalid url string.",
			expectedStringCall: 1,
			shortenStringCall:  0,
			getStringCallCount: 0,
		},
		{
			name:               "Trouble while shortening url.",
			url:                "https://ya.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("https://ya.ru"), error: nil},
			serviceReturn:      shortenURLReturn{result: "", error: errors.New("error from shorten service")},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Trouble while shortening url.",
			expectedStringCall: 1,
			shortenStringCall:  1,
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
			mockRequestContext.EXPECT().Header("Content-Type", "text/plain").
				Times(1)
			mockRequestContext.EXPECT().GetRawData().
				Return(tt.getRowDataReturn.result, tt.getRowDataReturn.error)
			mockRequestContext.EXPECT().GetString(gomock.Any()).
				Times(tt.getStringCallCount)

			h := &Handler{
				shortener: mockService,
			}
			h.ShortenURLText(mockRequestContext)
		})
	}
}
