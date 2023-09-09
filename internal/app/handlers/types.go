package handlers

import (
	"database/sql/driver"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
)

const key = "uid"
const template = "uid is empty"
const marshalErrorTemplate = "error while marshalling json data: %v"
const shorteningErrorTemplate = "error while shortening url: %v"
const readRequestErrorTemplate = "error while reading request: %v"
const urlAlreadyExist = "URL already exists"

// RequestContext represents a request context and hides gin.Context.
//
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

// ShortenService is interface for shortening URLs.
//
//go:generate mockgen -destination=mocks/shorten_service.go -package=mocks . ShortenService
type ShortenService interface {
	GetShortenedURLByID(id string) (storage.ShortenedURL, error)
	GetShortenedURLByUserID(userID string) ([]storage.UsersURLResponse, error)
	ShortenURL(url string, userID string) (string, error)
	ShortenURLBatch(request storage.ShortenURLBatchRequest, userID string) ([]storage.ShortenURLBatchResponse, error)
	DeleteURLs(req *[]string, userID string) error
}

// Handler is a struct of handler.
type Handler struct {
	shortener ShortenService
	pinger    driver.Pinger
}

// NewHandler returns a new Handler with shortener and pinger configured.
func NewHandler(shortener *service.ShortenService, pinger driver.Pinger) *Handler {
	return &Handler{
		shortener: shortener,
		pinger:    pinger,
	}
}
