package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/middleware/compressor"
	logger2 "github.com/keyjin88/shortener/internal/app/middleware/logger"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
	"go.uber.org/zap"
	"net/http"
)

// API is the Base server instance description
type API struct {
	config        *config.Config
	router        *gin.Engine
	shortener     *service.ShortenService
	urlRepository *inmem.URLRepositoryInMem
	handlers      *handlers.Handler
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

	logger.Log.Infow("Running server", zap.String("address", api.config.ServerAddress))
	return http.ListenAndServe(api.config.ServerAddress, api.router)
}

func (api *API) setupRouter() {
	if api.config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(compressor.CompressionMiddleware())
	router.Use(logger2.LoggingMiddleware())
	//Раскомментировать для перехода на штатный логгер gin
	//router.Use(gin.Logger())
	rootGroup := router.Group("/")
	{
		rootGroup.POST("", func(c *gin.Context) { api.handlers.ShortenURLText(c) })
		rootGroup.GET(":id", func(c *gin.Context) { api.handlers.GetShortenedURL(c) })
	}
	apiGroup := rootGroup.Group("/api")
	{
		apiGroup.POST("/shorten", func(c *gin.Context) { api.handlers.ShortenURLJSON(c) })
	}
	api.router = router
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.shortener, api.config)
}

func (api *API) configStorage() {
	//передаем в репозиторий только необходимую часть конфига
	api.urlRepository = inmem.NewURLRepositoryInMem(api.config.FileStoragePath)
	//пробуем восстановиться из файла
	if api.config.FileStoragePath != "" {
		data, err := file.RestoreFromFile(api.config.FileStoragePath)
		if err != nil {
			//Логируем ошибку и продолжаем работу
			logger.Log.Errorf("error while restoring DB from file: %v", err)
			return
		}
		api.urlRepository.RestoreData(data)
	}
}

func (api *API) configService() {
	api.shortener = service.NewShortenService(api.urlRepository)
}
