package handlers

import (
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
}

//go:generate mockgen -destination=mocks/shorten_service.go -package=mocks . ShortenService
type ShortenService interface {
	GetShortenedURLByID(id string) (string, error)
	ShortenURL(url string) (string, error)
	ShortenURLBatch(request storage.ShortenURLBatchRequest) ([]storage.ShortenURLBatchResponse, error)
}

type Handler struct {
	shortener ShortenService
	config    *Config
}

type Config struct {
	DataBaseDSN string
}

func NewHandler(shortener *service.ShortenService, dataBaseDSN string) *Handler {
	return &Handler{shortener: shortener, config: &Config{DataBaseDSN: dataBaseDSN}}
}
