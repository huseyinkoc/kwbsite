package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// CSRF token'leri saklamak için eşzamanlı map
var csrfTokens sync.Map // username -> csrfToken

// CSRF token oluşturma
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// CSRF token doğrulama
func ValidateCSRFToken(username, csrfToken string) bool {
	value, ok := csrfTokens.Load(username)
	if !ok {
		return false
	}
	return value == csrfToken
}

// CSRF token saklama
func StoreCSRFToken(username, csrfToken string) {
	csrfTokens.Store(username, csrfToken)
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			username := c.GetString("username") // Kullanıcı bilgisi
			csrfToken := c.GetHeader("X-CSRF-Token")

			if csrfToken == "" || !ValidateCSRFToken(username, csrfToken) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or missing CSRF token"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
