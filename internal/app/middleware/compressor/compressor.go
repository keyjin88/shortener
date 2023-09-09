package compressor

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const contentType = "Content-Type"
const compressType = "gzip"

// CompressionMiddleware is a middleware that compresses and decompress data.
func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error // Объявление переменной err
		encoding := c.GetHeader("Content-Encoding")
		if strings.Contains(encoding, compressType) {
			compressReader, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer func() {
				if closeErr := compressReader.Close(); closeErr != nil && err == nil {
					err = closeErr
				}
			}()
			c.Request.Body = compressReader
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, compressType)
		if supportsGzip {
			switch c.Writer.Header().Get(contentType) {
			case "application/json", "text/html":
				compressWriter := newCompressWriter(c.Writer)
				compressWriter.writer.Header().Set(contentType, compressType)
				defer func() {
					if closeErr := compressWriter.Close(); closeErr != nil && err == nil {
						err = closeErr
					}
				}()
			}
			// Продолжаем обработку запроса
			c.Next()
		} else {
			// Клиент не поддерживает сжатие, просто продолжаем обработку запроса
			c.Next()
		}
	}
}
