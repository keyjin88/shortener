package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"time"
)

// LoggingMiddleware добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый gin.HandlerFunc.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: c.Writer, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		c.Writer = &lw
		c.Next()
		duration := time.Since(start)
		logger.Log.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
}
