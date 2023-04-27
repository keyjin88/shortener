package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShortenedURL(t *testing.T) {
	testCases := []struct {
		method       string
		request      string
		expectedCode int
		body         string
		contentType  string
		expectedBody string
	}{
		{method: http.MethodGet, request: "/UUID", expectedCode: http.StatusBadRequest, body: "", expectedBody: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			GetShortenedURL(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestShortenURL(t *testing.T) {
	testCases := []struct {
		method       string
		request      string
		expectedCode int
		body         string
		contentType  string
		expectedBody string
	}{
		{method: http.MethodPost, request: "/", expectedCode: http.StatusCreated, body: "https://practicum.yandex.ru/",
			expectedBody: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			if tc.method == http.MethodPost {
				r.Header.Set("Content-Type", "text/plain")
				ShortenURL(w, r)
			}
			GetShortenedURL(w, r)

			require.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			//Проверим наличие тела ответа
			require.NotEmpty(t, w.Body)
		})
	}
}
