package api

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API is the Base server instance description
type API struct {
	logger    *logrus.Logger
	router    *gin.Engine
	shortener *service.ShortenService
	storage   *storage.Storage
}

// New is API constructor: build base API instance
func New() *API {
	return &API{
		logger: logrus.New(),
		router: SetupRouter(),
	}
}

// Start http server and configure it
func (api *API) Start() error {
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	api.logger.Info("logger configured successfully.")
	api.logger.Info("starting api server at port: ", "8080")
	return http.ListenAndServe("localhost:8080", api.router)
}

// Конфигурируем logger
func (api *API) configureLoggerField() error {
	logLevel, err := logrus.ParseLevel("debug")
	if err != nil {
		return err
	}
	api.logger.SetLevel(logLevel)
	return nil
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	//В Gin принято группировать ресурсы
	apiV1Group := router.Group("/")
	{
		apiV1Group.POST("/", ShortenURL)
		apiV1Group.GET("/:id", GetShortenedURL)
	}
	return router
}
