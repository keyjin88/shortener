package handlers

import (
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	shortener *service.ShortenService
	config    *config.Config
}

func NewHandler(shortener *service.ShortenService, config *config.Config) *Handler {
	return &Handler{shortener: shortener, config: config}
}

var (
	logger = logrus.New()
)
