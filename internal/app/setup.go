package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/middleware/compressor"
	loggerMiddleware "github.com/keyjin88/shortener/internal/app/middleware/logger"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
	"github.com/keyjin88/shortener/internal/app/storage/postgres"
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

	logger.Log.Infof("Running server. Address: %s |Base url: %s |DB DSN: %s |Gin release mode: %v |Log level: %s |Filestore path: %s",
		api.config.ServerAddress, api.config.BaseAddress, api.config.DataBaseDSN, api.config.GinReleaseMode, api.config.LogLevel, api.config.FileStoragePath)
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
	api.handlers = handlers.NewHandler(api.shortenService, api.urlRepository)
}

func (api *API) configStorage() {
	if api.config.DataBaseDSN != "" {
		dbPool, err := pgxpool.New(context.Background(), api.config.DataBaseDSN)
		if err != nil {
			logger.Log.Errorf("error while initialising DB Pool: %v", err)
			return
		}
		repository, err := postgres.NewPostgresRepository(dbPool, context.Background())
		if err != nil {
			logger.Log.Errorf("error while initialising DB: %v", err)
			return
		}
		api.urlRepository = repository
	} else if api.config.FileStoragePath != "" {
		repository, err := file.NewURLRepositoryFile(&api.config.FileStoragePath)
		if err != nil {
			logger.Log.Errorf("error while initialising DB: %v", err)
			return
		}
		api.urlRepository = repository
	} else {
		api.urlRepository = inmem.NewURLRepositoryInMem()
	}
}

func (api *API) configService() {
	api.shortenService = service.NewShortenService(api.urlRepository, api.config.BaseAddress)
}
