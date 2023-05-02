package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetShortenedURL(c *gin.Context) {
	id := c.Params.ByName("id")
	originalURL, ok := h.shortener.GetShortenedURLByID(id)
	if !ok {
		logger.Infof("URL not found by id: %s. Error while Api.GetShortenedURLByID()", id)
		c.Status(http.StatusBadRequest)
		_, err := c.Writer.Write([]byte(fmt.Sprintf("URL not found by id: %s", id)))
		if err != nil {
			return
		}
		return
	} else {
		logger.Infof("Запрос на получение URL по id: %s, originalURL: %s", id, originalURL)
		c.Redirect(http.StatusTemporaryRedirect, originalURL)
		return
	}
}
