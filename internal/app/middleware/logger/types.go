package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type (
	// Берём структуру для хранения сведений об ответе.
	responseData struct {
		status int
		size   int
	}

	// Добавляем реализацию gin.ResponseWriter.
	loggingResponseWriter struct {
		gin.ResponseWriter // встраиваем оригинальный gin.ResponseWriter
		responseData       *responseData
	}
)

// Write writes the response data to the response writer.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, errors.Wrap(err, "failed to write request data")
}

// WriteHeader write a header to response writer.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
