package api

import (
	"github.com/gorilla/mux"
	"github.com/keyjin88/shortener/internal/app/service"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

// API is the Base server instance description
type API struct {
	logger    *logrus.Logger
	router    *mux.Router
	shortener *service.ShortenService
	storage   *storage.Storage
}

// New is API constructor: build base API instance
func New() *API {
	return &API{
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start http server and configure it
func (api *API) Start() error {
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	api.logger.Info("logger configured successfully.")
	api.configureRouterField()
	api.logger.Info("router configured successfully.")

	api.logger.Info("starting api server at port: ", "8080")
	//На этапе валидного завршениея стратуем http-сервер
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

// Конфигурируем router API
func (api *API) configureRouterField() {
	api.router.HandleFunc("/", ShortenURL).Methods(http.MethodPost)
	api.router.HandleFunc("/{id}", GetShortenedURL).Methods(http.MethodGet)
}
