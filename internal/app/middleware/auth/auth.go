package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/keyjin88/shortener/internal/app/logger"
	"io"
	"net/http"
	"time"
)

const errorCreatingCipher = "error while creating cipher"

// Cookie represents a cookie used for authorization.
type Cookie struct {
	encryptedCookie string
	expiration      time.Time
	uid             string
}

// AuthenticationMiddleware is a middleware to authenticate users.
func AuthenticationMiddleware(secretKey *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const auth = "auth"
		const generateCookieErrorTemplate = "error while generate new cookie: %v"
		cookie, err := c.Cookie(auth)
		const key = "uid"
		if err != nil {
			// Куки не существует, выдаём новую
			generatedCookie, err := generateCookie(secretKey)
			if err != nil {
				withError := c.AbortWithError(http.StatusBadRequest, err)
				if withError != nil {
					logger.Log.Infof(generateCookieErrorTemplate, err)
				}
				return
			}
			c.SetCookie(auth, generatedCookie.encryptedCookie, int(generatedCookie.expiration.Unix()), "/", "", false, true)
			c.Set(key, generatedCookie.uid)
			c.Next()
			return
		}

		// Куки существует, проверяем подлинность
		uid, err := decryptCookie(cookie, []byte(*secretKey))
		if err != nil {
			// Куки не проходит проверку, выдаём новую
			generatedCookie, err := generateCookie(secretKey)
			if err != nil {
				withError := c.AbortWithError(http.StatusBadRequest, err)
				if withError != nil {
					logger.Log.Infof(generateCookieErrorTemplate, err)
				}
				return
			}
			c.SetCookie(auth, generatedCookie.encryptedCookie, int(generatedCookie.expiration.Unix()), "/", "", false, true)
		}
		// Передаём уникальный идентификатор пользователя в следующий middleware/handler
		c.Set(key, uid)
		c.Next()
	}
}

func encryptCookie(cookieValue string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", errors.Wrap(err, errorCreatingCipher)
	}
	ciphertext := make([]byte, aes.BlockSize+len(cookieValue))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "error while reading cookie")
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(cookieValue))
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptCookie(cipherText string, secretKey []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.Wrap(err, "error while decode string")
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", errors.Wrap(err, errorCreatingCipher)
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	initVector := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, initVector)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

func generateCookie(secretKey *string) (Cookie, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return Cookie{}, errors.Wrap(err, "error while generate UUID")
	}
	uid := newUUID.String()
	encryptedCookie, err := encryptCookie(uid, []byte(*secretKey))
	if err != nil {
		return Cookie{}, errors.Wrap(err, "error while encrypting cookie")
	}
	const oneDay = 24 * time.Hour
	expiration := time.Now().Add(oneDay)
	return Cookie{encryptedCookie, expiration, uid}, nil
}
