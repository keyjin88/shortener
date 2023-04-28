package helpers

import (
	"github.com/gin-gonic/gin"
)

type Message struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func RespondJSON(w *gin.Context, statusCode int, data interface{}) {
	var message Message

	message.StatusCode = statusCode
	message.Data = data
	w.JSON(statusCode, message)
}
