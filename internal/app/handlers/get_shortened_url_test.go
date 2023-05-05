package handlers

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers/mocks"
	"net/http"
	"testing"
)

type getShortenURLReturn struct {
	result string
	ok     bool
}

func TestHandler_GetShortenedURLWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name                 string
		shortenURL           string
		originalURL          string
		serviceReturn        getShortenURLReturn
		expectedStringCall   int
		expectedRedirectCall int
		expectedCode         int
		expectedBody         string
	}{
		{
			name:                 "Get successfully",
			shortenURL:           "ShortenURL",
			originalURL:          "https://www.test.ru",
			serviceReturn:        getShortenURLReturn{result: "https://www.test.ru", ok: true},
			expectedRedirectCall: 1,
			expectedStringCall:   0,
			expectedCode:         http.StatusTemporaryRedirect,
		},
		{
			name:                 "URL not found",
			shortenURL:           "ShortenURL",
			serviceReturn:        getShortenURLReturn{result: "", ok: false},
			expectedCode:         http.StatusBadRequest,
			expectedRedirectCall: 0,
			expectedStringCall:   1,
			expectedBody:         "URL not found by id: %s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockShortenService(ctrl)
			mockService.EXPECT().GetShortenedURLByID(tt.shortenURL).
				Return(tt.serviceReturn.result, tt.serviceReturn.ok)

			mockRequestContext := mocks.NewMockRequestContext(ctrl)
			mockRequestContext.EXPECT().Param("id").
				Return(tt.shortenURL)
			mockRequestContext.EXPECT().String(tt.expectedCode, fmt.Sprintf(tt.expectedBody, tt.shortenURL)).
				Times(tt.expectedStringCall)
			mockRequestContext.EXPECT().Redirect(tt.expectedCode, tt.originalURL).
				Times(tt.expectedRedirectCall)

			h := &Handler{
				shortener: mockService,
				config:    config.NewConfig(),
			}
			h.GetShortenedURL(mockRequestContext)
		})
	}
}