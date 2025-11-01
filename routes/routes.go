package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRoutes defines routes for authentication
func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/svc/auth")
	auth.Use(middlewares.RateLimitMiddleware())
	{
		auth.POST("/login-by-username", controllers.LoginByUsernameHandler)
		auth.POST("/login-by-email", controllers.LoginByEmailHandler)
		auth.POST("/login-by-phone", controllers.LoginByPhoneHandler)
		auth.POST("/refresh", controllers.RefreshHandler)
		auth.POST("/logout", controllers.LogoutHandler)
		auth.GET("/verify", controllers.VerifyEmailHandler)
		auth.POST("/send-verification/:userID", controllers.SendVerificationEmailHandler)
		auth.POST("/request-password-reset", controllers.RequestPasswordResetHandler)
		auth.POST("/reset-password", controllers.ResetPasswordHandler)

	}

	protected := router.Group("/admin")
	protected.Use(middlewares.AuthMiddleware()) // JWT Middleware
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to the admin dashboard"})
		})
	}
}

func MaintenanceRoutes(router *gin.Engine) {
	settings := router.Group("/maintenance")
	settings.Use(middlewares.AuthMiddleware()) // Yetkilendirme
	{
		settings.PUT("/", controllers.ToggleMaintenanceMode)
	}
}

// Health check (k8s readiness / liveness)
func HealthRoutes(router *gin.Engine) {
	router.GET("/healthz", func(c *gin.Context) {
		// Basit health check; DB/other servis check ekleyin
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
