package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"net/http"
	"testing"
)

type shortenUrlReturn struct {
	result string
	error  error
}

type getRowDataReturn struct {
	result []byte
	error  error
}

func TestHandler_ShortenURLWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name               string
		url                string
		getRowDataReturn   getRowDataReturn
		serviceReturn      shortenUrlReturn
		expectedCode       int
		expectedBody       string
		expectedStringCall int
		shortenStringCall  int
	}{
		{
			name:               "Create successfully",
			url:                "https://ya.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("https://ya.ru"), error: nil},
			serviceReturn:      shortenUrlReturn{result: "SHORTURL", error: nil},
			expectedCode:       http.StatusCreated,
			expectedBody:       "/SHORTURL",
			expectedStringCall: 1,
			shortenStringCall:  1,
		},
		{
			name:               "Invalid request body",
			url:                "ww.bad url",
			getRowDataReturn:   getRowDataReturn{result: nil, error: errors.New("error from GetRowData()")},
			serviceReturn:      shortenUrlReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Invalid request body.",
			expectedStringCall: 1,
			shortenStringCall:  0,
		},
		{
			name:               "Invalid url string.",
			url:                "htt=bad_url.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("htt=bad_url.ru"), error: nil},
			serviceReturn:      shortenUrlReturn{result: "", error: nil},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Invalid url string.",
			expectedStringCall: 1,
			shortenStringCall:  0,
		},
		{
			name:               "Trouble while shortening url.",
			url:                "https://ya.ru",
			getRowDataReturn:   getRowDataReturn{result: []byte("https://ya.ru"), error: nil},
			serviceReturn:      shortenUrlReturn{result: "", error: errors.New("error from shorten service")},
			expectedCode:       http.StatusBadRequest,
			expectedBody:       "Trouble while shortening url.",
			expectedStringCall: 1,
			shortenStringCall:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockShortenService(ctrl)
			mockService.EXPECT().ShortenString(tt.url).
				Times(tt.shortenStringCall).
				Return(tt.serviceReturn.result, tt.serviceReturn.error)

			mockRequestContext := mocks.NewMockRequestContext(ctrl)
			mockRequestContext.EXPECT().String(tt.expectedCode, tt.expectedBody).
				Times(tt.expectedStringCall)
			mockRequestContext.EXPECT().Header("Content-Type", "text/plain").
				Times(1)
			mockRequestContext.EXPECT().GetRawData().
				Return(tt.getRowDataReturn.result, tt.getRowDataReturn.error)

			h := &Handler{
				shortener: mockService,
				config:    config.NewConfig(),
			}
			h.ShortenURL(mockRequestContext)

		})
	}
}