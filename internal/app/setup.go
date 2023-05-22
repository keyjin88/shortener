package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"go.uber.org/zap"
	"net/http"
)

// API is the Base server instance description
type API struct {
	config        *config.Config
	router        *gin.Engine
	shortener     *service.ShortenService
	urlRepository *storage.URLRepositoryInMem
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
	err := api.configStorage()
	if err != nil {
		return err
	}
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
	router.Use(handlers.CompressionMiddleware())

	//Раскомментировать для перехода на штатный логгер gin
	//router.Use(gin.Logger())
	rootGroup := router.Group("/")
	{
		rootGroup.POST("", handlers.WithLogging(api.handlers.ShortenURLText))
		rootGroup.GET(":id", handlers.WithLogging(api.handlers.GetShortenedURL))
	}
	apiGroup := rootGroup.Group("/api")
	{
		apiGroup.POST("/shorten", handlers.WithLogging(api.handlers.ShortenURLJSON))
	}
	api.router = router
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.shortener, api.config)
}

func (api *API) configStorage() error {
	//передаем в репозиторий только необходимую часть конфига
	api.urlRepository = storage.NewURLRepositoryInMem(api.config.FileStoragePath)
	//пробуем восстановиться из файла
	err := api.urlRepository.RestoreFromFile()
	if err != nil {
		return err
	}
	return nil
}

func (api *API) configService() {
	api.shortener = service.NewShortenService(api.urlRepository)
}
