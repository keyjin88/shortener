package handlers

import (
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/service"
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
}

//go:generate mockgen -destination=mocks/shorten_srvice.go -package=mocks . ShortenService
type ShortenService interface {
	GetShortenedURLByID(id string) (string, bool)
	ShortenURL(url string) (string, error)
}

type Handler struct {
	shortener ShortenService
	config    *config.Config
}

func NewHandler(shortener *service.ShortenService, config *config.Config) *Handler {
	return &Handler{shortener: shortener, config: config}
}
