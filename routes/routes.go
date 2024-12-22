package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

// AuthRoutes defines routes for authentication
func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.LoginHandler)
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
