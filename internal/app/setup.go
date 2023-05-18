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
		config:    config.NewConfig(),
		shortener: service.NewShortenService(storage.NewURLRepositoryInMem()),
	}
}

// Start http server and configure it
func (api *API) Start() error {
	if err := logger.Initialize(api.config.LogLevel); err != nil {
		return err
	}
	api.config.InitConfig()
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
	//Использую стандартный логгер gin. В итоге нужно будет выбрать какой-то один
	//router.Use(gin.Logger())
	rootGroup := router.Group("/")
	rootGroup.POST("api/shorten", handlers.WithLogging(api.handlers.ShortenURL))
	rootGroup.POST("", handlers.WithLogging(api.handlers.ShortenURL))
	rootGroup.GET(":id", handlers.WithLogging(api.handlers.GetShortenedURL))
	api.router = router
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.shortener, api.config)
}
