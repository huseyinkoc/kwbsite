package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

// Yeni JWT Claims yapısı
type Claims struct {
	UserID            string   `json:"userID"`
	Username          string   `json:"username"`
	Email             string   `json:"email"`
	PreferredLanguage string   `json:"preferred_language"`
	Roles             []string `json:"roles"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your_secret_key") // JWT token için gizli anahtar

// AuthMiddleware JWT doğrulama middleware'i
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization başlığını kontrol et
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// 'Bearer ' kısmını temizle
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Token'ı doğrula
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Kullanıcı bilgilerini context'e ekle
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("preferred_language", claims.PreferredLanguage)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	// Yeni Tollbooth v7 sürümüne göre limiter tanımlama
	limiter := tollbooth.NewLimiter(3, nil) // 5 dakika içinde 3 istek limiti

	// Tollbooth v7 sürümünde zaman ayarlamaları **otomatik olarak IP başına hesaplanır**
	// Dolayısıyla ekstra bir TTL metodu çağırmaya gerek yok!

	return tollbooth_gin.LimitHandler(limiter)
}
