package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API is the Base server instance description
type API struct {
	config        *config.Config
	logger        *logrus.Logger
	router        *gin.Engine
	shortener     *service.ShortenService
	urlRepository *storage.URLRepositoryInMem
	handlers      *handlers.Handler
}

// New is API constructor: build base API instance
func New() *API {
	return &API{
		logger:    logrus.New(),
		config:    config.NewConfig(),
		shortener: service.NewShortenService(storage.NewURLRepositoryInMem()),
	}
}

// Start http server and configure it
func (api *API) Start() error {
	api.logger.Info("read configs")
	api.config.InitConfig()

	api.logger.Info("configure handlers")
	api.configureHandlers()

	api.logger.Info("setting up router")
	api.setupRouter()

	if err := api.configureLoggerField(); err != nil {
		return err
	}
	api.logger.Info("logger configured successfully.")
	serverAddress := api.config.ServerAddress
	api.logger.Info("starting api server at: ", serverAddress)
	return http.ListenAndServe(serverAddress, api.router)
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

func (api *API) setupRouter() {
	router := gin.Default()
	//В Gin принято группировать ресурсы
	apiV1Group := router.Group("/")
	apiV1Group.POST("/", api.handlers.ShortenURL)
	apiV1Group.GET("/:id", api.handlers.GetShortenedURL)
	api.router = router
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.shortener, api.config)
}