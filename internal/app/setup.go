package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/middleware/compressor"
	loggerMiddleware "github.com/keyjin88/shortener/internal/app/middleware/logger"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
	"github.com/keyjin88/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
	"net/http"
)

// API is the Base server instance description
type API struct {
	config         *config.Config
	router         *gin.Engine
	shortenService *service.ShortenService
	urlRepository  service.URLRepository
	handlers       *handlers.Handler
}

// New is API constructor: build base API instance
func New() *API {
	return &API{
		config: config.NewConfig(),
	}
}

// Start http server and configure it
func (api *API) Start() error {
	if err := logger.Initialize(api.config.LogLevel); err != nil {
		return err
	}
	api.config.InitConfig()
	api.configStorage()
	api.configService()
	api.configureHandlers()
	api.setupRouter()

	defer api.urlRepository.Close()

	logger.Log.Debug("Running server",
		zap.String("Address", api.config.ServerAddress),
		zap.String("Base addres", api.config.BaseAddress),
		zap.String("DB DSN", api.config.DataBaseDSN),
		zap.Bool("Gin release mode", api.config.GinReleaseMode),
		zap.String("Log level", api.config.LogLevel),
		zap.String("Filestore path", api.config.LogLevel),
	)
	return http.ListenAndServe(api.config.ServerAddress, api.router)
}

func (api *API) setupRouter() {
	if api.config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(compressor.CompressionMiddleware())
	router.Use(loggerMiddleware.LoggingMiddleware())
	//Раскомментировать для перехода на штатный логгер gin
	//router.Use(gin.Logger())
	rootGroup := router.Group("/")
	{
		rootGroup.POST("", func(c *gin.Context) { api.handlers.ShortenURLText(c) })
		rootGroup.GET(":id", func(c *gin.Context) { api.handlers.GetShortenedURL(c) })
		rootGroup.GET("ping", func(c *gin.Context) { api.handlers.DBPing(c) })
	}
	apiGroup := rootGroup.Group("/api")
	{
		apiGroup.POST("/shorten", func(c *gin.Context) { api.handlers.ShortenURLJSON(c) })
		apiGroup.POST("/shorten/batch", func(c *gin.Context) { api.handlers.ShortenURLBatch(c) })
	}
	api.router = router
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.shortenService, api.config)
}

func (api *API) configStorage() {
	if api.config.DataBaseDSN != "" {
		repository, err := postgres.InitPgRepository(context.Background(), api.config.DataBaseDSN)
		if err != nil {
			logger.Log.Errorf("error while initialising DB: %v", err)
			return
		}
		api.urlRepository = repository
	} else {
		//передаем в репозиторий только необходимую часть конфига
		api.urlRepository = inmem.NewURLRepositoryInMem()
		//пробуем восстановиться из файла
		if api.config.FileStoragePath != "" {
			data, err := file.RestoreFromFile(api.config.FileStoragePath)
			if err != nil {
				//Логируем ошибку и продолжаем работу
				logger.Log.Errorf("error while restoring DB from file: %v", err)
				return
			}
			for _, shortenedURL := range data {
				_, err := api.urlRepository.Save(shortenedURL.ShortURL, shortenedURL.OriginalURL)
				if err != nil {
					return
				}
			}
		}
	}
}

func (api *API) configService() {
	api.shortenService = service.NewShortenService(api.urlRepository, api.config.FileStoragePath)
}
