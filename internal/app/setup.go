package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/config"
	"github.com/keyjin88/shortener/internal/app/handlers"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/middleware/auth"
	"github.com/keyjin88/shortener/internal/app/middleware/compressor"
	loggerMiddleware "github.com/keyjin88/shortener/internal/app/middleware/logger"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/keyjin88/shortener/internal/app/storage/file"
	"github.com/keyjin88/shortener/internal/app/storage/inmem"
	"github.com/keyjin88/shortener/internal/app/storage/postgres"
	"github.com/pkg/errors"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const duration = 30 * time.Second

// API is the Base server instance description.
type API struct {
	config         *config.Config
	router         *gin.Engine
	shortenService *service.ShortenService
	urlRepository  service.URLRepository
	handlers       *handlers.Handler
}

// New is API constructor: build base API instance.
func New() *API {
	return &API{
		config: config.NewConfig(),
	}
}

// Start http server and configure it.
func (api *API) Start() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := logger.Initialize(api.config.LogLevel); err != nil {
		return fmt.Errorf("error while initialize logger: %w", err)
	}
	api.config.InitConfig()
	api.configStorage()
	defer api.urlRepository.Close()
	api.configService()
	api.configHandlers()
	api.configRouter()

	if api.config.HTTPSEnable {
		api.configCerts()
	}

	srv := &http.Server{
		Addr:    api.config.ServerAddress,
		Handler: api.router,
	}

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Errorf("Error while listening and serving: %v", err)
		}
	}()
	logger.Log.Infof("Server started")
	// Ожидаем получения сигнала остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Log.Infof("Stop signal received")
	// Отменяем контекст для graceful shutdown
	cancel()
	// Устанавливаем таймаут для graceful shutdown
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), duration)
	defer cancelShutdown()

	// Останавливаем HTTP-сервер
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Log.Errorf("Error shutting down: %v", err)
	}
	// Закрываем соединения с БД
	api.urlRepository.Close()
	logger.Log.Infof("Server stopped")
	return nil
}

func (api *API) configCerts() {
	// Настройка TLS-сертификата и ключа.
	cert, err := tls.LoadX509KeyPair(api.config.PathToCert, api.config.PathToKey)
	if err != nil {
		logger.Log.Errorf("failed to load TLS certificates: %s", err)
	}

	// Создание HTTP-сервера с поддержкой TLS.
	httpServer := &http.Server{
		Addr:    ":443",
		Handler: api.router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	// Запуск HTTP-сервера.
	if err := httpServer.ListenAndServeTLS("", ""); err != nil {
		logger.Log.Errorf("Failed to start HTTPS server: %v", err)
	}
}

func (api *API) configRouter() {
	if api.config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	pprof.Register(router)
	router.Use(auth.AuthenticationMiddleware(&api.config.SecretKey))
	router.Use(compressor.CompressionMiddleware())
	router.Use(loggerMiddleware.LoggingMiddleware())
	// Раскомментировать для перехода на штатный логгер gin
	// router.Use(gin.Logger()).
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
		apiGroup.GET("/user/urls", func(c *gin.Context) { api.handlers.GetUserURL(c) })
		apiGroup.DELETE("/user/urls", func(c *gin.Context) { api.handlers.DeleteURLs(c) })
	}
	api.router = router
}

func (api *API) configHandlers() {
	api.handlers = handlers.NewHandler(api.shortenService, api.urlRepository)
}

func (api *API) configStorage() {
	const template = "error while initialising DB: %v"
	switch {
	case api.config.DataBaseDSN != "":
		dbPool, err := pgxpool.New(context.Background(), api.config.DataBaseDSN)
		if err != nil {
			logger.Log.Errorf("error while initialising DB Pool: %v", err)
			return
		}
		ch := make(chan storage.UserURLs)
		repository, err := postgres.NewPostgresRepository(dbPool, context.Background(), ch)
		if err != nil {
			logger.Log.Errorf(template, err)
			return
		}
		api.urlRepository = repository
	case api.config.FileStoragePath != "":
		repository, err := file.NewURLRepositoryFile(&api.config.FileStoragePath)
		if err != nil {
			logger.Log.Errorf(template, err)
			return
		}
		api.urlRepository = repository
	default:
		api.urlRepository = inmem.NewURLRepositoryInMem()
	}
}

func (api *API) configService() {
	api.shortenService = service.NewShortenService(api.urlRepository, api.config.BaseAddress)
}
