package controllers

import (
	"admin-panel/middlewares"
	"admin-panel/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = []byte("your_secret_key") // JWT token için gizli anahtar

// JWT Claims yapısı
type Claims struct {
	UserID            string   `json:"userID"`             // Kullanıcı ID'si
	Username          string   `json:"username"`           // Kullanıcı adı
	Email             string   `json:"email"`              // E-posta adresi
	PreferredLanguage string   `json:"preferred_language"` // Dil tercihi
	Roles             []string `json:"roles"`
	jwt.StandardClaims
}

// LoginHandler authenticates a user
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginByUsername true "User login credentials"
// @Success 200 {object} map[string]interface{} "JWT token and user details"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /svc/auth/login-by-username [post]
func LoginByUsernameHandler(c *gin.Context) {
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

	preferredLanguage := user.PreferredLanguage
	if preferredLanguage == "" {
		preferredLanguage = "en" // Varsayılan dil
	}

	// JWT token oluştur
	expirationTime := time.Now().Add(24 * time.Hour) // 1 gün geçerli
	claims := &Claims{
		UserID:            user.ID.Hex(),
		Username:          user.Username,
		Email:             user.Email,
		Roles:             user.Roles,
		PreferredLanguage: preferredLanguage, // Dil tercihini ekledik
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
		"user": gin.H{
			"name":      user.Name,
			"surname":   user.Surname,
			"full_name": fmt.Sprintf("%s %s", user.Name, user.Surname),
		},
	})
}

// LoginHandler authenticates a user
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginByEmail true "User login credentials"
// @Success 200 {object} map[string]interface{} "JWT token and user details"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /svc/auth/login-by-email [post]
func LoginByEmailHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Kullanıcıyı veritabanından al
	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		log.Printf("Login failed: User %s not found", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Şifre doğrulaması
	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		log.Printf("Login failed: Incorrect password for user %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	preferredLanguage := user.PreferredLanguage
	if preferredLanguage == "" {
		preferredLanguage = "en" // Varsayılan dil
	}

	// JWT token oluştur
	expirationTime := time.Now().Add(24 * time.Hour) // 1 gün geçerli
	claims := &Claims{
		UserID:            user.ID.Hex(),
		Username:          user.Username,
		Email:             user.Email,
		Roles:             user.Roles,
		PreferredLanguage: preferredLanguage, // Dil tercihini ekledik
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
	middlewares.StoreCSRFToken(input.Email, csrfToken)

	// Yanıtı döndür
	log.Printf("Login successful: User %s logged in", input.Email)

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"csrf_token": csrfToken,
		"message":    "Login successful",
		"user": gin.H{
			"name":      user.Name,
			"surname":   user.Surname,
			"full_name": fmt.Sprintf("%s %s", user.Name, user.Surname),
		},
	})
}

// LoginHandler authenticates a user
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginByPhone true "User login credentials"
// @Success 200 {object} map[string]interface{} "JWT token and user details"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /svc/auth/login-by-phone [post]
func LoginByPhoneHandler(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Kullanıcıyı telefon numarasıyla al
	user, err := services.GetUserByPhone(input.PhoneNumber)
	if err != nil {
		log.Printf("Login failed: User %s not found", input.PhoneNumber)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Şifre doğrulama
	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		log.Printf("Login failed: Incorrect password for user %s", input.PhoneNumber)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWT Token oluştur
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
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

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"message": "Login successful",
	})
}

// SendVerificationEmailHandler sends a verification email to the user
// @Summary Send verification email
// @Description Sends a verification email to a specific user
// @Tags Authentication
// @Param userID path string true "User ID"
// @Success 200 {object} map[string]interface{} "Verification email sent"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 500 {object} map[string]interface{} "Failed to send verification email"
// @Router /svc/auth/send-verification/{userID} [post]
func SendVerificationEmailHandler(c *gin.Context) {
	userID := c.Param("userID")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	token, err := services.GenerateEmailVerificationToken(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	err = services.SendVerificationEmail(c.Request.Context(), objectID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// VerifyEmailHandler verifies a user's email
// @Summary Verify email
// @Description Verifies a user's email using a token
// @Tags Authentication
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid or expired token"
// @Router /svc/auth/verify-email [get]
func VerifyEmailHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	err := services.VerifyEmailToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// RequestPasswordResetHandler handles password reset requests
// @Summary Request password reset
// @Description Sends a password reset email to the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body models.RequestPasswordReset true "User email"
// @Success 200 {object} map[string]interface{} "Password reset email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 404 {object} map[string]interface{} "Email not found"
// @Failure 500 {object} map[string]interface{} "Failed to send password reset email"
// @Router /svc/auth/request-password-reset [post]
func RequestPasswordResetHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Kullanıcıyı email ile bulun
	userID, err := services.GetUserIDByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Reset token oluştur
	token, err := services.GeneratePasswordResetToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password reset token"})
		return
	}

	resetURL := "http://localhost:8080/auth/reset-password?token=" + token
	subject := "Password Reset Request"
	body := "Click the link to reset your password: " + resetURL

	err = services.SendEmail([]string{request.Email}, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// ResetPasswordHandler resets a user's password
// @Summary Reset password
// @Description Resets a user's password using a valid reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Password reset token"
// @Param request body models.ResetPasswordRequest true "New password"
// @Success 200 {object} map[string]interface{} "Password updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload or token"
// @Failure 500 {object} map[string]interface{} "Failed to update password"
// @Router /svc/auth/reset-password [post]
func ResetPasswordHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	var request struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Token'ı doğrula
	userID, err := services.VerifyPasswordResetToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Şifreyi güncelle
	err = services.UpdateUserPassword(c.Request.Context(), userID, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Token'ı sil
	_ = services.DeletePasswordResetToken(c.Request.Context(), token)

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
