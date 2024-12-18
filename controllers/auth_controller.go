package controllers

import (
	"admin-panel/middlewares"
	"admin-panel/services"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("your_secret_key") // JWT için gizli anahtar

// JWT Claims yapısı
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Kullanıcıyı veritabanından al
	user, err := services.GetUserByUsername(input.Username)
	if err != nil {
		log.Printf("Login failed: User %s not found", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Şifre doğrulaması
	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		log.Printf("Login failed: Incorrect password for user %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWT token oluştur
	expirationTime := time.Now().Add(24 * time.Hour) // 1 gün geçerli
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Login failed: Unable to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// CSRF token oluştur
	csrfToken, err := middlewares.GenerateCSRFToken()
	if err != nil {
		log.Println("Login failed: Unable to generate CSRF token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
		return
	}

	// CSRF token'i oturum bazlı saklama
	middlewares.StoreCSRFToken(input.Username, csrfToken)

	// Yanıtı döndür
	log.Printf("Login successful: User %s logged in", input.Username)

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"csrf_token": csrfToken,
		"message":    "Login successful",
	})
}
