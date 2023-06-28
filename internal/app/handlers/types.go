package handlers

import (
	"database/sql/driver"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
)

//go:generate mockgen -destination=mocks/get_shortened_url.go -package=mocks . RequestContext
type RequestContext interface {
	String(code int, format string, values ...any)
	Redirect(code int, location string)
	Param(key string) string
	Header(key, value string)
	GetRawData() ([]byte, error)
	ShouldBind(obj any) error
	JSON(code int, obj any)
	FullPath() string
	AbortWithStatus(code int)
	BindJSON(obj any) error
	GetString(key string) (s string)
}

//go:generate mockgen -destination=mocks/shorten_service.go -package=mocks . ShortenService
type ShortenService interface {
	GetShortenedURLByID(id string) (storage.ShortenedURL, error)
	GetShortenedURLByUserID(userID string) ([]storage.UsersURLResponse, error)
	ShortenURL(url string, userID string) (string, error)
	ShortenURLBatch(request storage.ShortenURLBatchRequest, userID string) ([]storage.ShortenURLBatchResponse, error)
	DeleteURLs(req *[]string, userID string) error
}

type Handler struct {
	shortener ShortenService
	pinger    driver.Pinger
}

func NewHandler(shortener *service.ShortenService, pinger driver.Pinger) *Handler {
	return &Handler{
		shortener: shortener,
		pinger:    pinger,
	}
}
