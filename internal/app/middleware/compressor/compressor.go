package compressor

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const contentType = "Content-Type"
const compressionType = "gzip"

// CompressionMiddleware is a middleware that compresses and decompress data.
func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Принимаем запросы в сжатом формате
		encoding := c.GetHeader("Content-Encoding")
		if strings.Contains(encoding, compressionType) {
			compressReader, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer func() {
				_ = compressReader.Close()
			}()
			c.Request.Body = compressReader
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, compressionType)
		if supportsGzip {
			switch c.Writer.Header().Get(contentType) {
			case "application/json", "text/html":
				compressWriter := newCompressWriter(c.Writer)
				compressWriter.writer.Header().Set(contentType, compressionType)
				defer func() {
					_ = compressWriter.Close()
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
