package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/shortener/internal/app/logger"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		gin.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData       *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый gin.HandlerFunc.
func WithLogging(f func(c RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: c.Writer, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		c.Writer = &lw
		f(c)
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
