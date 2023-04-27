package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// Вспомогательная структура для формирования сообщений
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func (api *API) ShortenURL(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	requestBytes, err := io.ReadAll(request.Body)
	if err != nil {
		api.logger.Error("Invalid url string. Error while Api.ShortenURL() :", err)
		msg := Message{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	shortenString, err := api.shortener.ShortenString(string(requestBytes))
	if err != nil {
		api.logger.Error("Trouble while shortening url. Error while shortener.ShortenString() :", err)
		msg := Message{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	api.logger.Infof("Запрос на сокращение URL: %s", string(requestBytes))
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write([]byte("http://localhost:8080/" + shortenString))
	if err != nil {
		return
	}
}

func (api *API) GetShortenedURL(writer http.ResponseWriter, request *http.Request) {
	//Логируем момент начала обработки запроса
	api.logger.Info("Get All Articles GET /api/v1/articles")
	id, ok := mux.Vars(request)["id"]
	if !ok {
		api.logger.Info("Invalid id. Error while Api.GetShortenedURL() :", mux.Vars(request)["id"])
		msg := Message{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid id",
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	originalUrl, ok, _ := api.shortener.GetShortenedUrl(id)
	if !ok {
		api.logger.Infof("URL not found by id: %s. Error while Api.GetShortenedURL()", id)
		msg := Message{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("URL not found by id: %s. Error while Api.GetShortenedURL()", id),
			IsError:    true,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(msg)
		return
	} else {
		api.logger.Infof("Запрос на получение URL по id: %s", id)
		writer.Header().Set("Location", originalUrl)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
