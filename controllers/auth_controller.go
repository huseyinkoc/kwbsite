package controllers

import (
	"admin-panel/services"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWT Claims yapısı (RegisteredClaims kullanılıyor)
type Claims struct {
	UserID            string   `json:"userID"`
	Username          string   `json:"username"`
	Email             string   `json:"email"`
	PreferredLanguage string   `json:"preferred_language"`
	Roles             []string `json:"roles"`
	jwt.RegisteredClaims
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

	// check lock (username based)
	user, err := services.GetUserByUsername(input.Username)
	if err == nil {
		locked, until, _ := services.IsAccountLockedByEmail(user.Email)
		if locked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account locked", "locked_until": until})
			return
		}
	}

	user, err = services.GetUserByUsernameWithPassword(input.Username)
	if err != nil {
		// do not reveal existence
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		_, _ = services.IncrementFailedLoginByEmail(user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// success -> reset failed attempts
	_ = services.ResetFailedAttempts(user.ID)

	// access token
	tokenString, _, err := services.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// If client already has a refresh cookie and it's valid for this user, reuse it.
	if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
		if uid, ok, _ := services.IsRefreshTokenValid(cookie.Value); ok && uid == user.ID {
			// reuse existing cookie — do not generate/rotate
			c.JSON(http.StatusOK, gin.H{
				"token":      tokenString,
				"expires_in": 15 * 60,
				"message":    "Login successful",
			})
			return
		}
	}

	// otherwise create and set a new refresh token
	refreshPlain, rtExpiry, err := services.GenerateAndStoreRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	setRefreshCookie(c.Writer, refreshPlain, rtExpiry)

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_in": 15 * 60,
		"message":    "Login successful",
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
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// check lock
	locked, until, _ := services.IsAccountLockedByEmail(input.Email)
	if locked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account locked", "locked_until": until})
		return
	}

	user, err := services.GetUserByEmailWithPassword(input.Email)
	if err != nil {
		// do not reveal existence
		_, _ = services.IncrementFailedLoginByEmail(input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		_, _ = services.IncrementFailedLoginByEmail(input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// success -> reset failed attempts
	_ = services.ResetFailedAttempts(user.ID)

	// access token
	tokenString, _, err := services.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// reuse existing cookie if valid
	if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
		if uid, ok, _ := services.IsRefreshTokenValid(cookie.Value); ok && uid == user.ID {
			c.JSON(http.StatusOK, gin.H{
				"token":      tokenString,
				"expires_in": 15 * 60,
				"message":    "Login successful",
			})
			return
		}
	}

	// create new refresh token
	refreshPlain, rtExpiry, err := services.GenerateAndStoreRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	setRefreshCookie(c.Writer, refreshPlain, rtExpiry)

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_in": 15 * 60,
		"message":    "Login successful",
	})
}

// LoginByPhoneHandler için Swagger tanımı
// @Summary Login by phone
// @Description Authenticates user with phone number and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginByPhone true "Phone login credentials"
// @Success 200 {object} map[string]interface{} "JWT token and user details"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 403 {object} map[string]interface{} "Account locked"
// @Failure 500 {object} map[string]interface{} "Server error"
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

	// Get user first to check lock status by email
	user, err := services.GetUserByPhone(input.PhoneNumber)
	if err == nil {
		locked, until, _ := services.IsAccountLockedByEmail(user.Email)
		if locked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account locked", "locked_until": until})
			return
		}
	}

	user, err = services.GetUserByPhoneWithPassword(input.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := services.CheckPassword(user.Password, input.Password); err != nil {
		_, _ = services.IncrementFailedLoginByEmail(user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// success -> reset failed attempts
	_ = services.ResetFailedAttempts(user.ID)

	// access token
	tokenString, _, err := services.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// reuse existing cookie if valid
	if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
		if uid, ok, _ := services.IsRefreshTokenValid(cookie.Value); ok && uid == user.ID {
			c.JSON(http.StatusOK, gin.H{
				"token":      tokenString,
				"expires_in": 15 * 60,
				"message":    "Login successful",
			})
			return
		}
	}

	// create new refresh token
	refreshPlain, rtExpiry, err := services.GenerateAndStoreRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	setRefreshCookie(c.Writer, refreshPlain, rtExpiry)

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_in": 15 * 60,
		"message":    "Login successful",
	})
}

// RefreshHandler için Swagger tanımı
// @Summary Refresh access token
// @Description Rotates refresh token and returns new access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refresh_token cookie string false "Refresh token (cookie)"
// @Param refresh_token body string false "Refresh token (body) (for Swagger/testing)"
// @Param X-Refresh-Token header string false "Refresh token (header) (for Swagger/testing)"
// @Success 200 {object} map[string]interface{} "New access token and refresh cookie set"
// @Failure 401 {object} map[string]interface{} "Invalid or missing refresh token"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /svc/auth/refresh [post]
func RefreshHandler(c *gin.Context) {
	// 1) Try cookie
	cookie, err := c.Request.Cookie("refresh_token")
	var plain string
	if err == nil {
		plain = cookie.Value
	}

	// 2) Fallback to header (useful for Swagger / tooling)
	if plain == "" {
		plain = c.GetHeader("X-Refresh-Token")
	}

	// 3) Fallback to JSON body (useful for Swagger 'Try it out')
	if plain == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&body); err == nil && body.RefreshToken != "" {
			plain = body.RefreshToken
		}
	}

	if plain == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	newPlain, userID, err := services.VerifyAndRotateRefreshToken(plain)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}

	accessToken, _, err := services.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// set rotated refresh cookie
	setRefreshCookie(c.Writer, newPlain, time.Now().Add(7*24*time.Hour))

	c.JSON(http.StatusOK, gin.H{"token": accessToken, "expires_in": 15 * 60})
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

// LogoutHandler için Swagger tanımı
// @Summary Logout user
// @Description Revokes refresh token and clears cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully logged out"
// @Router /svc/auth/logout [post]
func LogoutHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err == nil {
		plain := cookie.Value
		parts := strings.SplitN(plain, ":", 2)
		if len(parts) == 2 {
			_ = services.RevokeRefreshTokenByID(parts[0])
		}
	}
	// clear cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("COOKIE_SECURE") != "false",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func setRefreshCookie(w http.ResponseWriter, plain string, expiry time.Time) {
	cookieSecure := true
	if os.Getenv("COOKIE_SECURE") == "false" {
		cookieSecure = false
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    plain,
		Path:     "/",
		HttpOnly: true,
		Secure:   cookieSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
		MaxAge:   int(expiry.Sub(time.Now()).Seconds()),
	})
}
